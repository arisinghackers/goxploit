package main

import (
	"log"

	codegen "github.com/arisinghackers/goxploit/internal/generator"
)

func main() {
	generator := codegen.NewLibraryGenerator()
	if err := generator.GenerateLibrary(); err != nil {
		log.Fatal(err)
	}
}
