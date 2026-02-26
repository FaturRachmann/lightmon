package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/lightmon/battery"
	"github.com/lightmon/history"
	"github.com/lightmon/metrics"
	"github.com/lightmon/system"
)

// SystemSnapshot holds all system metrics at a point in time
type SystemSnapshot struct {
	CPU         metrics.CPUStats
	Memory      metrics.MemStats
	Disks       []metrics.DiskStats
	Processes   []metrics.ProcessInfo
	Network     []metrics.NetStats
	Temperatures []metrics.TempStats
	CPUTemp     float64
	Battery     battery.BatteryStats
	SystemInfo  system.SystemInfo
	Timestamp   time.Time
	Alerts      []Alert
}

// Alert represents a system alert
type Alert struct {
	Level     string // "warning" or "critical"
	Resource  string // "cpu", "memory", "disk", "temperature"
	Message   string
	Timestamp time.Time
}

// Updater collects system metrics and maintains history
type Updater struct {
	mu          sync.RWMutex
	Interval    time.Duration
	SortBy      metrics.SortBy
	ProcLimit   int
	ProcessFilter string
	Snapshots   chan SystemSnapshot
	stop        chan struct{}
	
	// History tracking
	CPUHistory    *history.History
	MemoryHistory *history.History
	
	// Thresholds
	CPUWarning    float64
	CPUCritical   float64
	MemWarning    float64
	MemCritical   float64
	DiskWarning   float64
	DiskCritical  float64
	
	// Alert cooldown
	lastAlertTime map[string]time.Time
	alertCooldown time.Duration
}

func NewUpdater(interval time.Duration) *Updater {
	return &Updater{
		Interval:      interval,
		SortBy:        metrics.SortByCPU,
		ProcLimit:     20,
		Snapshots:     make(chan SystemSnapshot, 1),
		stop:          make(chan struct{}),
		CPUHistory:    history.NewHistory(60, "cpu"),
		MemoryHistory: history.NewHistory(60, "memory"),
		CPUWarning:    70,
		CPUCritical:   90,
		MemWarning:    75,
		MemCritical:   90,
		DiskWarning:   80,
		DiskCritical:  95,
		lastAlertTime: make(map[string]time.Time),
		alertCooldown: 5 * time.Minute,
	}
}

func (u *Updater) Start() {
	go func() {
		// Collect immediately
		snap := u.collect()
		select {
		case u.Snapshots <- snap:
		default:
		}

		ticker := time.NewTicker(u.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-u.stop:
				return
			case <-ticker.C:
				snap := u.collect()
				select {
				case u.Snapshots <- snap:
				default:
					<-u.Snapshots
					u.Snapshots <- snap
				}
			}
		}
	}()
}

func (u *Updater) Stop() {
	close(u.stop)
}

func (u *Updater) collect() SystemSnapshot {
	cpu, _ := metrics.GetCPU()
	mem, _ := metrics.GetMemory()
	disks, _ := metrics.GetDisks()
	procs, _ := metrics.GetProcesses(u.ProcLimit, u.SortBy, u.ProcessFilter)
	net, _ := metrics.GetNetwork()
	temps, _ := metrics.GetTemperatures()
	cpuTemp, _ := metrics.GetCPUTemperature()
	bat, _ := battery.GetBattery()
	sysInfo, _ := system.GetSystemInfo()

	// Update history
	u.CPUHistory.Add(cpu.TotalPercent)
	u.MemoryHistory.Add(mem.UsedPercent)

	// Check alerts
	alerts := u.checkAlerts(cpu, mem, disks, cpuTemp)

	return SystemSnapshot{
		CPU:          cpu,
		Memory:       mem,
		Disks:        disks,
		Processes:    procs,
		Network:      net,
		Temperatures: temps,
		CPUTemp:      cpuTemp,
		Battery:      bat,
		SystemInfo:   sysInfo,
		Timestamp:    time.Now(),
		Alerts:       alerts,
	}
}

