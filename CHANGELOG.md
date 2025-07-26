# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
- 🌍 Custom port configuration
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