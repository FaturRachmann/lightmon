package metrics

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TempStats holds temperature sensor information
type TempStats struct {
	Sensor      string
	Temperature float64 // Celsius
}

// GetTemperatures reads CPU temperature from /sys/class/thermal (Linux only)
func GetTemperatures() ([]TempStats, error) {
	var stats []TempStats
	
	// Try reading from /sys/class/thermal/thermal_zone*/temp
	zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil {
		return stats, err
	}

	for _, zone := range zones {
		data, err := os.ReadFile(zone)
		if err != nil {
			continue
		}
		
		// Temperature is in millidegrees Celsius
		tempStr := strings.TrimSpace(string(data))
		tempMilli, err := strconv.ParseInt(tempStr, 10, 64)
		if err != nil {
			continue
		}
		
		temp := float64(tempMilli) / 1000.0
		
		// Get sensor name from directory
		sensor := filepath.Base(filepath.Dir(zone))
		
		stats = append(stats, TempStats{
			Sensor:      sensor,
			Temperature: temp,
		})
	}

	return stats, nil
}

// GetCPUTemperature returns the first available CPU temperature
func GetCPUTemperature() (float64, error) {
	temps, err := GetTemperatures()
	if err != nil || len(temps) == 0 {
		return 0, err
	}

	// Return the first temperature (usually the CPU package)
	return temps[0].Temperature, nil
}