// checkAlerts checks thresholds and generates alerts
func (u *Updater) checkAlerts(cpu metrics.CPUStats, mem metrics.MemStats, disks []metrics.DiskStats, cpuTemp float64) []Alert {
	var alerts []Alert
	now := time.Now()

	// CPU alerts
	if cpu.TotalPercent >= u.CPUCritical && u.canAlert("cpu_critical", now) {
		alerts = append(alerts, Alert{
			Level:     "critical",
			Resource:  "cpu",
			Message:   f("CPU usage critical: %.1f%%", cpu.TotalPercent),
			Timestamp: now,
		})
	} else if cpu.TotalPercent >= u.CPUWarning && u.canAlert("cpu_warning", now) {
		alerts = append(alerts, Alert{
			Level:     "warning",
			Resource:  "cpu",
			Message:   f("CPU usage high: %.1f%%", cpu.TotalPercent),
			Timestamp: now,
		})
	}

	// CPU Temperature alerts
	if cpuTemp > 0 {
		if cpuTemp >= 85 && u.canAlert("temp_critical", now) {
			alerts = append(alerts, Alert{
				Level:     "critical",
				Resource:  "temperature",
				Message:   f("CPU temperature critical: %.0f°C", cpuTemp),
				Timestamp: now,
			})
		} else if cpuTemp >= 75 && u.canAlert("temp_warning", now) {
			alerts = append(alerts, Alert{
				Level:     "warning",
				Resource:  "temperature",
				Message:   f("CPU temperature high: %.0f°C", cpuTemp),
				Timestamp: now,
			})
		}
	}

	// Memory alerts
	if mem.UsedPercent >= u.MemCritical && u.canAlert("mem_critical", now) {
		alerts = append(alerts, Alert{
			Level:     "critical",
			Resource:  "memory",
			Message:   f("Memory usage critical: %.1f%%", mem.UsedPercent),
			Timestamp: now,
		})
	} else if mem.UsedPercent >= u.MemWarning && u.canAlert("mem_warning", now) {
		alerts = append(alerts, Alert{
			Level:     "warning",
			Resource:  "memory",
			Message:   f("Memory usage high: %.1f%%", mem.UsedPercent),
			Timestamp: now,
		})
	}

	// Disk alerts
	for _, d := range disks {
		if d.UsedPercent >= u.DiskCritical && u.canAlert(f("disk_critical_%s", d.Path), now) {
			alerts = append(alerts, Alert{
				Level:     "critical",
				Resource:  "disk",
				Message:   f("Disk %s critical: %.1f%%", d.Path, d.UsedPercent),
				Timestamp: now,
			})
		} else if d.UsedPercent >= u.DiskWarning && u.canAlert(f("disk_warning_%s", d.Path), now) {
			alerts = append(alerts, Alert{
				Level:     "warning",
				Resource:  "disk",
				Message:   f("Disk %s high: %.1f%%", d.Path, d.UsedPercent),
				Timestamp: now,
			})
		}
	}

	return alerts
}

func (u *Updater) canAlert(key string, now time.Time) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	lastTime, exists := u.lastAlertTime[key]
	if !exists || now.Sub(lastTime) >= u.alertCooldown {
		u.lastAlertTime[key] = now
		return true
	}
	return false
}

// SetThresholds updates alert thresholds
func (u *Updater) SetThresholds(cpuWarn, cpuCrit, memWarn, memCrit, diskWarn, diskCrit float64) {
	u.CPUWarning = cpuWarn
	u.CPUCritical = cpuCrit
	u.MemWarning = memWarn
	u.MemCritical = memCrit
	u.DiskWarning = diskWarn
	u.DiskCritical = diskCrit
}

// SetAlertCooldown updates the alert cooldown period
func (u *Updater) SetAlertCooldown(cooldown time.Duration) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.alertCooldown = cooldown
}

// Simple sprintf helper
func f(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
