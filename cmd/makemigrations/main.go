package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Gahdloot/GormCase/pkg/migration/generator"
)

func main() {
	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Usage: makemigrations <model_file>")
		os.Exit(1)
	}

	modelPath := args[0]

	// Generate migration
	if err := generator.GenerateMigration(modelPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
