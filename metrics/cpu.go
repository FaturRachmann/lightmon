package metrics

import (
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUStats struct {
	TotalPercent float64
	PerCore      []float64
	CoreCount    int
}

func GetCPU() (CPUStats, error) {
	total, err := cpu.Percent(0, false)
	if err != nil {
		return CPUStats{}, err
	}

	perCore, err := cpu.Percent(0, true)
	if err != nil {
		perCore = []float64{}
	}

	info, _ := cpu.Info()
	count := len(perCore)
	if len(info) > 0 {
		count = int(info[0].Cores)
	}

	totalVal := 0.0
	if len(total) > 0 {
		totalVal = total[0]
	}

	return CPUStats{
		TotalPercent: totalVal,
		PerCore:      perCore,
		CoreCount:    count,
	}, nil
}
