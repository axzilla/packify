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

### Web App

Visit [stackpack.xyz](https://stackpack.xyz)

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

```bash
git clone https://github.com/axzilla/stackpack
cd stackpack
make dev
```

Requires Go 1.22+ and Node.js.

## License

MIT
