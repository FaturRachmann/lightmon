# ⌨️ lightmon Keyboard Shortcuts

## 📊 Monitoring

| Key | Action |
|-----|--------|
| `1` | Switch to **Overview** tab |
| `2` | Switch to **Processes** tab |
| `3` | Switch to **Help** tab |
| `e` | Toggle **per-core CPU** view |
| `b` | Toggle **battery** panel |
| `l` | Toggle **system load** (uptime, load avg) |

## 🔧 Process Management

| Key | Action |
|-----|--------|
| `k` | **Kill** process (opens signal menu) |
| `i` | Show process **Info/Details** |
| `+` | **Renice** process (change priority) |
| `s` | **Send Signal** menu (SIGUSR1, SIGHUP, etc) |

## 🔍 Search & Sort

| Key | Action |
|-----|--------|
| `/` | **Filter** process by name |
| `c` | Sort by **CPU%** (descending) |
| `m` | Sort by **Memory%** (descending) |
| `p` | Sort by **PID** (ascending) |
| `n` | Sort by **Name** (alphabetical) |

## ⚠️ General

| Key | Action |
|-----|--------|
| `h` | Toggle **help** bar |
| `q` | **Quit** application |
| `↑` `↓` | Navigate process list |
| `Enter` | Confirm selection |
| `Esc` | Cancel/Go back |

---

## 📋 Kill Signal Options

When you press `k` to kill a process, you can choose:

| Signal | Number | Description |
|--------|--------|-------------|
| **SIGTERM** | 15 | Graceful termination (recommended) |
| **SIGKILL** | 9 | Force kill (use with caution) |
| **SIGINT** | 2 | Interrupt (like Ctrl+C) |
| **SIGHUP** | 1 | Hangup (reload config) |
| **SIGQUIT** | 3 | Quit with core dump |
| **SIGUSR1** | - | User-defined signal 1 |
| **SIGUSR2** | - | User-defined signal 2 |
| **SIGSTOP** | - | Pause process |
| **SIGCONT** | - | Resume process |

---

## 🎯 Quick Start

```bash
./lightmon
```

Then press:
- `h` - See full help bar at bottom
- `/` - Type process name to filter
- `k` - Kill with signal selection
- `i` - View process details
- `q` - Quit anytime

---

## 💡 Tips

1. **Filter Mode**: Press `/`, type process name, press `Enter` to confirm or `Esc` to cancel
2. **Kill Menu**: Interactive menu with all signal options
3. **Renice**: Change process priority from -20 (highest) to 19 (lowest)
4. **Quit Anytime**: `q` works even in filter mode!

---

<div align="center">

**Press `h` in the application to see this guide anytime!**

</div>
