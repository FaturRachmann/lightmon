package metrics

import (
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo holds process information
type ProcessInfo struct {
	PID       int32
	PPID      int32
	Name      string
	Cmdline   string
	CPU       float64
	Memory    float32
	Status    string
	User      string
	Nice      int32
	NumThreads int32
	CpuNum    int32
}

// SortBy defines sorting method
type SortBy int

const (
	SortByCPU SortBy = iota
	SortByMemory
	SortByPID
	SortByName
)

// GetProcesses returns a list of processes with optional filtering
func GetProcesses(limit int, sortBy SortBy, filter string) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var list []ProcessInfo
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}
		
		// Apply filter if specified
		if filter != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(filter)) {
			continue
		}
		
		cpuPct, _ := p.CPUPercent()
		memPct, _ := p.MemoryPercent()
		status, _ := p.Status()
		statusStr := "?"
		if len(status) > 0 {
			statusStr = status[0]
		}
		
		// Get username
		user, _ := p.Username()
		if user == "" {
			user = "-"
		}
		
		// Get additional info
		ppid, _ := p.Ppid()
		cmdline, _ := p.Cmdline()
		if cmdline == "" {
			cmdline = name
		}
		nice, _ := p.Nice()
		numThreads, _ := p.NumThreads()

		list = append(list, ProcessInfo{
			PID:        p.Pid,
			PPID:       ppid,
			Name:       name,
			Cmdline:    cmdline,
			CPU:        cpuPct,
			Memory:     memPct,
			Status:     statusStr,
			User:       user,
			Nice:       nice,
			NumThreads: numThreads,
			CpuNum:     0,
		})
	}

	// Sort
	switch sortBy {
	case SortByCPU:
		sort.Slice(list, func(i, j int) bool { return list[i].CPU > list[j].CPU })
	case SortByMemory:
		sort.Slice(list, func(i, j int) bool { return list[i].Memory > list[j].Memory })
	case SortByPID:
		sort.Slice(list, func(i, j int) bool { return list[i].PID < list[j].PID })
	case SortByName:
		sort.Slice(list, func(i, j int) bool { return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name) })
	}

	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	return list, nil
}
