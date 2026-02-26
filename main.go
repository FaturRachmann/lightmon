package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lightmon/config"
	"github.com/lightmon/export"
	"github.com/lightmon/internal"
	"github.com/lightmon/ui"
)

var (
	version = "2.0.0"
	goos    = "unknown"
	goarch  = "unknown"
)

func main() {
	// CLI flags
	interval := flag.Duration("interval", time.Second, "Refresh interval (e.g. 500ms, 1s, 2s)")
	procs := flag.Int("procs", 20, "Max number of processes to show")
	configFile := flag.String("config", "", "Path to config file (default: ~/.lightmon/config.yaml)")
	showVersion := flag.Bool("version", false, "Show version and exit")
	exportFile := flag.String("export", "", "Export metrics to file (CSV or JSON)")
	exportFormat := flag.String("export-format", "json", "Export format: json or csv")
	noExport := flag.Bool("no-export", false, "Disable metrics export")
	
	flag.Parse()

	if *showVersion {
		fmt.Printf("lightmon v%s (%s/%s)\n", version, goos, goarch)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Parse refresh rate from config if not overridden by flag
	if *interval == time.Second && cfg.Display.RefreshRate != "" {
		if d, err := cfg.GetRefreshRate(); err == nil {
			*interval = d
		}
	}

	// Create exporter
	var exporter *export.Exporter
	if !*noExport && *exportFile != "" {
		exporter, err = export.NewExporter(*exportFile, *exportFormat, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create exporter: %v\n", err)
		}
	} else if !*noExport && cfg.Logging.Enabled {
		exporter, err = export.NewExporter(cfg.Logging.FilePath, cfg.Logging.Format, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create exporter: %v\n", err)
		}
	}

	updater := internal.NewUpdater(*interval)
	updater.ProcLimit = *procs
	
	// Apply thresholds from config
	updater.SetThresholds(
		cfg.Thresholds.CPUWarning,
		cfg.Thresholds.CPUCritical,
		cfg.Thresholds.MemWarning,
		cfg.Thresholds.MemCritical,
		cfg.Thresholds.DiskWarning,
		cfg.Thresholds.DiskCritical,
	)
	
	if cooldown, err := cfg.GetAlertCooldown(); err == nil {
		updater.SetAlertCooldown(cooldown)
	}

	dash := ui.NewDashboard(updater, cfg, exporter)
	if err := dash.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
