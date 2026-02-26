package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/FaturRachmann/lightmon/internal"
)

// Exporter handles metrics export to various formats
type Exporter struct {
	mu         sync.Mutex
	filePath   string
	format     string
	file       *os.File
	csvWriter  *csv.Writer
	enabled    bool
	lastExport time.Time
}

// MetricsRecord represents a single metrics snapshot for export
type MetricsRecord struct {
	Timestamp    time.Time `json:"timestamp"`
	CPU          float64   `json:"cpu_percent"`
	MemoryUsed   uint64    `json:"memory_used"`
	MemoryTotal  uint64    `json:"memory_total"`
	MemoryPercent float64  `json:"memory_percent"`
	DiskUsed     uint64    `json:"disk_used"`
	DiskTotal    uint64    `json:"disk_total"`
	NetSent      uint64    `json:"net_sent"`
	NetRecv      uint64    `json:"net_recv"`
	ProcCount    int       `json:"process_count"`
}

// NewExporter creates a new exporter
func NewExporter(filePath, format string, enabled bool) (*Exporter, error) {
	exp := &Exporter{
		filePath: filePath,
		format:   format,
		enabled:  enabled,
	}

	if !enabled {
		return exp, nil
	}

	// Expand ~ to home dir
	if filePath[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(home, filePath[2:])
	}

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	exp.file = file

	if format == "csv" {
		exp.csvWriter = csv.NewWriter(file)
		// Write header if file is empty
		stat, _ := file.Stat()
		if stat.Size() == 0 {
			exp.csvWriter.Write([]string{
				"timestamp", "cpu_percent", "memory_used", "memory_total",
				"memory_percent", "disk_used", "disk_total", "net_sent", "net_recv", "process_count",
			})
			exp.csvWriter.Flush()
		}
	}

	return exp, nil
}

// Export exports a snapshot to the configured format
func (e *Exporter) Export(snap internal.SystemSnapshot) error {
	if !e.enabled {
		return nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	record := MetricsRecord{
		Timestamp:     snap.Timestamp,
		CPU:           snap.CPU.TotalPercent,
		MemoryUsed:    snap.Memory.Used,
		MemoryTotal:   snap.Memory.Total,
		MemoryPercent: snap.Memory.UsedPercent,
		ProcCount:     len(snap.Processes),
	}

	// Aggregate disk stats
	for _, d := range snap.Disks {
		record.DiskUsed += d.Used
		record.DiskTotal += d.Total
	}

	// Aggregate network stats
	for _, n := range snap.Network {
		record.NetSent += n.BytesSent
		record.NetRecv += n.BytesRecv
	}

	switch e.format {
	case "json":
		return e.exportJSON(record)
	case "csv":
		return e.exportCSV(record)
	default:
		return e.exportJSON(record)
	}
}

func (e *Exporter) exportJSON(record MetricsRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	_, err = e.file.Write(append(data, '\n'))
	return err
}

func (e *Exporter) exportCSV(record MetricsRecord) error {
	if e.csvWriter == nil {
		return fmt.Errorf("CSV writer not initialized")
	}

	err := e.csvWriter.Write([]string{
		record.Timestamp.Format(time.RFC3339),
		fmt.Sprintf("%.2f", record.CPU),
		fmt.Sprintf("%d", record.MemoryUsed),
		fmt.Sprintf("%d", record.MemoryTotal),
		fmt.Sprintf("%.2f", record.MemoryPercent),
		fmt.Sprintf("%d", record.DiskUsed),
		fmt.Sprintf("%d", record.DiskTotal),
		fmt.Sprintf("%d", record.NetSent),
		fmt.Sprintf("%d", record.NetRecv),
		fmt.Sprintf("%d", record.ProcCount),
	})
	if err != nil {
		return err
	}

	e.csvWriter.Flush()
	return e.csvWriter.Error()
}

// Close closes the exporter and underlying file
func (e *Exporter) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.file != nil {
		return e.file.Close()
	}
	return nil
}

// Rotate rotates the log file (for log rotation)
func (e *Exporter) Rotate() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.file != nil {
		e.file.Close()
	}

	// Rotate: rename current file to .1, .2, etc.
	for i := 9; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", e.filePath, i)
		newPath := fmt.Sprintf("%s.%d", e.filePath, i+1)
		os.Rename(oldPath, newPath)
	}
	os.Rename(e.filePath, fmt.Sprintf("%s.1", e.filePath))

	// Reopen file
	file, err := os.OpenFile(e.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	e.file = file

	if e.format == "csv" {
		e.csvWriter = csv.NewWriter(file)
		e.csvWriter.Write([]string{
			"timestamp", "cpu_percent", "memory_used", "memory_total",
			"memory_percent", "disk_used", "disk_total", "net_sent", "net_recv", "process_count",
		})
		e.csvWriter.Flush()
	}

	return nil
}
