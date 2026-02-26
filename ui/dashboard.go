package ui

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/FaturRachmann/lightmon/config"
	"github.com/FaturRachmann/lightmon/export"
	"github.com/FaturRachmann/lightmon/internal"
	"github.com/FaturRachmann/lightmon/metrics"
	"github.com/rivo/tview"
)

const barWidth = 40

type Dashboard struct {
	app           *tview.Application
	updater       *internal.Updater
	config        *config.Config
	exporter      *export.Exporter
	sortBy        metrics.SortBy
	showCores     bool
	processFilter string
	filterMode    bool
	selectedRow   int
	firstRun      bool

	cpuView      *tview.TextView
	memView      *tview.TextView
	diskView     *tview.TextView
	netView      *tview.TextView
	sysInfoView  *tview.TextView
	procTable    *tview.Table
	statusBar    *tview.TextView
	shortcutView *tview.TextView
	helpView     *tview.TextView
	rootLayout   tview.Primitive
}

func NewDashboard(updater *internal.Updater, cfg *config.Config, exporter *export.Exporter) *Dashboard {
	d := &Dashboard{
		app:       tview.NewApplication(),
		updater:   updater,
		config:    cfg,
		exporter:  exporter,
		sortBy:    metrics.SortByCPU,
		showCores: cfg.Display.ShowCores,
		firstRun:  true,
	}
	d.build()
	return d
}

func (d *Dashboard) build() {
	borderColor := tcell.ColorDarkGray

	// TOP: CPU & Memory
	d.cpuView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.cpuView.SetBorder(true).SetTitle(" CPU ").SetTitleColor(tcell.ColorGreen).SetBorderColor(borderColor)

	d.memView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.memView.SetBorder(true).SetTitle(" Memory ").SetTitleColor(tcell.ColorGreen).SetBorderColor(borderColor)

	// SECONDARY: Disk & Network
	d.diskView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.diskView.SetBorder(true).SetTitle(" Disk ").SetTitleColor(tcell.ColorDarkGray).SetBorderColor(borderColor)

	d.netView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.netView.SetBorder(true).SetTitle(" Network ").SetTitleColor(tcell.ColorDarkGray).SetBorderColor(borderColor)

	// CENTER: Process Table
	d.procTable = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	d.procTable.SetBorder(true).SetTitle(" Processes ").SetTitleColor(tcell.ColorWhite).SetBorderColor(borderColor)

	// RIGHT: System Info + Shortcuts + Help
	d.sysInfoView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.sysInfoView.SetBorder(true).SetTitle(" System ").SetTitleColor(tcell.ColorDarkGray).SetBorderColor(tcell.ColorDarkGray)

	d.shortcutView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.shortcutView.SetBorder(true).SetTitle(" Quick Keys ").SetTitleColor(tcell.ColorAqua).SetBorderColor(borderColor)

	d.helpView = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	d.helpView.SetBorder(true).SetTitle(" Help & Tips ").SetTitleColor(tcell.ColorYellow).SetBorderColor(borderColor)

	// STATUS BAR
	d.statusBar = tview.NewTextView().SetDynamicColors(true).SetScrollable(false)

	// LAYOUT
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(d.cpuView, 0, 1, false).
		AddItem(d.memView, 0, 1, false)

	secondaryRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(d.diskView, 0, 1, false).
		AddItem(d.netView, 0, 1, false)

	leftSide := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topRow, 5, 0, false).
		AddItem(secondaryRow, 5, 0, false).
		AddItem(d.procTable, 0, 3, true)

	rightSidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.sysInfoView, 0, 1, false).
		AddItem(d.shortcutView, 0, 1, false).
		AddItem(d.helpView, 0, 1, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftSide, 0, 4, false).
		AddItem(rightSidebar, 35, 0, false)

	d.rootLayout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.tabBar(), 1, 0, false).
		AddItem(mainFlex, 0, 1, false).
		AddItem(d.statusBar, 1, 0, false)

	d.updateStatusBar()
	d.renderShortcuts()
	d.renderHelp()

	// Set focus to process table
	d.app.SetFocus(d.procTable)

	d.app.SetRoot(d.rootLayout, true)

	// Key bindings
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Filter mode
		if d.filterMode {
			if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
				if len(d.processFilter) > 0 {
					d.processFilter = d.processFilter[:len(d.processFilter)-1]
				}
				d.updater.ProcessFilter = d.processFilter
				return nil
			}
			if event.Rune() != 0 && event.Rune() >= 32 && event.Rune() <= 126 {
				d.processFilter += string(event.Rune())
				d.updater.ProcessFilter = d.processFilter
				return nil
			}
			if event.Key() == tcell.KeyEscape {
				d.filterMode = false
				d.processFilter = ""
				d.updater.ProcessFilter = ""
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.filterMode = false
				d.updater.ProcessFilter = d.processFilter
				return nil
			}
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				if d.exporter != nil {
					d.exporter.Close()
				}
				d.updater.Stop()
				d.app.Stop()
				return nil
			}
			return nil
		}

		// Handle character keys
		switch event.Rune() {
		case 'q', 'Q':
			if d.exporter != nil {
				d.exporter.Close()
			}
			d.updater.Stop()
			d.app.Stop()
			return nil
		case 'e', 'E':
			d.showCores = !d.showCores
			return nil
		case '/':
			d.filterMode = true
			d.processFilter = ""
			d.updater.ProcessFilter = ""
			return nil
		case 'c', 'C':
			d.sortBy = metrics.SortByCPU
			d.updater.SortBy = d.sortBy
			return nil
		case 'm', 'M':
			d.sortBy = metrics.SortByMemory
			d.updater.SortBy = d.sortBy
			return nil
		case 'p', 'P':
			d.sortBy = metrics.SortByPID
			d.updater.SortBy = d.sortBy
			return nil
		case 'n', 'N':
			d.sortBy = metrics.SortByName
			d.updater.SortBy = d.sortBy
			return nil
		case 'k', 'K':
			d.killSelected()
			return nil
		case 'i', 'I':
			d.showProcessDetails()
			return nil
		case '+':
			d.reniceSelected()
			return nil
		case 's', 'S':
			d.sendSignalMenu()
			return nil
		case 'H', '?':
			d.toggleHelp()
			return nil
		}

		// Let arrow keys and tab switching through
		return event
	})

	// Show welcome modal on first run
	if d.firstRun {
		d.showWelcome()
		d.firstRun = false
	}
}

