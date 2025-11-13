# GitHub Template Repository Setup

This directory contains GitHub-specific configurations for this template repository.

## Files in this Directory

### `workflows/build.yml`
- Automated CI/CD pipeline
- Builds on push and creates releases
- Supports multi-OS/arch builds

### `ISSUE_TEMPLATE/`
- Issue templates for bug reports and feature requests (can be added)

### `PULL_REQUEST_TEMPLATE.md`
- Pull request template (can be added)

## Enabling as GitHub Template

To enable this as a GitHub template repository:

1. Go to repository **Settings**
2. Scroll to **Template repository** section
3. Check the **"Template repository"** checkbox
4. Click **Save**

## Users Can Now

- Click **"Use this template"** button on the repository page
- Create a new repository from this template
- All files and directory structure will be copied
- Fresh git history (no commit history from template)

## What Gets Copied

✅ All source code files
✅ Directory structure
✅ Configuration files
✅ Documentation
✅ GitHub Actions workflows
✅ Makefile
✅ Docker files

❌ Git history (starts fresh)
❌ git-related files (.git directory)

## After Using the Template

Users should:

1. Update `go.mod` with their module path
2. Update imports in Go files
3. Customize `README.md`
4. Modify `Makefile` if needed
5. Update `.github/workflows/build.yml` for their repo
6. Remove `TEMPLATE.md` and `template.properties` if desired

## Testing the Template

Before publishing:

1. Make sure all files are properly committed
2. Verify directory structure is complete
3. Test that code compiles: `make build`
4. Check documentation is accurate
5. Verify GitHub Actions workflow syntax

## Template Visibility

- Repository must be public to use as a template
- Template will appear in GitHub's template search
- Use meaningful README for discoverability
- Add repository topics for better categorization
