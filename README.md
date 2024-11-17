````markdown
# Stackpack 🎁

Transform GitHub repositories into AI-friendly single files - perfect for LLMs like Claude, ChatGPT, and Gemini.

## Quick Start

### CLI

```bash
# Install
go install github.com/axzilla/stackpack/cmd/cli@latest

# Use
stackpack --remote=https://github.com/user/repo
stackpack --include="*.go,*.md"
```
````

### Web App

Visit [stackpack.xyz](https://stackpack.xyz)

> **Note**: The web app has a rate limit of 60 requests per hour per IP for GitHub API access. For higher limits, consider using the CLI with a personal GitHub token.

## Features

- 📦 Single-file output optimized for AI consumption
- 🔍 Include/exclude files via patterns (e.g. `*.go`, `*.md`)
- 🌐 Works with local & remote GitHub repositories
- 🌙 Dark mode support

## CLI Options

```bash
stackpack [flags]

--output  output.txt    # Output filename
--remote  URL           # GitHub repository URL
--include "*.go,*.md"   # Files to include
--exclude "*.png,*.jpg" # Files to exclude
```

## Development

Requirements:

- Go 1.22+
- Node.js
- GitHub Personal Access Token (for higher API rate limits)

```bash
# Clone and setup
git clone https://github.com/axzilla/stackpack
cd stackpack

# Create .env file
echo "GITHUB_TOKEN=your_github_token" > .env  # Get token from GitHub.com -> Settings -> Developer Settings -> PAT

# Start development
make dev
```

To get a GitHub token:

1. Go to GitHub.com -> Settings -> Developer Settings -> Personal Access Tokens
2. Generate new token with `public_repo` scope
3. Add to .env file as GITHUB_TOKEN

## License

MIT
