# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-07-28

### ✨ Added
- **Self-Update Command**: New `devshare update` command to automatically check for and install the latest release from GitHub, using the official installer scripts for each platform.
- **Cross-Platform Update**: Update logic downloads the correct archive, extracts the binary, and runs the appropriate installer (`install.sh` or `install.bat`).
- **Windows Version Info**: Windows binaries now include version information, company, product name, and more in the Properties dialog.
- **Windows App Icon**: Custom icon is now embedded in the Windows executable.
- **Embedded License Info**: License and copyright information are embedded in the Windows binary.
- **Improved Release Automation**: GoReleaser workflow now builds with `-ldflags="-s -w"` for smaller binaries and uses `windres` to embed resources.

### 🔐 Added
- **Password Protection**: New `--password` flag to secure development environments
- **Authentication Proxy**: Creates a secure proxy server with password authentication
- **Session Management**: Maintains authentication cookies for 1 hour
- **Beautiful Auth UI**: Modern, responsive authentication page
- **Secure Access**: Protects development work from unauthorized network access

### 🎯 Features
- **Reverse Proxy**: Forwards authenticated requests to the actual development server
- **Cookie-based Sessions**: Secure session management with HttpOnly cookies
- **Network Isolation**: Direct local access remains available while network access is protected
- **QR Code Integration**: QR codes point to the secure proxy URL when authentication is enabled
- **Port Management**: Automatically finds available ports for the authentication proxy

### 🛠️ Technical Details
- **Authentication Middleware**: Custom middleware for password verification
- **Proxy Implementation**: HTTP reverse proxy with header forwarding
- **Port Detection**: Automatic port availability checking
- **Error Handling**: Graceful handling of authentication failures
- **Resource Embedding**: Windows resource file (`versioninfo.rc`) and icon embedding via GoReleaser

### 🚀 Usage
```bash
# Basic usage (no authentication)
devshare

# With password protection
devshare --password=mysecret123

# Custom port with authentication
devshare 3000 --password=mysecret123

# Update DevShare to the latest version
devshare update
```

### 🔐 Authentication Flow
1. User starts DevShare with `--password` flag
2. DevShare starts the development server on the original port
3. DevShare creates an authentication proxy on a different port
4. Network access requires password authentication
5. Local access remains direct (no authentication required)
6. QR codes point to the secure proxy URL

---

## [1.0.0] - 2025-07-26

### 🎉 Added
- ✨ Initial release of DevShare CLI tool
- 🔍 **Project Detection**: Auto-detects project types (React, Next.js, Node.js, Go, Laravel)
- 🛠️ **Package Manager Detection**: Automatically detects npm, yarn, or pnpm
- 🌐 **LAN Sharing**: Exposes development servers to local network
- 📱 **QR Code Generation**: Generates QR codes for easy mobile access
- 🚀 **Smart Server Startup**: Automatically starts appropriate dev servers
- 💡 **Dependency Management**: Auto-installs dependencies if needed
- 🎨 **Beautiful CLI Interface**: Colored output with emojis and ASCII art
- 📊 **Network Information**: Displays local and network URLs
- 🔄 **Retry Logic**: Attempts multiple startup methods for Node.js projects

### 🎯 Features
- **React Support**: Works with Vite, Create React App, and custom React setups
- **Next.js Support**: Detects and runs Next.js development servers
- **Node.js Support**: Handles various Node.js project structures and entry points
- **Go Support**: Runs Go applications with `go run .`
- **Laravel Support**: Starts Laravel development server with `php artisan serve`
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Version Information**: Built-in version command with `--version` flag

### 🛠️ Technical Details
- **CLI Framework**: Built with Cobra for robust command-line interface
- **QR Code Library**: Uses qrterminal for terminal-based QR code generation
- **Network Detection**: Automatically finds local IP addresses
- **Error Handling**: Comprehensive error handling and user feedback
- **Modular Architecture**: Clean separation of concerns with dedicated packages

### 📦 Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output
- `github.com/mdp/qrterminal/v3` - QR code generation
- Standard Go libraries for file system and network operations

### 🚀 Usage
```bash
# Basic usage
devshare

# Check version
devshare --version
```

### 🎨 UI/UX
- ASCII art logo with DevShare branding
- Color-coded output for better readability
- Emoji indicators for different project types and actions
- Clear network information display
- User-friendly error messages

---

## [Unreleased]

### 🔮 Planned Features
- 🔧 Configuration file support
- 🔐 HTTPS support
- 📊 Usage analytics
---

## Release Notes

### v1.0.0 - Initial Release
This is the first stable release of DevShare, providing a complete solution for sharing development environments over LAN. The tool is production-ready and supports all major web development frameworks.

**Breaking Changes**: None (initial release)

**Migration Guide**: Not applicable (initial release)

---

*For more information, see [GitHub Repository](https://github.com/abdorrahmani/devshare)* 