func (d *Dashboard) showWelcome() {
	modal := tview.NewModal()
	modal.SetTitle(" 🎉 Welcome to LightMon! ")
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	modal.SetTextColor(tcell.ColorWhite)

	welcomeText := `[white]Lightweight Terminal System Monitor

[yellow]Quick Start:[white]
  • Press [green]h[white] or [green]?[white] to toggle help panel
  • Press [green]/[white] to filter/search processes
  • Press [green]k[white] to kill a selected process
  • Press [green]q[white] to quit

[yellow]Navigation:[white]
  • Use [green]↑↓[white] arrows to navigate processes
  • Press [green]Enter[white] to select, [green]Esc[white] to cancel

[yellow]View Options:[white]
  • [green]e[white] - Toggle CPU cores view
  • [green]c/m/p/n[white] - Sort processes

Press [green]Enter[white] to start monitoring or [green]Esc[white] to quit`

	modal.SetText(welcomeText)
	modal.AddButtons([]string{"Start Monitoring", "Quit"})
	modal.SetFocus(0)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			d.app.Stop()
		} else {
			d.app.SetRoot(d.rootLayout, true)
			d.app.SetFocus(d.procTable)
		}
	})

	d.app.SetRoot(modal, false)
}

func (d *Dashboard) killSelected() {
	row, _ := d.procTable.GetSelection()
	if row < 1 {
		return
	}
	cell := d.procTable.GetCell(row, 0)
	if cell == nil {
		return
	}
	var pid int32
	fmt.Sscanf(cell.Text, "%d", &pid)
	nameCell := d.procTable.GetCell(row, 1)
	procName := "unknown"
	if nameCell != nil {
		procName = nameCell.Text
	}
	ShowKillMenu(d.app, pid, procName, func(signal syscall.Signal) {
		err := syscall.Kill(int(pid), signal)
		if err != nil {
			d.statusBar.SetText(fmt.Sprintf("[red]Error:[white] %v", err))
		} else {
			d.statusBar.SetText(fmt.Sprintf("Sent signal %d to PID %d", signal, pid))
		}
	}, func() {
		d.app.SetRoot(d.rootLayout, true)
		d.app.SetFocus(d.procTable)
	})
}

