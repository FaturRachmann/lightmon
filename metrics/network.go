package metrics

import (
	"github.com/shirou/gopsutil/v3/net"
)

type NetStats struct {
	Interface string
	BytesSent uint64
	BytesRecv uint64
	SendRate  float64 // bytes/sec
	RecvRate  float64 // bytes/sec
}

var prevNet map[string]net.IOCountersStat

func GetNetwork() ([]NetStats, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var stats []NetStats
	for _, c := range counters {
		if c.Name == "lo" {
			continue
		}
		stat := NetStats{
			Interface: c.Name,
			BytesSent: c.BytesSent,
			BytesRecv: c.BytesRecv,
		}
		if prevNet != nil {
			if prev, ok := prevNet[c.Name]; ok {
				stat.SendRate = float64(c.BytesSent-prev.BytesSent) / 1.0
				stat.RecvRate = float64(c.BytesRecv-prev.BytesRecv) / 1.0
			}
		}
		stats = append(stats, stat)
	}

	// Save current as prev
	prevNet = make(map[string]net.IOCountersStat)
	for _, c := range counters {
		prevNet[c.Name] = c
	}

	return stats, nil
}
