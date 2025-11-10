# Makefile Usage Guide for Windows

## Problem: Colors Not Working on Windows

If you see garbled text like `\033[0;32m` when running `make` commands on Windows, it's because Windows Command Prompt doesn't support ANSI color codes properly.

## Solutions

### Option 1: Use Windows-Compatible Makefile (Recommended for Windows)

Replace your current `Makefile` with `Makefile.windows`:

```bash
# Backup original
copy Makefile Makefile.original

# Use Windows version
copy Makefile.windows Makefile
```

The Windows version uses simple text markers instead of colors:
- `[+]` - Action in progress
- `[OK]` - Success
- `[WARNING]` - Warning
- `[ERROR]` - Error

### Option 2: Use Git Bash (Best for Colors)

Install [Git for Windows](https://git-scm.com/download/win) and use Git Bash instead of Command Prompt or PowerShell. Git Bash supports ANSI colors.

```bash
# In Git Bash
make deps
make docker-up
```

### Option 3: Use Windows Terminal (Modern Windows)

Install [Windows Terminal](https://apps.microsoft.com/store/detail/windows-terminal/9N0DX20HK701) from Microsoft Store. It supports ANSI colors.

```powershell
# In Windows Terminal
make deps
make docker-up
```

### Option 4: Use WSL (Linux Subsystem)

Use Windows Subsystem for Linux for the full Linux experience:

```bash
# Install WSL first
wsl --install

# Then use make normally
make deps
make docker-up
```

## Common Windows Issues & Fixes

### Issue 1: "make: command not found"

**Solution:** Install Make for Windows

```bash
# Using Chocolatey
choco install make

# Or download from: http://gnuwin32.sourceforge.net/packages/make.htm
```

### Issue 2: Shell script comments cause errors

```
process_begin: CreateProcess(NULL, # comment, ...) failed.
```

**Solution:** This is fixed in `Makefile.windows` - no inline comments in recipe lines.

### Issue 3: Path separator issues

Windows uses `\` but Makefile uses `/`. 

**Solution:** Use forward slashes in Makefile (works on both systems):
```makefile
# Good (works everywhere)
./cmd/main.go

# Bad (Windows only)
.\cmd\main.go
```

### Issue 4: Date command format

Windows `date` command differs from Unix.

**Solution:** Use Git Bash or install GNU coreutils:
```bash
# Install GNU date on Windows
choco install gnuwin32-coreutils.install
```

## Quick Reference

### Basic Commands (All Platforms)

```bash
# Show help
make help

# Install dependencies
make deps

# Run locally
make run

# Build binary
make build

# Run tests
make test

# Docker operations
make docker-up          # Start services
make docker-down        # Stop services
make docker-logs        # View logs

# Database operations
make migrate-up         # Apply migrations
make db-backup          # Backup database
make db-connect         # Connect to DB
```

### Development Workflow

```bash
# 1. Setup (first time)
make deps
make install-tools

# 2. Start development
make docker-up          # Start services
make dev                # Run with hot reload

# 3. Before committing
make check              # Run all quality checks
make test               # Run tests

# 4. Cleanup
make docker-down        # Stop services
```

## File Structure

```
project/
├── Makefile                 # Unix/Linux/Mac version (with colors)
├── Makefile.windows        # Windows-compatible version (no colors)
└── MAKEFILE_GUIDE.md       # This guide
```

## Recommendation by Platform

| Platform | Recommended Makefile | Terminal |
|----------|---------------------|----------|
| Windows 11 | `Makefile` | Windows Terminal |
| Windows 10 | `Makefile.windows` | Git Bash |
| macOS | `Makefile` | Terminal.app |
| Linux | `Makefile` | Any terminal |
| WSL | `Makefile` | WSL terminal |

## Testing Your Setup

Run this command to test if colors work:

```bash
make version
```

**Expected output (with colors):**
```
Version:    1.0.0
Build Time: 2025-11-10_14:30:00
Git Commit: abc123d
```

**If you see garbled text:**
```
\033[0;32mVersion:    1.0.0\033[0m
```

→ Use `Makefile.windows` or switch to Git Bash/Windows Terminal

## Need Help?

1. Check which make you're using: `make --version`
2. Check your shell: `echo $SHELL` (bash/zsh) or `$PSVersionTable` (PowerShell)
3. Verify Docker is running: `docker --version`
4. Check if services are up: `docker-compose ps`

## Summary

- **Quick fix:** Use `Makefile.windows`
- **Best experience:** Use Git Bash or Windows Terminal
- **Long term:** Consider WSL for full Linux compatibility