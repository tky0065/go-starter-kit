# Installation Guide

This guide details all methods to install `create-go-starter` on your system.

!!! note "Translation in progress"
    This page is being translated from French. For the complete documentation, please refer to the [French version](../installation/).

## Prerequisites

### System Requirements

- **Go 1.25 or higher** - Check your version with `go version`
- **Git** - To clone the repository
- **Disk space** - ~50MB for the tool and its dependencies

### Check Go Installation

```bash
go version
# Should display: go version go1.25.x ...
```

If Go is not installed, download it from [golang.org/dl](https://golang.org/dl/).

## Method 1: Direct Installation (Recommended)

Global installation with a single command, **without cloning the repository**.

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

This command:

- Automatically downloads the code from GitHub
- Compiles the binary
- Installs it in `$GOPATH/bin` (usually `~/go/bin`)
- Makes it globally available

### Verification

```bash
create-go-starter --help
```

**Note**: Make sure `$GOPATH/bin` is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Method 2: Build from Sources

Recommended for contributors or customization:

```bash
# Clone the repository
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit

# Build
go build -o create-go-starter ./cmd/create-go-starter

# Or with Makefile
make build
```

## Updating

### Via go install

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

### From sources

```bash
cd /path/to/go-starter-kit
git pull origin main
make build
```

## Troubleshooting

### Problem: "command not found: create-go-starter"

**Solutions**:

1. Reload shell cache: `hash -r`
2. Restart your terminal
3. Check PATH: `echo $PATH | grep "$(go env GOPATH)/bin"`

### Problem: "permission denied" on macOS

```bash
chmod +x create-go-starter
xattr -d com.apple.quarantine create-go-starter
```

## Next Steps

- [Usage Guide](usage.md) - Learn to use the tool
- [Generated Project Guide](generated-project-guide.md) - Develop with generated projects
- [Contributing Guide](contributing.md) - Contribute to the project
