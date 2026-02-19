package msfrpc

import codegen "github.com/arisinghackers/goxploit/internal/generator"

// MsfPayloadScraper is a compatibility alias for the internal generator scraper.
// External users should not rely on documentation scraping as a runtime dependency.
type MsfPayloadScraper struct {
	inner *codegen.PayloadScraper
}

func NewMsfPayloadScraper() *MsfPayloadScraper {
	return &MsfPayloadScraper{inner: codegen.NewPayloadScraper()}
}

func (s *MsfPayloadScraper) GetArraysPayloadsFromWebsite() ([][]string, error) {
	if s == nil || s.inner == nil {
		return codegen.NewPayloadScraper().GetArraysPayloadsFromWebsite()
	}
	return s.inner.GetArraysPayloadsFromWebsite()
}
