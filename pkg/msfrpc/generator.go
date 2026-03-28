//go:generate go run github.com/arisinghackers/goxploit/cmd/generator
package msfrpc

import codegen "github.com/arisinghackers/goxploit/internal/generator"

// MsfLibraryGenerator is a compatibility wrapper around the internal generator package.
// Prefer running `go run ./cmd/generator` directly.
type MsfLibraryGenerator struct{}

func (g *MsfLibraryGenerator) GenerateLibrary() error {
	return codegen.NewLibraryGenerator().GenerateLibrary()
}
