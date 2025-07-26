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
- 🌐 **LAN sharing**: exposes your app to your local network
- 📱 **QR code**: scan to open your app instantly on any device
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

### 1. Install Go (if not already installed)
- [Download Go](https://golang.org/dl/)

### 2. Install DevShare

#### Option 1: Download from releases
Download the latest release for your platform from [GitHub Releases](https://github.com/abdorrahmani/devshare/releases):

- **Windows**: `devshare_windows_amd64.exe`
- **macOS**: `devshare_darwin_amd64` or `devshare_darwin_arm64` (Apple Silicon)
- **Linux**: `devshare_linux_amd64`

#### Option 2: Install via Go
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
devshare
```

1. DevShare will:
   - Detect your project type
   - Start the appropriate dev server
   - Show your LAN IP and a QR code

2. Scan the QR code with your phone or share the link with others on your LAN!

---

## 📝 Example Output

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
