[![GitHub Repo](https://img.shields.io/badge/GitHub-Repo-blue?logo=github)](https://github.com/abdorrahmani/devshare)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?logo=open-source-initiative)](https://github.com/abdorrahmani/devshare/blob/master/LICENSE)
![Static Badge](https://img.shields.io/badge/Built%20with-%F0%9F%92%99%20by%20Anophel-blue)

---

# DevShare

> **Share your development environment over LAN in seconds!**

DevShare is a CLI tool that lets you instantly share your local development environment with others on your local network. It detects your project type (React, Next.js, Node.js, Go, Laravel), starts the appropriate server, and displays a QR code for easy access from any device on your LAN.

---

## 🚀 Features

- 🔍 **Auto-detects** project type (React, Next.js, Node.js, Go, Laravel)
- 🛠️ **Runs** the correct dev server for your project
- 🌍 **Custom port configuration** simple and flexible
- 🌐 **LAN sharing**: exposes your app to your local network
- 📱 **QR code**: scan to open your app instantly on any device
- 🔐 **Password protection**: secure your development environment with authentication
- 🧠 **Smart dependency install**: auto-installs dependencies if needed
- 💡 **Simple CLI**: just run `devshare` in your project folder

---

## 📦 Supported Frameworks

- [x] React (Vite, Create React App)
- [x] Next.js
- [x] Node.js (with `start` or `dev` script, or entry file)
- [x] Go
- [x] Laravel

---

## 🛠️ Installation

### 1. Download and Install DevShare

#### Option 1: Download from Releases (Recommended)
Download the latest release for your platform from [GitHub Releases](https://github.com/abdorrahmani/devshare/releases) and extract the archive:

- **Windows**: `devshare_windows_amd64.zip`
- **macOS**: `devshare_darwin_amd64.tar.gz` or `devshare_darwin_arm64.tar.gz` (Apple Silicon)
- **Linux**: `devshare_linux_amd64.tar.gz`

Each archive contains:
- The DevShare binary (`devshare` or `devshare.exe`)
- The installer script (`install.sh` for Linux/macOS, `install.ps1` for Windows)

**To install:**
- On **Linux/macOS**:
  ```sh
  cd <extracted-folder>
  sudo bash install.sh
  ```
- On **Windows**:
  - **Easiest:** Double-click `install.bat` in the extracted folder. This will request Administrator privileges, then run the installer in a new window.
  - **Advanced:** If you prefer to use Windows Terminal or PowerShell, first right-click Windows Terminal and select "Run as administrator", then:
    ```powershell
    cd <extracted-folder>
    .\install.ps1
    ```

- To check installation success:
    ```sh
    devshare --version
    DevShare version: vX.X.X
    ```

The installer will copy the binary to the appropriate location and add it to your PATH if needed. No build step is required.

**Troubleshooting:**
- If you see `devshare.exe not found in current directory`, make sure you are running the installer from the same folder where `devshare.exe` is located (the extracted archive folder). If you double-clicked `install.bat`, this should happen automatically.

#### Option 2: Install via Go (for Go users)
```sh
go install github.com/abdorrahmani/devshare@latest
```

#### Option 3: Build from source
```sh
git clone https://github.com/abdorrahmani/devshare.git
cd devshare
go mod download
go build -o devshare
```

---

## 🚦 Usage

1. Open your terminal in your project directory.
2. Run:

```sh
# Basic usage (no authentication)
devshare

# With authentication (password protection)
devshare --password=your_password

# Specify custom port
devshare 3000 --password=your_password
```

1. DevShare will:
   - Detect your project type
   - Start the appropriate dev server
   - Show your LAN IP and a QR code
   - If password is provided, create a secure proxy with authentication

2. Scan the QR code with your phone or share the link with others on your LAN!

### 🔐 Authentication

When you use the `--password` flag, DevShare creates a secure proxy server that:
- Requires a password to access your development environment
- Provides a beautiful authentication page
- Maintains session cookies for 1 hour
- Protects your development work from unauthorized access

**Example with authentication:**
```sh
devshare --password=mysecret123
```

This will show:
- Local access: `http://localhost:3000` (direct access)
- Secure network access: `http://192.168.1.42:3001` (password protected)
- QR code for the secure access URL

---

## 📝 Example Output

### Without Authentication
```
 ____              _____ _                    
|  _ \  _____   __/ ___|| |__   __ _ _ __ ___ 
| | | |/ _ \ \ / /\___ \| '_ \ / _ | '__/ _ \
| |_| |  __/\ V /  ___) | | | | (_| | | |  __/
|____/ \___| \_/  |____/|_| |_|\__,_|_|  \___|

DevShare - Share your dev environment over LAN
🔗 Served by DevShare - Built with 💙 by Anophel

Current working directory: /your/project/path
✅ React project detected!
📦 Using package manager: npm
🚀 Starting React app...
🌐 Your LAN IP is: 192.168.1.42
Local:   http://localhost:5173
Network: http://192.168.1.42:5173
📱 Scan this on your phone (React/Vite default port 5173):
[QR CODE]
```

### With Authentication
```
 ____              _____ _                    
|  _ \  _____   __/ ___|| |__   __ _ _ __ ___ 
| | | |/ _ \ \ / /\___ \| '_ \ / _ | '__/ _ \
| |_| |  __/\ V /  ___) | | | | (_| | | |  __/
|____/ \___| \_/  |____/|_| |_|\__,_|_|  \___|

DevShare - Share your dev environment over LAN
🔗 Served by DevShare - Built with 💙 by Anophel

Current working directory: /your/project/path
✅ React project detected!
📦 Using package manager: npm
🔐 Authentication enabled - Password required to access the app
🚀 Starting React app...
🌐 Your LAN IP is: 192.168.1.42
Local:   http://localhost:5173
🔗 Auth Proxy: http://192.168.1.42:5174
📱 Scan this on your phone (React/Vite default port 5173):
[QR CODE]
```

---

## 🔄 Updating DevShare

To update DevShare to the latest version, simply run:

```sh
devshare update
```

This command will:
- Check for the latest release on GitHub
- Download the correct archive for your platform
- Extract the new binary and run the official installer script (`install.sh` on Linux/macOS, `install.bat` on Windows)
- Replace your current DevShare installation with the latest version

**Note:**
- You may be prompted for administrator/root permissions during the update, depending on your system configuration.
- If the update fails, you can always manually download and install the latest release as described in the installation section above.

---

## Security Note for Node.js Users

**WARNING:** Your Node.js app may be listening on all interfaces (`0.0.0.0`). For security, ensure your app binds to `127.0.0.1` to prevent bypassing authentication.

If you control the app, update your server code to listen only on `127.0.0.1`, for example:

```js
// Node.js (Express example):
const host = process.env.HOST || '127.0.0.1';
const port = process.env.PORT || 3000;
app.listen(port, host, () => {
  console.log(`Server running at http://${host}:${port}/`);
});
```

---

## 🤝 Contributing

Contributions are welcome! Please open issues or pull requests on [GitHub](https://github.com/abdorrahmani/devshare).

---

## 📄 License

This project is licensed under the [MIT License](https://github.com/abdorrahmani/devshare/blob/master/LICENSE).

---

## 🔗 Links

- [GitHub Repository](https://github.com/abdorrahmani/devshare)
- [MIT License](https://github.com/abdorrahmani/devshare/blob/master/LICENSE)

---

> Built with 💙 by [Anophel](https://github.com/abdorrahmani)
