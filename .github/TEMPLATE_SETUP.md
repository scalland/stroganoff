# GitHub Template Repository Setup Guide

This guide explains how to set up stroganoff as a GitHub template repository.

## What is a Template Repository?

A GitHub template repository allows other developers to quickly start new projects using your repository as a base. When users click "Use this template" on your repository, they get:

- ✅ All files and directory structure
- ✅ All commits from main branch
- ✅ README and documentation
- ✅ GitHub workflows and configurations
- ❌ Git history from the template
- ❌ Issue/PR history

## Prerequisites

1. Repository must be on GitHub
2. Repository must be public
3. You must have admin access to the repository

## Step 1: Prepare the Repository

### Update README.md
Ensure your README clearly explains:
- What the project is
- How to use it as a template
- What customizations are needed

### Include Documentation
- ✅ TEMPLATE.md - Template usage guide
- ✅ GETTING_STARTED.md - Setup instructions
- ✅ README.md - Project overview
- ✅ DEPLOYMENT.md - Deployment guide

### Verify Code Compiles
```bash
make build
make test
```

### Clean Up Template-Specific Files (Optional)
Files that might be template-specific:
- `TEMPLATE.md` - Can be removed after using template
- `template.properties` - Can be removed after using template
- `.github/TEMPLATE_*.md` - Can be removed after using template

But it's actually helpful to keep them for guidance!

## Step 2: Enable Template Repository on GitHub

1. Go to your repository on GitHub
2. Click **Settings** (gear icon)
3. Scroll down to **General** section
4. Check the box: **"Template repository"**
5. Click **Save** (if there's a button)

That's it! The repository is now a template.

## Step 3: Add Repository Metadata

### Topics
Go to **Settings → About** and add topics:
- `go`
- `golang`
- `cli`
- `cobra`
- `gin`
- `template`
- `boilerplate`
- `production-ready`
- `docker`
- `github-actions`

### Description
Update repository description to:
```
🚀 Production-ready Go CLI application template with Cobra, Gin, YAML config,
authentication, rate limiting, and comprehensive documentation.
```

### Website (Optional)
If you have documentation hosted somewhere, add the URL

## Step 4: Configure GitHub Features

### Enable Issues
- Settings → Features → Issues ✅

### Enable Discussions (Optional)
- Settings → Features → Discussions ✅

### Disable Unused Features (Optional)
- Settings → Features → Wikis ❌ (if not needed)
- Settings → Features → Projects ❌ (if not needed)

## Step 5: Set Up Security

### Enable Security Features
1. Settings → Code security and analysis
2. Enable:
   - ✅ Dependabot alerts
   - ✅ Secret scanning
   - ✅ Code scanning with CodeQL (optional)

## Step 6: Configure Branch Protection (Optional)

To maintain template quality:

1. Settings → Branches
2. Add rule for `main` branch:
   - Require pull request reviews: ✓
   - Require status checks: ✓
   - Require branches up to date: ✓
   - Include administrators: ✓

## Step 7: Document Usage

Create/update the following:

### `TEMPLATE.md`
Guide for using the template:
- How to use "Use this template"
- What to customize
- Quick start for new projects

### `.github/TEMPLATE_README.md`
Overview of template configuration

### GitHub Template Guide
Add to README or create separate guide:
```markdown
## Using This as a Template

This repository is a GitHub template repository. You can create a new repository
from this template using the "Use this template" button.

See [TEMPLATE.md](TEMPLATE.md) for detailed instructions.
```

## Step 8: Test the Template

### Before Publishing
1. Verify template creates repos correctly
2. Test that new repos can build: `make build`
3. Check that imports work correctly
4. Verify all files are present

### How to Test
1. Create a test repository from the template
2. Clone it locally
3. Update module name in go.mod
4. Run `make build` - should work
5. Run `make serve` - should work
6. Delete test repository

## Step 9: Announce the Template

Once template is live:

1. Update GitHub repository about section
2. Add to portfolio/personal site
3. Share on relevant communities:
   - r/golang
   - Go Slack channels
   - Dev communities
   - Twitter/X

## Maintenance

### Keep Template Updated
- Regularly update Go version
- Update dependencies: `go get -u ./...`
- Fix bugs and add features
- Update documentation

### Monitor Usage
- Check if users are using the template (via GitHub metrics)
- Collect feedback
- Fix issues reported

### Version Management
- Keep a clear CHANGELOG
- Use semantic versioning
- Document breaking changes

## Users Can Now

Users will see the **"Use this template"** button and can:

1. Click "Use this template"
2. Enter new repository name
3. Choose public/private
4. Create repository
5. Clone and customize:
   ```bash
   git clone https://github.com/username/my-app
   cd my-app

   # Update module name in go.mod
   sed -i 's/yourusername\/stroganoff/yourname\/my-app/g' go.mod

   # Build
   make build
   ```

## Example: What Users See

### Before Template Configuration
- "Use this template" button is **not visible**

### After Template Configuration
- GitHub shows: "Use this template" button
- Template shows up in search
- Users can generate new repos quickly

## Troubleshooting

### "Use this template" button not showing?
- Repository must be public
- Must have checked "Template repository" checkbox
- Might need to wait a few minutes for GitHub to update

### New repositories have template history?
- This shouldn't happen - GitHub creates fresh git history
- If it does, check repository settings

### Files missing in new repository?
- Make sure all files are committed and pushed to main branch
- Don't commit to other branches

## Next Steps

1. ✅ Enable template repository setting
2. ✅ Add topics and description
3. ✅ Enable security features
4. ✅ Create comprehensive documentation
5. ✅ Test creating a new repository from the template
6. ✅ Share with the community

## References

- [GitHub Template Repositories Docs](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-template-repository)
- [GitHub CLI for template repos](https://cli.github.com/)
- [Creating repositories from templates](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template)

---

**Your template repository is now ready for others to use!** 🚀
