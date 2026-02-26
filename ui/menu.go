package ui

import (
	"fmt"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var signalOptions = []struct {
	Name   string
	Signal syscall.Signal
}{
	{"SIGTERM", syscall.SIGTERM},
	{"SIGKILL", syscall.SIGKILL},
	{"SIGINT", syscall.SIGINT},
	{"SIGHUP", syscall.SIGHUP},
	{"SIGQUIT", syscall.SIGQUIT},
	{"SIGUSR1", syscall.SIGUSR1},
	{"SIGUSR2", syscall.SIGUSR2},
	{"SIGSTOP", syscall.SIGSTOP},
	{"SIGCONT", syscall.SIGCONT},
}

// ShowKillMenu shows an interactive kill menu
func ShowKillMenu(app *tview.Application, pid int32, procName string, onKill func(signal syscall.Signal), done func()) {
	modal := tview.NewModal()
	modal.SetTitle(" 📋 Kill Process ")
	modal.SetBackgroundColor(tcell.ColorMaroon)
	modal.SetTextColor(tcell.ColorWhite)

	buttons := make([]string, len(signalOptions)+1)
	for i, opt := range signalOptions {
		buttons[i] = fmt.Sprintf("%s (%d)", opt.Name, opt.Signal)
	}
	buttons[len(buttons)-1] = "❌ Cancel"

	modal.SetText(fmt.Sprintf("Process: [white]%s[white]\nPID: [yellow]%d[white]\n\nSelect signal to send:", procName, pid))
	modal.AddButtons(buttons)
	modal.SetFocus(0)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "❌ Cancel" {
			done()
			return
		}
		for _, opt := range signalOptions {
			if buttonLabel == fmt.Sprintf("%s (%d)", opt.Name, opt.Signal) {
				onKill(opt.Signal)
				break
			}
		}
		done()
	})

	app.SetRoot(modal, false)
}

// ShowProcessDetails shows detailed process information
func ShowProcessDetails(app *tview.Application, proc ProcessInfo, done func()) {
	modal := tview.NewModal()
	modal.SetTitle(" ℹ️ Process Details ")
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	modal.SetTextColor(tcell.ColorWhite)

	statusText := getStatusDescription(proc.Status)

	text := fmt.Sprintf(`[white]Process:[white]  %s
[white]PID:[white]        %d
[white]User:[white]       %s
[white]Status:[white]     %s (%s)
[white]CPU%%:[white]       %.1f%%
[white]Memory%%:[white]    %.1f%%

Press Enter or click OK to close`,
		proc.Name,
		proc.PID,
		proc.User,
		proc.Status,
		statusText,
		proc.CPU,
		proc.Memory,
	)

	modal.SetText(text)
	modal.AddButtons([]string{"OK"})
	modal.SetFocus(0)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		done()
	})

	app.SetRoot(modal, false)
}

// ShowReniceMenu shows a menu to change process priority
func ShowReniceMenu(app *tview.Application, pid int32, currentNice int32, onRenice func(nice int32), done func()) {
	modal := tview.NewModal()
	modal.SetTitle(" ⚖️ Renice Process ")
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	modal.SetTextColor(tcell.ColorWhite)

	buttons := []string{
		"-20 (Highest)", "-15", "-10", "-5", "0 (Default)",
		"5", "10", "15", "19 (Lowest)",
		"❌ Cancel",
	}

	modal.SetText(fmt.Sprintf("Process PID: [yellow]%d[white]\nCurrent nice: [yellow]%d[white]\n\nSelect new nice value:", pid, currentNice))
	modal.AddButtons(buttons)
	modal.SetFocus(0)

	niceValues := []int32{-20, -15, -10, -5, 0, 5, 10, 15, 19}

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "❌ Cancel" {
			done()
			return
		}
		if buttonIndex >= 0 && buttonIndex < len(niceValues) {
			onRenice(niceValues[buttonIndex])
		}
		done()
	})

	app.SetRoot(modal, false)
}

// ShowSignalsMenu shows a menu with all available signals
func ShowSignalsMenu(app *tview.Application, pid int32, procName string, onSignal func(signal syscall.Signal), done func()) {
	modal := tview.NewModal()
	modal.SetTitle(" 📡 Send Signal ")
	modal.SetBackgroundColor(tcell.ColorDarkGreen)
	modal.SetTextColor(tcell.ColorWhite)

	buttons := make([]string, len(signalOptions)+1)
	for i, opt := range signalOptions {
		buttons[i] = fmt.Sprintf("%s (%d)", opt.Name, opt.Signal)
	}
	buttons[len(buttons)-1] = "❌ Cancel"

	modal.SetText(fmt.Sprintf("Process: [white]%s[white]\nPID: [yellow]%d[white]\n\nSelect signal:", procName, pid))
	modal.AddButtons(buttons)
	modal.SetFocus(0)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "❌ Cancel" {
			done()
			return
		}
		for _, opt := range signalOptions {
			if buttonLabel == fmt.Sprintf("%s (%d)", opt.Name, opt.Signal) {
				onSignal(opt.Signal)
				break
			}
		}
		done()
	})

	app.SetRoot(modal, false)
}

func getStatusDescription(status string) string {
	switch status {
	case "R", "running":
		return "Running"
	case "S", "sleep":
		return "Sleeping"
	case "D", "disk-sleep":
		return "Disk Sleep"
	case "Z", "zombie":
		return "Zombie"
	case "T", "stopped":
		return "Stopped"
	default:
		return "Unknown"
	}
}

type ProcessInfo struct {
	PID    int32
	Name   string
	CPU    float64
	Memory float32
	Status string
	User   string
}
