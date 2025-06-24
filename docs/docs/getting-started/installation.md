# Installation Guide

This guide covers various ways to install and run Librebucket on your system.

## Prerequisites

Before installing Librebucket, ensure you have:

- **Go 1.22 or later** - [Download Go](https://golang.org/dl/)
- **Git 2.20 or later** - [Download Git](https://git-scm.com/downloads)
- **A supported operating system** - Linux, macOS, Windows, or ARM devices (Raspberry Pi, Orange Pi, etc.)

## Installation Methods

### From Source (Recommended)

Building from source gives you the latest features and allows customization.

#### 1. Clone the Repository

```bash
git clone https://github.com/standard-group/librebucket.git
cd librebucket
```

#### 2. Build the Application

```bash
go build -o librebucket main.go
```

#### 3. Run Librebucket

```bash
./librebucket
```

The server will start on `http://localhost:3000` by default.

### Using Docker (Coming Soon)

Docker provides an easy way to run Librebucket in a containerized environment.

```bash
# Pull the image
docker pull librebucket/librebucket:latest

# Run with persistent data
docker run -d \
  -p 3000:3000 \
  -v ./data:/app/config/data \
  --name librebucket \
  librebucket/librebucket:latest
```

### Pre-built Binaries (Coming Soon)

For convenience, we'll provide pre-built binaries for common platforms:

1. Visit the [releases page](https://github.com/standard-group/librebucket/releases)
2. Download the binary for your platform
3. Make it executable and run

## Platform-Specific Instructions

### Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go git

# CentOS/RHEL/Fedora
sudo dnf install golang git
# or
sudo yum install golang git

# Arch Linux
sudo pacman -S go git
```

### macOS

```bash
# Using Homebrew
brew install go git

# Using MacPorts
sudo port install go git
```

### Windows

1. Download and install Go from [golang.org](https://golang.org/dl/)
2. Download and install Git from [git-scm.com](https://git-scm.com/downloads)
3. Use PowerShell or Command Prompt to run the build commands

### ARM Devices (Raspberry Pi, etc.)

Librebucket runs great on ARM devices:

```bash
# Raspberry Pi OS
sudo apt update
sudo apt install golang-go git

# Build (same as other platforms)
go build -o librebucket main.go
./librebucket
```

## Verification

After installation, verify everything is working:

1. **Check the server is running:**
   ```bash
   curl http://localhost:3000
   ```

2. **Test the API:**
   ```bash
   curl http://localhost:3000/api/v1/health
   ```

3. **Open the web interface:**
   Navigate to `http://localhost:3000` in your browser

## Next Steps

After installation:

1. [Configure your server](configuration.md)
2. [Create your first repository](../user-guide/repositories.md)
3. [Set up user accounts](../user-guide/users.md)

## Troubleshooting

### Common Issues

#### Port Already in Use
If port 3000 is already in use, set a different port:

```bash
LIBREBUCKET_PORT=8080 ./librebucket
```

#### Permission Denied
Make sure the binary is executable:

```bash
chmod +x ./librebucket
```

#### Go Version Too Old
Update Go to version 1.22 or later:

```bash
go version  # Check current version
```

### Getting Help

If you encounter issues:

1. Check our [FAQ](../about/faq.md)
2. Search [existing issues](https://github.com/standard-group/librebucket/issues)
3. Create a [new issue](https://github.com/standard-group/librebucket/issues/new)
