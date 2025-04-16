package msfrpc

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type MsfPayloadScraper struct{}

func NewMsfPayloadScraper() *MsfPayloadScraper {
	return &MsfPayloadScraper{}
}

func (s *MsfPayloadScraper) GetArraysPayloadsFromWebsite() ([][]string, error) {
	var payloads [][]string
	seen := map[string]bool{}

	c := colly.NewCollector()

	c.OnHTML(".language-bash > .token-line > .code-line-content", func(e *colly.HTMLElement) {
		line := e.Text
		var result []string

		if strings.Contains(line, "[") && strings.Contains(line, "{") {
			result = s.stringToArrayProcessor(strings.ReplaceAll(line, "{", `"Options" ]`))
		} else if !strings.Contains(line, "[") || strings.Contains(line, "Bad") {
			return
		} else {
			result = s.stringToArrayProcessor(line)
		}

		joined := strings.Join(result, ",")
		if !seen[joined] {
			seen[joined] = true
			payloads = append(payloads, result)
		}
	})

	err := c.Visit("https://docs.rapid7.com/metasploit/standard-api-methods-reference")
	if err != nil {
		return nil, err
	}

	return payloads, nil
}

func (s *MsfPayloadScraper) stringToArrayProcessor(input string) []string {
	replacements := map[string]string{
		"MyUserName": "userName", "MyPassword": "userPassword",
		"ThreadID": "threadId", "JobID": "jobId",
		"ConsoleID": "consoleId", "SessionID": "sessionId",
		"0": "ConsoleId", "versionn": "InputCommand",
		`"ReadPointer ]`: "InputCommand", "idn": "InputCommand",
		"1.2.3.4": "IpAddress", "4444": "Port", "1": "SessionId",
		"scriptname": "scriptName",
	}

	var result []string
	split := strings.Split(input, "\"")
	re := regexp.MustCompile(`[#$%^&*()+=\-\[\];,\/{}|":?~\\ ]`)

	for _, val := range split {
		clean := re.ReplaceAllString(val, "")
		if clean == "" {
			continue
		}
		if repl, found := replacements[clean]; found {
			result = append(result, repl)
		} else {
			result = append(result, clean)
		}
	}

	return result
}
