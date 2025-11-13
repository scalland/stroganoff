package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ProjectName:  "myapp",
				ModuleName:   "github.com/user/myapp",
				AppName:      "MyApp",
				OutputDir:    tmpDir,
				TemplatePath: tmpDir,
			},
			wantErr: false,
		},
		{
			name: "missing project name",
			config: Config{
				ModuleName:   "github.com/user/myapp",
				AppName:      "MyApp",
				OutputDir:    tmpDir,
				TemplatePath: tmpDir,
			},
			wantErr: true,
		},
		{
			name: "missing module name",
			config: Config{
				ProjectName:  "myapp",
				AppName:      "MyApp",
				OutputDir:    tmpDir,
				TemplatePath: tmpDir,
			},
			wantErr: true,
		},
		{
			name: "missing app name",
			config: Config{
				ProjectName:  "myapp",
				ModuleName:   "github.com/user/myapp",
				OutputDir:    tmpDir,
				TemplatePath: tmpDir,
			},
			wantErr: true,
		},
		{
			name: "missing output dir",
			config: Config{
				ProjectName:  "myapp",
				ModuleName:   "github.com/user/myapp",
				AppName:      "MyApp",
				TemplatePath: tmpDir,
			},
			wantErr: true,
		},
		{
			name: "missing template path",
			config: Config{
				ProjectName: "myapp",
				ModuleName:  "github.com/user/myapp",
				AppName:     "MyApp",
				OutputDir:   tmpDir,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeneratorNew(t *testing.T) {
	tmpDir := t.TempDir()

	config := Config{
		ProjectName:  "myapp",
		ModuleName:   "github.com/user/myapp",
		AppName:      "MyApp",
		OutputDir:    tmpDir,
		TemplatePath: tmpDir,
	}

	gen, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if gen.config.ProjectName != "myapp" {
		t.Errorf("ProjectName = %v, want %v", gen.config.ProjectName, "myapp")
	}

	// Check replacements
	if _, ok := gen.replacements["github.com/yourusername/stroganoff"]; !ok {
		t.Error("replacements missing module name mapping")
	}
	if _, ok := gen.replacements["stroganoff"]; !ok {
		t.Error("replacements missing project name mapping")
	}
}

func TestReplaceContent(t *testing.T) {
	tmpDir := t.TempDir()

	config := Config{
		ProjectName:  "myapp",
		ModuleName:   "github.com/user/myapp",
		AppName:      "MyApp",
		OutputDir:    tmpDir,
		TemplatePath: tmpDir,
	}

	gen, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "replace module name",
			input:    "module github.com/yourusername/stroganoff",
			expected: "module github.com/user/myapp",
		},
		{
			name:     "replace project name",
			input:    "stroganoff version 1.0",
			expected: "myapp version 1.0",
		},
		{
			name:     "multiple replacements",
			input:    "import github.com/yourusername/stroganoff/pkg and run stroganoff",
			expected: "import github.com/user/myapp/pkg and run myapp",
		},
		{
			name:     "replace uppercase",
			input:    "Welcome to STROGANOFF",
			expected: "Welcome to MYAPP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.replaceContent(tt.input)
			if result != tt.expected {
				t.Errorf("replaceContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name    string
		entry   string
		isRoot  bool
		wantSkip bool
	}{
		{
			name:    "skip .git at root",
			entry:   ".git",
			isRoot:  true,
			wantSkip: true,
		},
		{
			name:    "don't skip .git in subdirs",
			entry:   ".git",
			isRoot:  false,
			wantSkip: false,
		},
		{
			name:    "skip .github at root",
			entry:   ".github",
			isRoot:  true,
			wantSkip: true,
		},
		{
			name:    "skip generator at root",
			entry:   "cmd/generate",
			isRoot:  true,
			wantSkip: true,
		},
		{
			name:    "don't skip regular files",
			entry:   "main.go",
			isRoot:  true,
			wantSkip: false,
		},
		{
			name:    "skip .DS_Store everywhere",
			entry:   ".DS_Store",
			isRoot:  true,
			wantSkip: true,
		},
		{
			name:    "skip .DS_Store in subdirs",
			entry:   ".DS_Store",
			isRoot:  false,
			wantSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkip(tt.entry, tt.isRoot)
			if result != tt.wantSkip {
				t.Errorf("shouldSkip() = %v, want %v", result, tt.wantSkip)
			}
		})
	}
}

func TestShouldProcessFile(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		should bool
	}{
		{
			name:   "process .go file",
			path:   "/path/to/main.go",
			should: true,
		},
		{
			name:   "process go.mod",
			path:   "/path/to/go.mod",
			should: true,
		},
		{
			name:   "process Makefile",
			path:   "/path/to/Makefile",
			should: true,
		},
		{
			name:   "process Dockerfile",
			path:   "/path/to/Dockerfile",
			should: true,
		},
		{
			name:   "process .yaml file",
			path:   "/path/to/config.yaml",
			should: true,
		},
		{
			name:   "process .yml file",
			path:   "/path/to/config.yml",
			should: true,
		},
		{
			name:   "don't process binary file",
			path:   "/path/to/image.png",
			should: false,
		},
		{
			name:   "don't process executable",
			path:   "/path/to/binary",
			should: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldProcessFile(tt.path)
			if result != tt.should {
				t.Errorf("shouldProcessFile() = %v, want %v", result, tt.should)
			}
		})
	}
}

func TestGenerateIntegration(t *testing.T) {
	// Create a temporary template directory
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	// Create some template files
	goModContent := `module github.com/yourusername/stroganoff

go 1.21
`
	mainGoContent := `package main

import (
	"github.com/yourusername/stroganoff/pkg/example"
)

func main() {
	example.Run()
}
`

	if err := os.WriteFile(filepath.Join(templateDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create test go.mod: %v", err)
	}

	mainDir := filepath.Join(templateDir, "cmd", "stroganoff")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		t.Fatalf("failed to create cmd directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(mainDir, "main.go"), []byte(mainGoContent), 0644); err != nil {
		t.Fatalf("failed to create test main.go: %v", err)
	}

	// Create generator
	config := Config{
		ProjectName:  "myapp",
		ModuleName:   "github.com/user/myapp",
		AppName:      "MyApp",
		OutputDir:    outputDir,
		TemplatePath: templateDir,
	}

	gen, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Generate project
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify generated files
	generatedGoMod := filepath.Join(outputDir, "go.mod")
	if _, err := os.Stat(generatedGoMod); err != nil {
		t.Errorf("go.mod not found in output directory: %v", err)
	}

	// Check content was replaced
	content, err := os.ReadFile(generatedGoMod)
	if err != nil {
		t.Fatalf("failed to read generated go.mod: %v", err)
	}

	if !strings.Contains(string(content), "github.com/user/myapp") {
		t.Error("module name not replaced in go.mod")
	}

	if strings.Contains(string(content), "github.com/yourusername/stroganoff") {
		t.Error("template module name not replaced in go.mod")
	}

	// Check cmd/myapp directory was created
	newCmdDir := filepath.Join(outputDir, "cmd", "myapp")
	if _, err := os.Stat(newCmdDir); err != nil {
		t.Errorf("cmd/myapp directory not found: %v", err)
	}
}
