package battery

import (
	"fmt"

	battery "github.com/distatus/battery"
)

// BatteryStats holds battery information
type BatteryStats struct {
	Percent       float64
	State         string
	StateSymbol   string
	TimeRemaining int64 // seconds
	HasBattery    bool
	IsCharging    bool
}

// GetBattery returns battery statistics
func GetBattery() (BatteryStats, error) {
	b, err := battery.Get(0)
	if err != nil {
		// No battery or error
		return BatteryStats{HasBattery: false}, err
	}

	stats := BatteryStats{
		HasBattery: true,
		Percent:    0,
		State:      b.State.String(),
	}

	// Calculate percentage
	if b.Full > 0 && b.Current > 0 {
		stats.Percent = (b.Current / b.Full) * 100
	}

	// Get time remaining (in seconds)
	// ChargeRate is in mW (milliwatts), Current is in mWh (milliwatt-hours)
	// Time = Current / ChargeRate * 3600 (convert hours to seconds)
	if b.ChargeRate > 0 && b.Current > 0 {
		stats.TimeRemaining = int64(b.Current / b.ChargeRate * 3600)
	}

	// Check if charging - State wraps AgnosticState
	stats.IsCharging = b.State.String() == "Charging"

	// Set state symbol based on string representation
	switch b.State.String() {
	case "Charging":
		stats.StateSymbol = "⚡"
	case "Full":
		stats.StateSymbol = "🔋"
	case "Empty":
		stats.StateSymbol = "🪫"
	default:
		if stats.Percent > 50 {
			stats.StateSymbol = "🔋"
		} else if stats.Percent > 20 {
			stats.StateSymbol = "🔋"
		} else {
			stats.StateSymbol = "🪫"
		}
	}

	return stats, nil
}

// FormatTimeRemaining formats seconds to human readable string
func FormatTimeRemaining(seconds int64) string {
	if seconds <= 0 {
		return "Unknown"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