func (d *Dashboard) showProcessDetails() {
	row, _ := d.procTable.GetSelection()
	if row < 1 {
		return
	}
	cell := d.procTable.GetCell(row, 0)
	if cell == nil {
		return
	}
	var pid int32
	fmt.Sscanf(cell.Text, "%d", &pid)
	nameCell := d.procTable.GetCell(row, 1)
	procName := "unknown"
	if nameCell != nil {
		procName = nameCell.Text
	}
	ShowProcessDetails(d.app, ProcessInfo{PID: pid, Name: procName, User: "-"}, func() {
		d.app.SetRoot(d.rootLayout, true)
		d.app.SetFocus(d.procTable)
	})
}

func (d *Dashboard) reniceSelected() {
	row, _ := d.procTable.GetSelection()
	if row < 1 {
		return
	}
	cell := d.procTable.GetCell(row, 0)
	if cell == nil {
		return
	}
	var pid int32
	fmt.Sscanf(cell.Text, "%d", &pid)
	ShowReniceMenu(d.app, pid, 0, func(nice int32) {
		d.statusBar.SetText(fmt.Sprintf("Renice PID %d → %d", pid, nice))
	}, func() {
		d.app.SetRoot(d.rootLayout, true)
		d.app.SetFocus(d.procTable)
	})
}

func (d *Dashboard) sendSignalMenu() {
	row, _ := d.procTable.GetSelection()
	if row < 1 {
		return
	}
	cell := d.procTable.GetCell(row, 0)
	if cell == nil {
		return
	}
	var pid int32
	fmt.Sscanf(cell.Text, "%d", &pid)
	nameCell := d.procTable.GetCell(row, 1)
	procName := "unknown"
	if nameCell != nil {
		procName = nameCell.Text
	}
	ShowSignalsMenu(d.app, pid, procName, func(signal syscall.Signal) {
		err := syscall.Kill(int(pid), signal)
		if err != nil {
			d.statusBar.SetText(fmt.Sprintf("[red]Error:[white] %v", err))
		} else {
			d.statusBar.SetText(fmt.Sprintf("Signal %d sent to PID %d", signal, pid))
		}
	}, func() {
		d.app.SetRoot(d.rootLayout, true)
		d.app.SetFocus(d.procTable)
	})
}

func (d *Dashboard) updateStatusBar() {
	sortNames := map[metrics.SortBy]string{
		metrics.SortByCPU:    "CPU",
		metrics.SortByMemory: "MEM",
		metrics.SortByPID:    "PID",
		metrics.SortByName:   "NAME",
	}[d.sortBy]

	filterInfo := ""
	if d.processFilter != "" {
		filterInfo = fmt.Sprintf("  [yellow]/%s[white]", d.processFilter)
	}

	d.statusBar.SetText(fmt.Sprintf(
		" [dim]%s[white]  [dim]│[white]  Sort: [cyan]%s[white]  [dim]│[white]  [green]%d[white] procs  [dim]│[white]  [yellow]h[white]:help  [yellow]q[white]:quit  [yellow]/[white]:filter%s",
		time.Now().Format("15:04:05"),
		sortNames,
		d.updater.ProcLimit,
		filterInfo,
	))
}

func (d *Dashboard) renderShortcuts() {
	d.shortcutView.Clear()
	fmt.Fprint(d.shortcutView, "\n")
	fmt.Fprint(d.shortcutView, "  [green]e[white]  Toggle CPU cores\n")
	fmt.Fprint(d.shortcutView, "  [green]c[white]  Sort by CPU\n")
	fmt.Fprint(d.shortcutView, "  [green]m[white]  Sort by MEM\n")
	fmt.Fprint(d.shortcutView, "  [green]k[white]  Kill process\n")
	fmt.Fprint(d.shortcutView, "  [green]i[white]  Process info\n")
	fmt.Fprint(d.shortcutView, "  [green]/[white] Filter/Search\n")
}

