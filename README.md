# ⚡ LightMon

<div align="center">

**Lightweight Terminal System Monitor**

Fast • Advanced • User-Friendly

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](#license)

</div>

---

## 📖 Tentang LightMon

LightMon adalah system monitor terminal yang **ringan, cepat, dan kaya fitur**. Terinspirasi oleh tools seperti `htop`, `btop`, dan `glances`, LightMon mengambil pendekatan berbeda: **ambil yang penting, eksekusi dengan sempurna, dan jaga agar tetap sederhana**.

### ✨ Keunggulan

- 🚀 **Fast Startup** - Mulai dalam < 100ms
- 💾 **Lightweight** - Binary ~8MB, zero runtime dependencies
- 🎨 **Modern TUI** - UI yang clean dan intuitif
- ⚡ **Real-time** - Monitoring real-time dengan update cepat
- 🎯 **User-Friendly** - Onboarding interaktif untuk user baru
- 🔧 **Advanced Features** - Kill, renice, signals, filter, dan lebih banyak lagi

---

## 📸 Preview

```
╭──────────────────────────────────────────────────────────────╮
│  [1] Overview  [2] Processes  [3] Help                       │
├──────────────────────────────────────────────────────────────┤
│  CPU                          │  Memory                      │
│    63.5%  [████████░░]        │    58.1%  [███████░░]        │
│    Temp: 72°C                 │    12.4G / 16.0G             │
│    ██ ██ ██ ██                │    Swap: [██░░░░] 2.0%       │
├───────────────────────────────┴──────────────────────────────┤
│  Disk                        │  Network                      │
│    /      [████░░] 41%       │    eth0: ↑1.2M/s ↓840K/s     │
│    /home  [██░░░░] 18%       │    wlan0: ↑500K/s ↓200K/s    │
├──────────────────────────────────────────────────────────────┤
│  PID       NAME                      CPU%      MEM%   STATUS │
│  1234      firefox                   45.2      11.8     ●    │
│  5678      code                      22.1       8.4     ○    │
│  ...                                                         │
├──────────────┬──────────────┬────────────────────────────────┤
│  System      │  Quick Keys  │  Help & Tips                   │
│    Uptime:   │    e cores   │    MONITORING:                 │
│    2h 45m    │    c CPU     │      1  Overview               │
│    Load:     │    m MEM     │      2  Processes              │
│    1.2 0.8   │    k kill    │    PROCESS MGMT:               │
│    Procs:    │    i info    │      k  Kill                   │
│    142       │    / filter  │      i  Info                   │
└──────────────┴──────────────┴────────────────────────────────┘
  15:04:05  │  Sort: CPU  │  20 procs  │  h:help  q:quit  /:filter
```

---

## 🚀 Fitur

### 📊 Monitoring

| Fitur | Deskripsi |
|-------|-----------|
| **CPU** | Total usage %, per-core view, temperature monitoring |
| **Memory** | RAM usage, Swap usage, used/total display |
| **Disk** | All mounted partitions with usage % |
| **Network** | Per-interface TX/RX rates + totals |
| **System Info** | Uptime, load averages, process count, boot time |
| **Battery** | Battery percentage, charging status, time remaining |

### 🔧 Process Management

| Fitur | Deskripsi |
|-------|-----------|
| **Sort** | CPU%, Memory%, PID, Name (toggle dengan c/m/p/n) |
| **Kill** | Interactive signal selection (SIGTERM, SIGKILL, dll) |
| **Info** | View detailed process information |
| **Renice** | Change process priority (-20 to 19) |
| **Signals** | Send custom signals (SIGUSR1, SIGUSR2, SIGHUP, dll) |
| **Filter** | Search/filter processes by name (/) |

### 🎨 User Experience

| Fitur | Deskripsi |
|-------|-----------|
| **Welcome Guide** | Interactive onboarding untuk first-time users |
| **Help Panel** | Always-visible shortcuts reference |
| **Quick Keys** | Dedicated panel untuk common shortcuts |
| **Modal Dialogs** | Clean, interactive dialogs untuk actions |
| **Color Coding** | 🟢 Green (<70%), 🟡 Yellow (70-89%), 🔴 Red (≥90%) |
| **Sparklines** | CPU/Memory history visualization |

---

## 📦 Instalasi

### Prerequisites

- **Go 1.21+** (untuk build dari source)
- **Linux** atau **macOS**

### Metode 1: Build dari Source (Recommended)

```bash
# Clone repository
git clone https://github.com/youruser/lightmon
cd lightmon

# Download dependencies
go mod tidy

# Build
make build

# Install ke system PATH (optional)
sudo make install
```

Binary akan tersedia di:
- `./lightmon` (setelah build)
- `/usr/local/bin/lightmon` (setelah install)

### Metode 2: Download Pre-built Binary

Download binary untuk platform Anda dari [releases page](https://github.com/youruser/lightmon/releases):

```bash
# Linux AMD64
wget https://github.com/youruser/lightmon/releases/download/v2.0.0/lightmon-linux-amd64
chmod +x lightmon-linux-amd64
sudo mv lightmon-linux-amd64 /usr/local/bin/lightmon

# macOS ARM64 (Apple Silicon)
wget https://github.com/youruser/lightmon/releases/download/v2.0.0/lightmon-darwin-arm64
chmod +x lightmon-darwin-arm64
sudo mv lightmon-darwin-arm64 /usr/local/bin/lightmon

# macOS AMD64 (Intel)
wget https://github.com/youruser/lightmon/releases/download/v2.0.0/lightmon-darwin-amd64
chmod +x lightmon-darwin-amd64
sudo mv lightmon-darwin-amd64 /usr/local/bin/lightmon
```

### Metode 3: Go Install

```bash
go install github.com/youruser/lightmon@latest
```

Binary akan tersedia di `$GOPATH/bin/lightmon`.

---

## 🎮 Penggunaan

### Basic Usage

```bash
# Jalankan dengan default settings
lightmon

# Refresh interval lebih cepat
lightmon --interval 500ms

# Tampilkan lebih banyak processes
lightmon --procs 50

# Export metrics ke file JSON
lightmon --export metrics.json

# Export ke CSV
lightmon --export metrics.csv --export-format csv
```

### Command Line Flags

| Flag | Default | Deskripsi |
|------|---------|-----------|
| `--interval` | `1s` | Refresh interval (e.g., `500ms`, `1s`, `2s`) |
| `--procs` | `20` | Max processes to display |
| `--config` | `~/.lightmon/config.yaml` | Path to config file |
| `--export` | — | Export metrics to file (JSON/CSV) |
| `--export-format` | `json` | Export format: `json` or `csv` |
| `--no-export` | `false` | Disable metrics export |
| `--version` | — | Show version and exit |

### Examples

```bash
# Default monitoring
lightmon

# Fast refresh (500ms)
lightmon --interval 500ms

# Monitor dengan export metrics
lightmon --export /var/log/lightmon.json --interval 1s

# Custom config file
lightmon --config /etc/lightmon/config.yaml

# Show top 50 processes
lightmon --procs 50
```

---

## ⌨️ Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `↑` `↓` | Navigate process list |
| `1` | Switch to **Overview** tab |
| `2` | Switch to **Processes** tab |
| `3` | Switch to **Help** tab |
| `Enter` | Confirm selection |
| `Esc` | Cancel operation |

### Monitoring

| Key | Action |
|-----|--------|
| `e` | Toggle CPU per-core view |
| `b` | Toggle battery panel |
| `l` | Toggle system load info |

### Process Management

| Key | Action |
|-----|--------|
| `k` | **Kill** process (interactive signal menu) |
| `i` | View process **Info/Details** |
| `+` | **Renice** process (change priority) |
| `s` | Send **Signal** (SIGUSR1, SIGUSR2, dll) |

### Sorting

| Key | Action |
|-----|--------|
| `c` | Sort by **CPU%** (descending) |
| `m` | Sort by **Memory%** (descending) |
| `p` | Sort by **PID** (ascending) |
| `n` | Sort by **Name** (alphabetical) |

### Search & General

| Key | Action |
|-----|--------|
| `/` | **Filter** process by name |
| `h` | Toggle **Help** panel |
| `q` | **Quit** application |

---

## 🎯 Kill Signals

Saat menekan `k` untuk kill process, tersedia signal berikut:

| Signal | Number | Deskripsi |
|--------|--------|-----------|
| **SIGTERM** | 15 | Graceful termination (default, recommended) |
| **SIGKILL** | 9 | Force kill (use with caution) |
| **SIGINT** | 2 | Interrupt (seperti Ctrl+C) |
| **SIGHUP** | 1 | Hangup (reload config) |
| **SIGQUIT** | 3 | Quit dengan core dump |
| **SIGUSR1** | - | User-defined signal 1 |
| **SIGUSR2** | - | User-defined signal 2 |
| **SIGSTOP** | - | Pause process |
| **SIGCONT** | - | Resume process |

---

## 📁 Struktur Project

```
lightmon/
├── main.go                 # Entry point, CLI flags
├── Makefile                # Build automation
├── README.md               # Documentation
├── go.mod                  # Go module definition
├── metrics/
│   ├── cpu.go              # CPU metrics
│   ├── memory.go           # Memory + Swap metrics
│   ├── disk.go             # Disk partition metrics
│   ├── network.go          # Network I/O metrics
│   ├── process.go          # Process list + sorting
│   └── temperature.go      # CPU temperature
├── internal/
│   └── updater.go          # Metrics polling + history + alerts
├── ui/
│   ├── dashboard.go        # Main TUI layout
│   └── menu.go             # Interactive menus (Kill, Info, Renice)
├── config/
│   └── config.go           # YAML configuration
├── history/
│   └── history.go          # Metrics history + sparklines
├── export/
│   └── exporter.go         # CSV/JSON export
├── battery/
│   └── battery.go          # Battery monitoring
└── system/
    └── system.go           # System info (uptime, load avg)
```

---

## ⚙️ Configuration

LightMon menggunakan konfigurasi YAML di `~/.lightmon/config.yaml`:

```yaml
# Alert thresholds
thresholds_config:
  cpu_warning: 70
  cpu_critical: 90
  mem_warning: 75
  mem_critical: 90
  disk_warning: 80
  disk_critical: 95

# Display options
display:
  show_cores: false
  show_battery: true
  show_load_avg: true
  process_limit: 20
  refresh_rate: "1s"

# Metrics logging
logging:
  enabled: false
  file_path: "~/.lightmon/metrics.log"
  format: "json"
```

Generate default config:

```bash
# Config akan dibuat otomatis saat first run
# Atau manual:
mkdir -p ~/.lightmon
cp config.example.yaml ~/.lightmon/config.yaml
```

---

## 🔧 Development

### Requirements

- Go 1.21+
- Make (optional, untuk build automation)

### Build

```bash
# Download dependencies
go mod tidy

# Build untuk current platform
go build -o lightmon .

# Build dengan version info
go build -ldflags="-s -w -X main.version=2.0.0" -o lightmon .

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o lightmon-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o lightmon-darwin-arm64 .
```

### Run from Source

```bash
go run main.go --interval 500ms
```

### Testing

```bash
# Run tests
go test ./...

# Check code quality
go vet ./...
go fmt ./...
```

---

## 🗺️ Roadmap

### ✅ Completed (v2.0)

- [x] CPU, Memory, Disk, Network monitoring
- [x] Process management (Kill, Info, Renice, Signals)
- [x] Interactive menus dengan modal dialogs
- [x] Process filter/search
- [x] CPU/Memory history sparklines
- [x] Temperature monitoring
- [x] Battery monitoring
- [x] System info (uptime, load avg)
- [x] CSV/JSON export
- [x] YAML configuration
- [x] Welcome guide untuk first-time users
- [x] Help panel dengan shortcuts

### 🚧 Planned

- [ ] Process tree view
- [ ] Docker container monitoring
- [ ] Custom dashboard layouts
- [ ] Remote monitoring (client-server mode)
- [ ] Plugin system
- [ ] Web UI (optional)

---

## 🤝 Contributing

Contributions are welcome! Silakan:

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Code Style

- Follow Go best practices
- Run `go fmt` dan `go vet` sebelum commit
- Write meaningful commit messages
- Add comments untuk complex logic

---

## 📄 License

Distributed under the **MIT License**. Lihat `LICENSE` untuk lebih detail.

---

## 🙏 Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) - System metrics library
- [tview](https://github.com/rivo/tview) - Terminal UI library
- [tcell](https://github.com/gdamore/tcell) - Terminal cell library
- Inspired by [htop](https://hisham.hm/htop/), [btop](https://github.com/aristocratos/btop), [glances](https://nicolargo.github.io/glances/)

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/youruser/lightmon/issues)
- **Discussions:** [GitHub Discussions](https://github.com/youruser/lightmon/discussions)
- **Email:** your.email@example.com

---

<div align="center">

**Made with ❤️ using Go**

*"The best tool is the one that gets out of your way."*

⭐ Star this repo if you find it useful!

</div>
