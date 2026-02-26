package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

// Theme defines color scheme for the UI
type Theme struct {
	BorderColor      string `yaml:"border_color"`
	TitleColor       string `yaml:"title_color"`
	TextColor        string `yaml:"text_color"`
	HighlightColor   string `yaml:"highlight_color"`
	CriticalColor    string `yaml:"critical_color"`
	WarningColor     string `yaml:"warning_color"`
	NormalColor      string `yaml:"normal_color"`
	BackgroundColor  string `yaml:"background_color"`
}

// Thresholds define alert levels for resources
type Thresholds struct {
	CPUWarning    float64 `yaml:"cpu_warning"`    // Percentage for warning (yellow)
	CPUCritical   float64 `yaml:"cpu_critical"`   // Percentage for critical (red)
	MemWarning    float64 `yaml:"mem_warning"`
	MemCritical   float64 `yaml:"mem_critical"`
	DiskWarning   float64 `yaml:"disk_warning"`
	DiskCritical  float64 `yaml:"disk_critical"`
	TempWarning   float64 `yaml:"temp_warning"`   // Celsius
	TempCritical  float64 `yaml:"temp_critical"`
}

// Display settings
type Display struct {
	ShowCores      bool   `yaml:"show_cores"`
	ShowNetwork    bool   `yaml:"show_network"`
	ShowBattery    bool   `yaml:"show_battery"`
	ShowLoadAvg    bool   `yaml:"show_load_avg"`
	ShowUptime     bool   `yaml:"show_uptime"`
	ProcessLimit   int    `yaml:"process_limit"`
	RefreshRate    string `yaml:"refresh_rate"`
	TimeFormat     string `yaml:"time_format"`
	TemperatureUnit string `yaml:"temperature_unit"` // "celsius" or "fahrenheit"
}

// Logging configuration
type Logging struct {
	Enabled    bool   `yaml:"enabled"`
	FilePath   string `yaml:"file_path"`
	Format     string `yaml:"format"` // "json" or "csv"
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
}

// Alerts configuration
type Alerts struct {
	Enabled       bool     `yaml:"enabled"`
	SoundEnabled  bool     `yaml:"sound_enabled"`
	VisualEnabled bool     `yaml:"visual_enabled"`
	LogToFile     bool     `yaml:"log_to_file"`
	WebhookURL    string   `yaml:"webhook_url"`
	Cooldown      string   `yaml:"cooldown"` // Minimum time between same alert
}

// Config holds all configuration
type Config struct {
	Theme      Thresholds `yaml:"thresholds"`
	Thresholds Thresholds `yaml:"thresholds_config"`
	Display    Display    `yaml:"display"`
	Logging    Logging    `yaml:"logging"`
	Alerts     Alerts     `yaml:"alerts"`
	
	// Runtime config (not saved)
	configPath string
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		Thresholds: Thresholds{
			CPUWarning:   70,
			CPUCritical:  90,
			MemWarning:   75,
			MemCritical:  90,
			DiskWarning:  80,
			DiskCritical: 95,
			TempWarning:  75,
			TempCritical: 85,
		},
		Display: Display{
			ShowCores:       false,
			ShowNetwork:     true,
			ShowBattery:     true,
			ShowLoadAvg:     true,
			ShowUptime:      true,
			ProcessLimit:    20,
			RefreshRate:     "1s",
			TimeFormat:      "15:04:05",
			TemperatureUnit: "celsius",
		},
		Logging: Logging{
			Enabled:    false,
			FilePath:   "~/.lightmon/metrics.log",
			Format:     "json",
			MaxSizeMB:  100,
			MaxBackups: 3,
		},
		Alerts: Alerts{
			Enabled:       true,
			SoundEnabled:  false,
			VisualEnabled: true,
			LogToFile:     true,
			Cooldown:      "5m",
		},
	}
}

// Load reads config from file or creates default
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	
	// Expand ~ to home dir
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, fmt.Errorf("failed to get home dir: %w", err)
		}
		path = filepath.Join(home, ".lightmon", "config.yaml")
	}
	
	cfg.configPath = path
	
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config file
		if err := cfg.Save(); err != nil {
			return cfg, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}
	
	// Read existing config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}
	
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return cfg, nil
}

// Save writes config to file
func (c *Config) Save() error {
	if c.configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		c.configPath = filepath.Join(home, ".lightmon", "config.yaml")
	}
	
	// Create directory if needed
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	
	return os.WriteFile(c.configPath, data, 0644)
}

// GetRefreshRate parses refresh rate string to duration
func (c *Config) GetRefreshRate() (time.Duration, error) {
	return time.ParseDuration(c.Display.RefreshRate)
}

// GetAlertCooldown parses cooldown string to duration
func (c *Config) GetAlertCooldown() (time.Duration, error) {
	return time.ParseDuration(c.Alerts.Cooldown)
}

// GetColor parses color name to tcell.Color
func GetColor(name string) tcell.Color {
	colors := map[string]tcell.Color{
		"black":   tcell.ColorBlack,
		"red":     tcell.ColorRed,
		"green":   tcell.ColorGreen,
		"yellow":  tcell.ColorYellow,
		"blue":    tcell.ColorBlue,
		"aqua":    tcell.ColorAqua,
		"purple":  tcell.ColorPurple,
		"white":   tcell.ColorWhite,
		"gray":    tcell.ColorGray,
		"orange":  tcell.ColorOrange,
		"default": tcell.ColorDefault,
	}
	if c, ok := colors[name]; ok {
		return c
	}
	return tcell.ColorWhite
}

// Path returns the config file path
func (c *Config) Path() string {
	return c.configPath
}
