package metrics

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStats struct {
	Path        string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetDisks() ([]DiskStats, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var stats []DiskStats
	seen := map[string]bool{}

	for _, p := range partitions {
		if seen[p.Device] {
			continue
		}
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}
		seen[p.Device] = true
		stats = append(stats, DiskStats{
			Path:        p.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}
	return stats, nil
}
