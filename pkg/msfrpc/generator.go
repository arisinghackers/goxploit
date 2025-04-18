//go:generate go run github.com/arisinghackers/goxploit/cmd/generator
package msfrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MsfLibraryGenerator struct{}

func (g *MsfLibraryGenerator) GenerateLibrary() error {
	scraper := MsfPayloadScraper{}
	outputDir := "../../pkg/msfrpc/generated"

	apiMethods, err := scraper.GetArraysPayloadsFromWebsite()
	if err != nil {
		return err
	}

	filesMap := map[string]bool{}
	for _, method := range apiMethods {
		group := strings.Split(method[0], ".")[0]
		filesMap[group] = true
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for methodsGroup := range filesMap {
		className := strings.Title(methodsGroup) // + "ApiMethods"
		builder := &strings.Builder{}

		builder.WriteString("// Code generated by goxploit msfrpc LibraryGenerator. DO NOT EDIT.\n\n")
		builder.WriteString("package metasploit\n\n")
		builder.WriteString("import (")
		builder.WriteString("\n\t\"errors\"")
		builder.WriteString("\n\t\"github.com/arisinghackers/goxploit/pkg/msfrpc\"")
		builder.WriteString("\n)\n\n")
		builder.WriteString(fmt.Sprintf("type %s struct {\n\t*msfrpc.MsfRpcClient\n}\n\n", className))

		builder.WriteString(fmt.Sprintf("func New%s(client *msfrpc.MsfRpcClient) *%s {\n", className, className))
		builder.WriteString(fmt.Sprintf("\treturn &%s{client}\n", className))
		builder.WriteString("}\n\n")

		// Métodos
		for _, method := range apiMethods {
			if strings.Split(method[0], ".")[0] != methodsGroup {
				continue
			}

			rawMethod := strings.SplitN(method[0], ".", 2)[1]
			funcName := toCamelCase(rawMethod)

			params := []string{}
			args := []string{fmt.Sprintf(`"%s"`, method[0])}

			for _, val := range method[1:] {
				param := strings.Trim(val, "<>")
				if strings.Contains(val, "<") {
					if strings.Contains(param, "token") {
						args = append(args, "c.Token")
					} else {
						args = append(args, param)
					}
					params = append(params, fmt.Sprintf("%s string", param))
				} else {
					if methodsGroup == "console" && strings.ToLower(val) == "inputCommand" {
						args = append(args, fmt.Sprintf("%s+\"\\n\"", val))
					} else {
						args = append(args, val)
					}
					params = append(params, fmt.Sprintf("%s string", val))
				}
			}

			builder.WriteString(fmt.Sprintf("func (c *%s) %s(%s) (map[string]interface{}, error) {\n", className, funcName, strings.Join(params, ", ")))
			builder.WriteString(fmt.Sprintf("\tresp, err := c.MsfRequest(%s)\n", buildArgsString(args)))
			builder.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
			builder.WriteString("\tif result, ok := resp[\"result\"]; ok && result == \"failure\" {\n")
			builder.WriteString("\t\treturn nil, errors.New(\"Unprocessable Content\")\n\t}\n")
			builder.WriteString("\tif errMsg, ok := resp[\"error_message\"]; ok {\n")
			builder.WriteString("\t\treturn nil, errors.New(errMsg.(string))\n\t}\n")
			builder.WriteString("\treturn resp, nil\n}\n\n")
		}

		// Define the file path
		filePath := filepath.Join(outputDir, fmt.Sprintf("%s.go", strings.ToLower(className)))
		if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
			return err
		}

	}

	return nil
}

func buildArgsString(args []string) string {
	return "[]interface{}{" + strings.Join(args, ", ") + "}"
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
