# ⚡ lightmon

<div align="center">

```
██╗     ██╗ ██████╗ ██╗  ██╗████████╗███╗   ███╗ ██████╗ ███╗   ██╗
██║     ██║██╔════╝ ██║  ██║╚══██╔══╝████╗ ████║██╔═══██╗████╗  ██║
██║     ██║██║  ███╗███████║   ██║   ██╔████╔██║██║   ██║██╔██╗ ██║
██║     ██║██║   ██║██╔══██║   ██║   ██║╚██╔╝██║██║   ██║██║╚██╗██║
███████╗██║╚██████╔╝██║  ██║   ██║   ██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
```

**Lightweight. Advanced. Terminal-native system monitor.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](#license)
[![Binary Size](https://img.shields.io/badge/Binary-~8MB-blue?style=flat-square)](#)

</div>

---

## Mengapa lightmon?

Kebanyakan system monitor modern seperti `htop`, `btop`, dan `glances` sangat powerful — tapi juga berat, lambat startup, dan penuh fitur yang jarang dipakai. **lightmon** mengambil pendekatan berbeda: ambil hanya yang penting, eksekusi dengan sempurna, dan jaga agar binary tetap kecil.

> *"Fast by default. Advanced by design."*

---

## Preview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  CPU                    Memory               Disk               Network     │
├──────────────────┬──────────────────┬────────────────────┬──────────────────┤
│ Total            │ RAM              │ /                  │ eth0             │
│ [████████░░] 63% │ [███████░░] 58%  │ [████░░░░░░] 41%   │ ▲  1.2 MB/s     │
│                  │ 12.4 GB / 16 GB  │ 136 GB / 256 GB    │ ▼  840 KB/s     │
│ [e] per-core     │                  │                    │                  │
│                  │ Swap             │ /home              │ Tx: 4.2 GB      │
│                  │ [░░░░░░░░░░]  2% │ [██░░░░░░░░] 18%   │ Rx: 12.1 GB     │
│                  │ 0.4 GB / 16 GB   │ 89 GB / 512 GB     │                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  Processes                                                                  │
│  PID     NAME                   STATUS    CPU%      MEM%                   │
│ ▶ 1234   firefox                S         45.2%     11.8%                  │
│   5678   code                   S         22.1%      8.4%                  │
│   9012   node                   R         15.0%      3.2%                  │
│   3456   postgres                S          4.3%      6.1%                  │
│   7890   python3                R          2.8%      1.9%                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ Updated: 14:32:01  Sort: CPU  Procs: 142/142                               │
│ [c] CPU  [m] MEM  [p] PID  [n] NAME  [k] kill  [e] cores  [q] quit        │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Color coding:**
- 🟢 `< 70%` — Normal
- 🟡 `70–89%` — Perhatian
- 🔴 `≥ 90%` — Kritis

---

## Fitur

### 🖥️ CPU Monitor
- Usage total sistem secara real-time
- **Toggle per-core** — lihat beban tiap core secara individual (tekan `e`)
- Deteksi jumlah core otomatis via `gopsutil`
- Progress bar dengan warna adaptif berdasarkan threshold

### 🧠 Memory Monitor
- RAM: used / total dengan persentase
- **Swap memory** — monitoring terpisah, muncul otomatis jika swap aktif
- Format bytes cerdas: otomatis tampil dalam KB / MB / GB sesuai ukuran

### 💾 Disk Monitor
- Deteksi **semua partisi** yang ter-mount secara otomatis
- Deduplikasi per device (tidak ada duplikat entry untuk partisi yang sama)
- Tampil path mount point, used, free, dan persentase pemakaian
- Skip partisi virtual (total = 0)

### 🌐 Network Monitor
- **TX / RX rate** per detik (bytes/sec) — dihitung dari delta antar polling
- Akumulasi total bytes sent dan received sejak boot
- Support multiple interface (eth0, wlan0, dst.)
- Loopback (`lo`) di-skip otomatis

### ⚙️ Process Manager
- Tampilkan top N proses (default 20, configurable via `--procs`)
- **4 mode sorting** — langsung ganti tanpa restart:
  - `c` → sort by CPU usage
  - `m` → sort by Memory usage
  - `p` → sort by PID
  - `n` → sort by Name (alphabetical)
- Status proses color-coded: Running (hijau), Sleep (cyan), Zombie (merah)
- **Kill process** — pilih baris, tekan `k`, konfirmasi → kirim SIGTERM

### 🎨 TUI & UX
- Layout 4-panel atas + process table + status bar
- Navigasi keyboard penuh — tidak perlu mouse
- Refresh rate configurable (`--interval 500ms` sampai beberapa detik)
- Non-blocking render — UI tidak freeze saat pengumpulan data
- Padding, alignment, dan truncation rapi di semua ukuran terminal

---

## Quick Start

```bash
# Clone
git clone https://github.com/youruser/lightmon
cd lightmon

# Install dependencies & build
make install

# Jalankan
lightmon
```

Atau tanpa install:

```bash
go mod tidy
go run main.go
```

---

## Instalasi

### Dari Source (Rekomendasi)

```bash
git clone https://github.com/youruser/lightmon
cd lightmon
go mod tidy
make build
```

Binary siap di `./lightmon`.

### Install ke System PATH

```bash
make install
# Binary tersalin ke /usr/local/bin/lightmon
```

### Cross-compile

```bash
make build-all

# Menghasilkan:
# lightmon-linux-amd64
# lightmon-darwin-arm64    (Apple Silicon)
# lightmon-darwin-amd64    (Intel Mac)
```

### Manual (tanpa Makefile)

```bash
go mod tidy
go build -ldflags="-s -w" -o lightmon .
```

---

## Penggunaan

```
lightmon [flags]
```

| Flag | Default | Deskripsi |
|------|---------|-----------|
| `--interval` | `1s` | Interval refresh. Contoh: `500ms`, `2s` |
| `--procs` | `20` | Jumlah maksimal proses yang ditampilkan |
| `--version` | — | Tampilkan versi dan keluar |

**Contoh:**

```bash
lightmon                          # Default — 1 detik, 20 proses
lightmon --interval 500ms         # Refresh lebih cepat
lightmon --procs 50               # Tampilkan top 50 proses
lightmon --interval 2s --procs 10 # Refresh lambat, proses sedikit
lightmon --version                # lightmon v1.0.0
```

---

## Keyboard Shortcuts

| Tombol | Aksi |
|--------|------|
| `c` | Sort proses berdasarkan CPU% |
| `m` | Sort proses berdasarkan Memory% |
| `p` | Sort proses berdasarkan PID |
| `n` | Sort proses berdasarkan Nama |
| `k` | Kill proses yang dipilih (SIGTERM + konfirmasi) |
| `e` | Toggle tampilan per-core CPU |
| `r` | Force refresh |
| `↑` / `↓` | Navigasi list proses |
| `q` | Keluar |

---

## Arsitektur

```
lightmon/
├── main.go                 # Entry point, parsing CLI flags
│
├── metrics/                # Layer pengumpulan data (murni data, no UI)
│   ├── cpu.go              # CPU total + per-core via gopsutil
│   ├── memory.go           # RAM + Swap + helper FormatBytes()
│   ├── disk.go             # Semua partisi + deduplikasi per device
│   ├── network.go          # I/O counters + kalkulasi rate delta
│   └── process.go          # Daftar proses + multi-mode sorting
│
├── internal/
│   └── updater.go          # Goroutine polling + SystemSnapshot channel
│
└── ui/
    └── dashboard.go        # Layout TUI, render, key bindings (tview)
```

### Alur Data

```
┌──────────────┐    ticker     ┌──────────────┐    channel    ┌──────────────┐
│   metrics/*  │ ────────────► │   Updater    │ ────────────► │  Dashboard   │
│              │               │              │               │              │
│ GetCPU()     │               │ collect()    │               │ render()     │
│ GetMemory()  │               │              │   snapshot    │ renderCPU()  │
│ GetDisks()   │               │ SystemSnapshot ──────────────► renderMem()  │
│ GetNetwork() │               │              │               │ renderDisk() │
│ GetProcesses │               │ stop chan    │               │ renderNet()  │
└──────────────┘               └──────────────┘               └──────────────┘
```

- **metrics/** — stateless, pure functions, tidak ada side effect UI
- **Updater** — satu goroutine, polling dengan `time.Ticker`, push ke buffered channel (size 1)
- **Dashboard** — consume channel via `app.QueueUpdateDraw()`, thread-safe dengan tview

### Design Decisions

**Mengapa buffered channel size 1?**
Jika render lebih lambat dari polling (terminal lambat), snapshot lama dibuang dan diganti snapshot terbaru — tidak ada backpressure, tidak ada memory leak.

**Mengapa `gopsutil`?**
Single dependency untuk semua OS. Abstraksi `/proc/` (Linux) dan `sysctl` (macOS) secara transparan — kode metrics tidak perlu tahu platform.

**Mengapa `tview`?**
Built on `tcell`, support keyboard navigation, table selection, dan modal out-of-the-box. Lebih high-level dari raw ANSI escape codes tapi jauh lebih ringan dari Bubble Tea ecosystem.

---

## Dependencies

| Package | Versi | Fungsi |
|---------|-------|--------|
| [`gopsutil/v3`](https://github.com/shirou/gopsutil) | v3.23.x | CPU, Memory, Disk, Network, Process metrics |
| [`tview`](https://github.com/rivo/tview) | latest | TUI framework (panels, tables, modals) |
| [`tcell/v2`](https://github.com/gdamore/tcell) | v2.7.x | Terminal colors, key events (dep of tview) |

**Total binary size (setelah strip):** ~8 MB

---

## Perbandingan

| | lightmon | htop | btop | glances |
|---|---|---|---|---|
| Bahasa | Go | C | C++ | Python |
| Binary size | ~8 MB | ~500 KB | ~4 MB | N/A (Python) |
| Startup time | < 100ms | < 50ms | ~200ms | ~1–2 detik |
| CPU panel | ✅ total + per-core | ✅ | ✅ | ✅ |
| Memory + Swap | ✅ | ✅ | ✅ | ✅ |
| Disk per-partisi | ✅ auto-detect | ❌ | ✅ | ✅ |
| Network rate | ✅ TX/RX/s | ❌ | ✅ | ✅ |
| Kill process | ✅ | ✅ | ✅ | ✅ |
| Sort proses | ✅ 4 mode | ✅ | ✅ | ✅ |
| Cross-compile | ✅ single binary | ❌ | ❌ | ❌ |
| Zero runtime deps | ✅ | ✅ | ✅ | ❌ |
| Config file | ❌ YAGNI | ✅ | ✅ | ✅ |

---

## Roadmap

### v1.0 — Selesai ✅
- [x] CPU, RAM, Swap, Disk, Network, Process
- [x] Color-coded thresholds
- [x] Sort proses (CPU/MEM/PID/NAME)
- [x] Kill process dengan konfirmasi
- [x] Per-core CPU toggle
- [x] Network TX/RX rate per detik
- [x] CLI flags: `--interval`, `--procs`, `--version`
- [x] Cross-compile Linux + macOS

### v1.1 — Planned
- [ ] CPU history sparkline (mini grafik dalam panel)
- [ ] Filter proses by name (`/` untuk search)
- [ ] Temperature sensor (CPU/GPU jika tersedia)

### v1.2 — Future
- [ ] Log snapshot ke CSV untuk analisis historis
- [ ] Config file `~/.config/lightmon/config.toml`
- [ ] Docker container awareness (nama container sebagai process name)

---

## Kontribusi

Pull request welcome. Untuk perubahan besar, buka issue dulu.

```bash
# Setup dev
git clone https://github.com/youruser/lightmon
cd lightmon
go mod tidy

# Run
go run main.go --interval 500ms

# Test
go vet ./...
go build ./...
```

---

## License

MIT — bebas digunakan, dimodifikasi, dan didistribusikan.

---

<div align="center">

Built with Go · Powered by gopsutil · TUI by tview

*"The best tool is the one that gets out of your way."*

</div>
