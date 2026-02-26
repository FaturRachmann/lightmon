package system

import (
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemInfo holds system-level information
type SystemInfo struct {
	Uptime        uint64    // seconds
	UptimeString  string
	Load1         float64   // 1 minute load average
	Load5         float64   // 5 minute load average
	Load15        float64   // 15 minute load average
	ProcessCount  int
	BootTime      uint64
	Hostname      string
	OS            string
	Platform      string
	KernelVersion string
	KernelArch    string
}

// GetSystemInfo returns system information
func GetSystemInfo() (SystemInfo, error) {
	uptime, err := host.Uptime()
	if err != nil {
		uptime = 0
	}

	loadAvg, err := load.Avg()
	if err != nil {
		loadAvg = &load.AvgStat{}
	}

	info, err := host.Info()
	if err != nil {
		info = &host.InfoStat{}
	}

	// Get process count from process package
	procs, err := process.Processes()
	procCount := 0
	if err == nil {
		procCount = len(procs)
	}

	return SystemInfo{
		Uptime:        uptime,
		UptimeString:  formatUptime(uptime),
		Load1:         loadAvg.Load1,
		Load5:         loadAvg.Load5,
		Load15:        loadAvg.Load15,
		ProcessCount:  procCount,
		BootTime:      info.BootTime,
		Hostname:      info.Hostname,
		OS:            info.OS,
		Platform:      info.Platform,
		KernelVersion: info.KernelVersion,
		KernelArch:    info.KernelArch,
	}, nil
}

// formatUptime converts seconds to human-readable string
func formatUptime(seconds uint64) string {
	const (
		day  = 24 * 60 * 60
		hour = 60 * 60
		min  = 60
	)

	if seconds < min {
		return "< 1m"
	}

	days := seconds / day
	hours := (seconds % day) / hour
	minutes := (seconds % hour) / min

	if days > 0 {
		return formatDays(days, hours)
	}
	if hours > 0 {
		return formatHours(hours, minutes)
	}
	return formatMinutes(minutes)
}

func formatDays(days, hours uint64) string {
	if days == 1 {
		return "1d " + formatHours(hours, 0)
	}
	return string(rune('0'+days/10)) + string(rune('0'+days%10)) + "d " + formatHours(hours, 0)
}

func formatHours(hours, minutes uint64) string {
	if hours == 1 {
		return "1h " + formatMinutes(minutes)
	}
	return string(rune('0'+hours/10)) + string(rune('0'+hours%10)) + "h " + formatMinutes(minutes)
}

func formatMinutes(minutes uint64) string {
	if minutes < 10 {
		return string(rune('0'+minutes)) + "m"
	}
	return string(rune('0'+minutes/10)) + string(rune('0'+minutes%10)) + "m"
}

// GetBootTime returns the system boot time
func GetBootTime() (time.Time, error) {
	bootTime, err := host.BootTime()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(bootTime), 0), nil
}

// LoadInfo represents load average information
type LoadInfo struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// GetLoad returns load average information
func GetLoad() (LoadInfo, error) {
	avg, err := load.Avg()
	if err != nil {
		return LoadInfo{}, err
	}
	return LoadInfo{
		Load1:  avg.Load1,
		Load5:  avg.Load5,
		Load15: avg.Load15,
	}, nil
}
