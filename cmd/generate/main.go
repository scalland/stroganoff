package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/stroganoff/internal/generator"
)

func main() {
	var (
		projectName  = flag.String("name", "", "Project name (required)")
		moduleName   = flag.String("module", "", "Go module name, e.g., github.com/username/projectname (required)")
		appName      = flag.String("app", "", "Application name (required)")
		outputDir    = flag.String("output", "", "Output directory path (required)")
		templatePath = flag.String("template", "", "Path to stroganoff template directory (default: current stroganoff)")
	)

	flag.Parse()

	// Validate required flags
	if *projectName == "" || *moduleName == "" || *appName == "" || *outputDir == "" {
		fmt.Fprintf(os.Stderr, "Error: all flags are required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -name PROJECT_NAME -module MODULE_NAME -app APP_NAME -output OUTPUT_DIR [-template TEMPLATE_PATH]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -name myapp -module github.com/user/myapp -app MyApp -output ./myapp\n", os.Args[0])
		os.Exit(1)
	}

	// If template path not provided, use current directory
	if *templatePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		*templatePath = cwd
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Get absolute paths
	templateAbs, err := filepath.Abs(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving template path: %v\n", err)
		os.Exit(1)
	}

	outputAbs, err := filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output path: %v\n", err)
		os.Exit(1)
	}

	// Create generator
	gen, err := generator.New(generator.Config{
		ProjectName:  *projectName,
		ModuleName:   *moduleName,
		AppName:      *appName,
		OutputDir:    outputAbs,
		TemplatePath: templateAbs,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
		os.Exit(1)
	}

	// Generate project
	fmt.Printf("Generating project '%s' from template...\n", *projectName)
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Project generated successfully at %s\n", outputAbs)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", outputAbs)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  make build\n")
}
