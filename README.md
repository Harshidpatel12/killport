# 🚀 killport

[![Go Version](https://img.shields.io/github/go-mod/go-version/Harshidpatel12/killport?logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/Harshidpatel12/killport/actions/workflows/test.yml/badge.svg)](https://github.com/Harshidpatel12/killport/actions/workflows/test.yml)
[![Lint Status](https://github.com/Harshidpatel12/killport/actions/workflows/lint.yml/badge.svg)](https://github.com/Harshidpatel12/killport/actions/workflows/lint.yml)

A lightweight, blazing-fast CLI utility written in Go to quickly find and terminate processes running on specific network ports. 

Never write complex `lsof -i :3000` and `kill -9 <PID>` chains again.

---

## ✨ Features

- ⚡ **Zero Dependencies**: Compiles to a single static binary with no external runtimes needed.
- 🔍 **Reliable Lookup**: Uses system-level port diagnostics to locate active process PIDs.
- 🛡️ **Safety Guard**: Automatically prevents self-termination if the tool is somehow matching its own process.
- 🎨 **Clear Console Logs**: Reports how many processes were found, their PIDs, and their termination status.

---

## 🚀 Installation & Usage

Since `killport` is compiled into a single static binary with no external dependencies, **you do not need Go installed** on your machine to run it.

### 1. Installation

#### Option A: One-Line Shell Installer (Recommended for Linux/macOS)
Run the following command to automatically download the correct pre-compiled binary for your operating system and architecture, and install it to `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/Harshidpatel12/killport/main/install.sh | bash
```

*(Note: The install script will become active once the automated release pipeline is configured).*

#### Option B: Manual Binary Download (Zero Dependencies)
1. Navigate to the [Releases](https://github.com/Harshidpatel12/killport/releases) page.
2. Download the compressed archive matching your operating system and CPU architecture (e.g., `killport-linux-amd64.tar.gz`).
3. Extract the archive and copy the binary to your system PATH:
   ```bash
   tar -xvf killport-linux-amd64.tar.gz
   chmod +x killport
   sudo mv killport /usr/local/bin/
   ```

#### Option C: Via Package Managers (Coming Soon)
*   **macOS (Homebrew):**
    ```bash
    brew install Harshidpatel12/tap/killport
    ```
*   **Ubuntu/Debian (APT):**
    ```bash
    sudo apt install killport
    ```

#### Option D: Via Go (For developers with Go installed)
If you have Go installed on your machine, you can download, compile, and install it directly from source:
```bash
go install github.com/Harshidpatel12/killport@latest
```

---

### 2. Usage

To kill any process running on port `3000`:

```bash
killport 3000
```

#### Example Output:
```text
Searching for processes on port 3000...
Found 1 process(es) on port 3000: 12345
Killing process with PID 12345...
Successfully killed process 12345.
```


---

## 🛠️ Development & Linting

This project uses standard Go tools and `pre-commit` hooks to maintain code style and quality.

### Local Setup
1. Clone the project.
2. Initialize pre-commit hooks (requires `pre-commit` installed):
   ```bash
   pre-commit install
   ```

### Run Tests Locally
```bash
go test -v ./...
```

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