func (d *Dashboard) renderHelp() {
	d.helpView.Clear()
	helpText := `
[dim]MONITORING:[white]
  [green]1[white]  Overview tab
  [green]2[white]  Processes tab

[dim]PROCESS MGMT:[white]
  [green]k[white]  Kill (choose signal)
  [green]i[white]  View details
  [green]+[white] Renice (priority)
  [green]s[white]  Send signal

[dim]TIPS:[white]
  • Click process to select
  • [green]↑↓[white] to navigate
  • [green]Esc[white] to cancel
  • Colors: 🟢OK 🟡⚠ 🔴Critical
`
	fmt.Fprint(d.helpView, helpText)
}

func (d *Dashboard) toggleHelp() {
	d.renderHelp()
}

func (d *Dashboard) tabBar() *tview.TextView {
	tab := tview.NewTextView().SetDynamicColors(true).SetScrollable(false)
	tabs := []string{"Overview", "Processes", "Help"}
	for i, t := range tabs {
		if i == 0 {
			fmt.Fprintf(tab, " [black:white:bold] %d: %s [:-:-:-] ", i+1, t)
		} else {
			fmt.Fprintf(tab, " [dim]%d: %s [:-:-:-] ", i+1, t)
		}
	}
	return tab
}

func (d *Dashboard) Run() error {
	d.updater.Start()
	go func() {
		for snap := range d.updater.Snapshots {
			snapCopy := snap
			d.app.QueueUpdateDraw(func() {
				d.render(snapCopy)
			})
			if d.exporter != nil {
				d.exporter.Export(snapCopy)
			}
		}
	}()
	return d.app.Run()
}

func (d *Dashboard) render(snap internal.SystemSnapshot) {
	d.renderCPU(snap)
	d.renderMemory(snap)
	d.renderDisk(snap)
	d.renderNetwork(snap)
	d.renderProcesses(snap)
	d.renderSysInfo(snap)
	d.updateStatusBar()
}

func (d *Dashboard) renderCPU(snap internal.SystemSnapshot) {
	d.cpuView.Clear()
	cpu := snap.CPU

	bar := d.progressBar(cpu.TotalPercent, barWidth)
	color := d.colorForPercent(cpu.TotalPercent)

	fmt.Fprintf(d.cpuView, "\n")
	fmt.Fprintf(d.cpuView, "  [bold]%6.1f%%[white]  ", cpu.TotalPercent)
	fmt.Fprintf(d.cpuView, "[%s]%s[white]\n", color, bar)

	if snap.CPUTemp > 0 {
		tempColor := "green"
		if snap.CPUTemp >= 85 {
			tempColor = "red"
		} else if snap.CPUTemp >= 75 {
			tempColor = "yellow"
		}
		fmt.Fprintf(d.cpuView, "  [dim]Temp:[white] [%s]%3.0f°C[white]", tempColor, snap.CPUTemp)
	}

	if d.showCores && len(cpu.PerCore) > 0 {
		fmt.Fprintf(d.cpuView, "\n  [dim]")
		for i, c := range cpu.PerCore {
			if i > 15 {
				break
			}
			mini := d.progressBar(c, 2)
			col := d.colorForPercent(c)
			fmt.Fprintf(d.cpuView, "[%s]%s ", col, mini)
		}
	}
}

func (d *Dashboard) renderMemory(snap internal.SystemSnapshot) {
	d.memView.Clear()
	m := snap.Memory

	bar := d.progressBar(m.UsedPercent, barWidth)
	color := d.colorForPercent(m.UsedPercent)

	fmt.Fprintf(d.memView, "\n")
	fmt.Fprintf(d.memView, "  [bold]%6.1f%%[white]  ", m.UsedPercent)
	fmt.Fprintf(d.memView, "[%s]%s[white]\n", color, bar)
	fmt.Fprintf(d.memView, "  [dim]%s[white] / %s", metrics.FormatBytes(m.Used), metrics.FormatBytes(m.Total))

	if m.SwapTotal > 0 {
		swapBar := d.progressBar(m.SwapPercent, barWidth)
		swapColor := d.colorForPercent(m.SwapPercent)
		fmt.Fprintf(d.memView, "\n  [dim]Swap:[white] [%s]%s[white] %5.1f%%", swapColor, swapBar, m.SwapPercent)
	}
}

