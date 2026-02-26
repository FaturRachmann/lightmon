package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemStats struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapPercent float64
}

func GetMemory() (MemStats, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return MemStats{}, err
	}

	sw, _ := mem.SwapMemory()
	swapPercent := 0.0
	swapTotal := uint64(0)
	swapUsed := uint64(0)
	if sw != nil {
		swapPercent = sw.UsedPercent
		swapTotal = sw.Total
		swapUsed = sw.Used
	}

	return MemStats{
		Total:       vm.Total,
		Used:        vm.Used,
		Free:        vm.Free,
		UsedPercent: vm.UsedPercent,
		SwapTotal:   swapTotal,
		SwapUsed:    swapUsed,
		SwapPercent: swapPercent,
	}, nil
}

func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
