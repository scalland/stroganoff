# Stroganoff Template Generator

The stroganoff project now includes a powerful template generator that allows you to create new Go CLI projects based on the stroganoff template. This tool automatically handles all module path replacements, directory structure creation, and file processing.

## Features

- **Automatic Module Path Replacement**: Replaces all Go import paths and module names throughout the project
- **Project Structure Cloning**: Creates a complete project structure from the template
- **Smart File Processing**: Only processes text files (Go, YAML, Markdown, etc.), leaves binary files untouched
- **Directory Renaming**: Automatically renames directories to match your project name (e.g., `cmd/stroganoff` → `cmd/myapp`)
- **Comprehensive Validation**: Validates all input parameters before generation

## Installation

### Build the Generator

```bash
make build-generator
```

This creates the `stroganoff-generate` binary in the `dist/` directory.

### Add to PATH (Optional)

```bash
cp dist/stroganoff-generate /usr/local/bin/
# or add dist/ to your PATH
```

## Usage

### Basic Usage

```bash
stroganoff-generate \
  -name myapp \
  -module github.com/username/myapp \
  -app MyApp \
  -output ./myapp
```

### Parameters

- **`-name`** (required): Project name in lowercase (e.g., `myapp`)
  - Used for directory names and binary naming
  - Example: `myapp`

- **`-module`** (required): Go module name
  - Standard Go module path format
  - Example: `github.com/username/myapp`

- **`-app`** (required): Application name for display purposes
  - Used in help text, documentation, and comments
  - Example: `MyApp`

- **`-output`** (required): Output directory where the project will be created
  - Will be created if it doesn't exist
  - Example: `/path/to/new/project`

- **`-template`** (optional): Path to the template directory
  - Defaults to current directory if not specified
  - Use this if you want to use stroganoff from a different location

### Examples

#### Example 1: Basic Project Generation

```bash
stroganoff-generate \
  -name hello-world \
  -module github.com/myusername/hello-world \
  -app HelloWorld \
  -output ~/projects/hello-world
```

#### Example 2: Using a Specific Template Directory

```bash
stroganoff-generate \
  -name myservice \
  -module github.com/company/myservice \
  -app MyService \
  -output ./myservice \
  -template ~/stroganoff
```

#### Example 3: In Current Directory

```bash
cd stroganoff
./dist/stroganoff-generate \
  -name newproject \
  -module github.com/user/newproject \
  -app NewProject \
  -output ../newproject
```

## What Gets Generated

The template generator creates a complete project with:

### Directory Structure
```
myapp/
├── cmd/
│   └── myapp/
│       ├── main.go
│       └── commands/
│           ├── root.go
│           ├── version.go
│           ├── serve.go
│           └── ... (other commands)
├── internal/
│   ├── config/
│   ├── web/
│   ├── upgrade/
│   └── ... (other packages)
├── pkg/
│   ├── auth/
│   ├── ratelimit/
│   └── version/
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── config.example.yaml
├── go.mod
├── go.sum
└── ... (other files)
```

### What's Replaced

The generator automatically replaces:

1. **Module Path**: `github.com/yourusername/stroganoff` → `github.com/username/myapp`
2. **Project Name**: `stroganoff` → `myapp`
3. **Uppercase Name**: `STROGANOFF` → `MYAPP`
4. **Title Case**: `Stroganoff` → `Myapp`

### Files Processed

The generator processes these file types:
- `.go` - Go source files
- `.mod` - Go module files
- `.yaml`, `.yml` - YAML configuration
- `.json` - JSON files
- `.toml` - TOML files
- `.md` - Markdown documentation
- `.txt` - Text files
- `.sh` - Shell scripts
- `Makefile`, `Dockerfile` - Build files

Binary files (images, compiled binaries) are copied without modification.

### What's Skipped

These files/directories are skipped during generation:
- `.git/` - Version control directory
- `.github/` - GitHub specific files
- `cmd/generate/` - The generator itself
- `internal/generator/` - Generator package
- `node_modules/` - Node dependencies
- `dist/` - Build artifacts
- `bin/` - Binary files
- `.gitignore` - Git ignore file
- `template.properties` - Template configuration

## After Generation

After running the generator, follow these steps to complete setup:

```bash
cd myapp

# Download and verify dependencies
go mod tidy

# Build the project
make build

# Run the CLI
./dist/myapp --help

# Start the server (optional)
make serve

# Run tests
make test
```

## Template Replacements

The generator uses a straightforward replacement strategy. Here's what gets replaced in your generated code:

### Go Module Paths

Before:
```go
import "github.com/yourusername/stroganoff/pkg/auth"
```

After:
```go
import "github.com/username/myapp/pkg/auth"
```

### Configuration Files

Before:
```yaml
name: stroganoff
module: github.com/yourusername/stroganoff
```

After:
```yaml
name: myapp
module: github.com/username/myapp
```

### Documentation

Before:
```markdown
# Stroganoff

Welcome to stroganoff - A professional Go CLI template
```

After:
```markdown
# MyApp

Welcome to myapp - A professional Go CLI template
```

## Architecture

The generator package consists of:

- **Generator**: Main orchestrator that handles the generation process
- **Config**: Validates all input parameters
- **File Processing**: Handles text file processing and replacement
- **Directory Copy**: Recursively copies and processes directory trees

See [internal/generator/](internal/generator/) for implementation details.

## Customization

### Extending the Generator

To add additional replacements, modify the `replacements` map in `generator.go`:

```go
g := &Generator{
    config: config,
    replacements: map[string]string{
        "github.com/yourusername/stroganoff": config.ModuleName,
        "stroganoff": config.ProjectName,
        // Add custom replacements here
    },
}
```

### Changing Skipped Files

To change which files are skipped, modify the `shouldSkip()` function in `generator.go`.

### Changing Processed File Types

To add or remove file types that should be processed for replacements, modify the `shouldProcessFile()` function.

## Troubleshooting

### "Permission denied" error

Make sure the binary is executable:
```bash
chmod +x dist/stroganoff-generate
```

### Module path not being replaced

Ensure you're using the correct module path format:
- ✓ `github.com/username/projectname`
- ✗ `github.com/username/projectname/` (trailing slash)
- ✗ `http://github.com/username/projectname` (protocol prefix)

### Directories not renamed

Project name must be lowercase and match your directory naming convention:
- ✓ `myapp`, `my-app`, `my_app`
- ✗ `MyApp`, `MYAPP`

## Testing

Run the generator tests:

```bash
go test -v ./internal/generator
```

Run all tests including generator:

```bash
make test
```

The test suite includes:
- Configuration validation
- File processing logic
- Directory creation and copying
- Integration tests for complete generation

## Development

### Building for Multiple Platforms

To build the generator for multiple platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dist/stroganoff-generate-linux ./cmd/generate

# macOS
GOOS=darwin GOARCH=arm64 go build -o dist/stroganoff-generate-mac-arm64 ./cmd/generate

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/stroganoff-generate.exe ./cmd/generate
```

## License

This template generator is part of the stroganoff project and is licensed under the MIT License.