func (d *Dashboard) renderDisk(snap internal.SystemSnapshot) {
	d.diskView.Clear()
	fmt.Fprintf(d.diskView, "\n")

	for i, disk := range snap.Disks {
		if i > 2 {
			break
		}
		bar := d.progressBar(disk.UsedPercent, barWidth)
		color := d.colorForPercent(disk.UsedPercent)
		label := disk.Path
		if len(label) > 12 {
			label = "..." + label[len(label)-9:]
		}
		fmt.Fprintf(d.diskView, "  [dim]%s[white] [%s]%s[white] %5.0f%%\n", label, color, bar, disk.UsedPercent)
	}
}

func (d *Dashboard) renderNetwork(snap internal.SystemSnapshot) {
	d.netView.Clear()
	fmt.Fprintf(d.netView, "\n")

	if len(snap.Network) == 0 {
		fmt.Fprintf(d.netView, "  [dim]No network[white]")
		return
	}

	for i, iface := range snap.Network {
		if i > 1 {
			break
		}
		sendRate := metrics.FormatBytes(uint64(iface.SendRate))
		recvRate := metrics.FormatBytes(uint64(iface.RecvRate))
		fmt.Fprintf(d.netView, "  [dim]%s:[white] [green]↑%s/s[white] [blue]↓%s/s[white]\n", iface.Interface, sendRate, recvRate)
	}
}

func (d *Dashboard) renderProcesses(snap internal.SystemSnapshot) {
	table := d.procTable
	table.Clear()

	headers := []string{"PID", "NAME", "CPU%", "MEM%", "STATUS"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorAqua).
			SetSelectable(false).
			SetExpansion(1).
			SetAttributes(tcell.AttrBold)
		table.SetCell(0, col, cell)
	}

	for row, p := range snap.Processes {
		r := row + 1
		cpuColor := d.colorForPercent(p.CPU)
		memColor := d.colorForPercent(float64(p.Memory))

		table.SetCell(r, 0, tview.NewTableCell(fmt.Sprintf("%5d", p.PID)).SetTextColor(tcell.ColorDarkGray))
		table.SetCell(r, 1, tview.NewTableCell(d.truncate(p.Name, 28)).SetTextColor(tcell.ColorWhite).SetExpansion(2))
		table.SetCell(r, 2, tview.NewTableCell(fmt.Sprintf("[%s]%7.1f", cpuColor, p.CPU)).SetAlign(tview.AlignRight))
		table.SetCell(r, 3, tview.NewTableCell(fmt.Sprintf("[%s]%7.1f", memColor, p.Memory)).SetAlign(tview.AlignRight))
		table.SetCell(r, 4, tview.NewTableCell(d.statusSymbol(p.Status)).SetAlign(tview.AlignCenter))
	}
}

func (d *Dashboard) renderSysInfo(snap internal.SystemSnapshot) {
	d.sysInfoView.Clear()
	info := snap.SystemInfo

	fmt.Fprintf(d.sysInfoView, "\n")
	fmt.Fprintf(d.sysInfoView, "  [dim]Uptime:[white]   %s\n", info.UptimeString)
	fmt.Fprintf(d.sysInfoView, "  [dim]Load:[white]     %.2f %.2f %.2f\n", info.Load1, info.Load5, info.Load15)
	fmt.Fprintf(d.sysInfoView, "  [dim]Procs:[white]    %d\n", info.ProcessCount)
	if info.BootTime > 0 {
		bootTime := time.Unix(int64(info.BootTime), 0).Format("Jan 02 15:04")
		fmt.Fprintf(d.sysInfoView, "  [dim]Booted:[white]   %s\n", bootTime)
	}
}

func (d *Dashboard) progressBar(pct float64, width int) string {
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func (d *Dashboard) colorForPercent(pct float64) string {
	switch {
	case pct >= 90:
		return "red"
	case pct >= 70:
		return "yellow"
	default:
		return "green"
	}
}

func (d *Dashboard) truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (d *Dashboard) statusSymbol(s string) string {
	switch s {
	case "R", "running":
		return "[green]●[white]"
	case "S", "sleep":
		return "[dim]○[white]"
	case "Z", "zombie":
		return "[red]●[white]"
	default:
		return "[gray]○[white]"
	}
}
