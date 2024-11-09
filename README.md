# Stackpack 🎁

Transform your GitHub repositories into AI-friendly formats, perfect for LLMs like Claude, ChatGPT, and Gemini.

## Features

- 📦 Pack entire repositories into a single, structured file
- 🔍 Include/exclude files using glob patterns
- 🌐 Support for local and remote GitHub repositories
- 🌙 Dark/Light mode support
- 🎯 Optimized output format for AI consumption

## Usage

### Web App

Visit [stackpack.xyz](https://stackpack.xyz) to use the web interface.

1. Enter a GitHub repository URL
2. (Optional) Specify include/exclude patterns
3. Click "Generate Pack"

### CLI

Install:

```bash
go install github.com/axzilla/stackpack/cmd/cli@latest
```

Basic usage:

```bash
# Pack local directory
stackpack --output=output.txt

# Pack remote repository
stackpack --remote=https://github.com/user/repo

# With include/exclude patterns
stackpack --include="*.go,*.md" --exclude="test/*"
```

Available flags:

- `--output`: Output file name (default: "stackpack.txt")
- `--include`: Glob patterns to include (comma-separated), e.g. "_\*.go,_\*.md"
- `--exclude`: Glob patterns to ignore (comma-separated), e.g. "_\*.svg,_\*.png"
- `--remote`: GitHub repository URL to pack

## Examples

Pack only Go and Markdown files:

```bash
stackpack --include="*.go,*.md"
```

Pack a remote repository excluding images:

```bash
stackpack --remote=https://github.com/user/repo --exclude="*.png,*.jpg"
```

## Development

Requirements:

- Go 1.22+
- Node.js (for web development)

Setup:

```bash
# Clone repository
git clone https://github.com/axzilla/stackpack
cd stackpack

# Install dependencies
go mod download
npm install

# Run web development server
make dev

# Build CLI
go build -o stackpack cmd/cli/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
