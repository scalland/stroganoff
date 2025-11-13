package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Config contains the parameters for generating a new project
type Config struct {
	ProjectName  string
	ModuleName   string
	AppName      string
	OutputDir    string
	TemplatePath string
}

// Generator handles project generation from template
type Generator struct {
	config Config
	// Replacement mapping
	replacements map[string]string
}

// New creates a new generator with the given config
func New(config Config) (*Generator, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	g := &Generator{
		config: config,
		replacements: map[string]string{
			// Replace template module name with new module name
			"github.com/yourusername/stroganoff": config.ModuleName,
			// Replace template project name with new project name
			"stroganoff": config.ProjectName,
			// Replace template app name (for display purposes)
			"STROGANOFF": strings.ToUpper(config.ProjectName),
			"Stroganoff": toTitleCase(config.ProjectName),
		},
	}

	return g, nil
}

// Validate checks if the config is valid
func (c *Config) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if c.ModuleName == "" {
		return fmt.Errorf("module name is required")
	}
	if c.AppName == "" {
		return fmt.Errorf("app name is required")
	}
	if c.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}
	if c.TemplatePath == "" {
		return fmt.Errorf("template path is required")
	}

	// Check if template path exists
	if info, err := os.Stat(c.TemplatePath); err != nil {
		return fmt.Errorf("template path does not exist: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("template path is not a directory")
	}

	return nil
}

// Generate creates a new project from the template
func (g *Generator) Generate() error {
	// Copy all files from template to output directory
	return g.copyDir(g.config.TemplatePath, g.config.OutputDir, true)
}

// copyDir recursively copies a directory with template processing
func (g *Generator) copyDir(srcDir, dstDir string, isRoot bool) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcDir, err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dstDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstName := entry.Name()

		// Skip certain directories and files that shouldn't be copied
		if shouldSkip(entry.Name(), isRoot) {
			continue
		}

		// Replace template project name in directory names
		for old, new := range g.replacements {
			if strings.Contains(dstName, old) && dstName != "go.mod" && dstName != "go.sum" {
				dstName = strings.ReplaceAll(dstName, old, new)
				break
			}
		}

		dstPath := filepath.Join(dstDir, dstName)

		if entry.IsDir() {
			if err := g.copyDir(srcPath, dstPath, false); err != nil {
				return err
			}
		} else {
			if err := g.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file, processing templates
func (g *Generator) copyFile(srcPath, dstPath string) error {
	// Check if this is a text file that needs processing
	if shouldProcessFile(srcPath) {
		return g.copyAndProcessFile(srcPath, dstPath)
	}

	// For binary files, just copy directly
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dstPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file %s: %w", srcPath, err)
	}

	return nil
}

// copyAndProcessFile copies a file and processes template variables
func (g *Generator) copyAndProcessFile(srcPath, dstPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}

	processedContent := g.replaceContent(string(content))

	if err := os.WriteFile(dstPath, []byte(processedContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dstPath, err)
	}

	return nil
}

// replaceContent replaces all template variables in content
func (g *Generator) replaceContent(content string) string {
	result := content
	for old, new := range g.replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// shouldSkip determines if a file/directory should be skipped during copy
func shouldSkip(name string, isRoot bool) bool {
	// Skip these at root level
	if isRoot {
		skipList := map[string]bool{
			".git":                 true,
			".github":              true,
			"node_modules":         true,
			"dist":                 true,
			"bin":                  true,
			"cmd/generate":         true, // Don't copy the generator itself
			"internal/generator":   true, // Don't copy generator package
			".gitignore":           true,
			".env":                 true,
			".DS_Store":            true,
			"template.properties":  true, // Don't copy template config
		}
		return skipList[name]
	}

	// Skip common non-essential files everywhere
	skipList := map[string]bool{
		".DS_Store": true,
		".gitkeep":  true,
	}

	return skipList[name]
}

// shouldProcessFile determines if a file should have template variables replaced
func shouldProcessFile(path string) bool {
	ext := filepath.Ext(path)
	processableExts := map[string]bool{
		".go":         true,
		".mod":        true,
		".yaml":       true,
		".yml":        true,
		".json":       true,
		".toml":       true,
		".md":         true,
		".txt":        true,
		".sh":         true,
		".Makefile":   true,
		".dockerfile": true,
	}

	if processableExts[ext] {
		return true
	}

	// Check filename for special cases
	name := filepath.Base(path)
	if name == "Makefile" || name == "Dockerfile" || name == "go.mod" {
		return true
	}

	return false
}

// toTitleCase converts a string to title case
func toTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
