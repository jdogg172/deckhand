package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Context    string
	Namespace  string
	Mode       string
	ReadOnly   bool
	Kubeconfig string
	LogLevel   string
	NoColor    bool

	UI struct {
		ShowHeader             bool
		ShowFooter             bool
		RefreshIntervalSeconds int
	}
}

func Load(flags Flags) (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigType("yaml")

	configPath := defaultConfigPath()
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			v.SetConfigFile(configPath)
			if err := v.ReadInConfig(); err != nil {
				return cfg, fmt.Errorf("read config file: %w", err)
			}
		}
	}

	cfg.Mode = firstNonEmpty(flags.Mode, v.GetString("defaultMode"), "ops")
	cfg.ReadOnly = flags.ReadOnly || v.GetBool("readonly")
	cfg.Context = firstNonEmpty(flags.Context)
	cfg.Namespace = firstNonEmpty(flags.Namespace)
	cfg.Kubeconfig = firstNonEmpty(flags.Kubeconfig, os.Getenv("KUBECONFIG"))
	cfg.LogLevel = firstNonEmpty(flags.LogLevel, "info")
	cfg.NoColor = flags.NoColor

	cfg.UI.ShowHeader = true
	cfg.UI.ShowFooter = true
	cfg.UI.RefreshIntervalSeconds = 5

	if v.IsSet("ui.showHeader") {
		cfg.UI.ShowHeader = v.GetBool("ui.showHeader")
	}
	if v.IsSet("ui.showFooter") {
		cfg.UI.ShowFooter = v.GetBool("ui.showFooter")
	}
	if v.IsSet("ui.refreshIntervalSeconds") {
		cfg.UI.RefreshIntervalSeconds = v.GetInt("ui.refreshIntervalSeconds")
	}

	return cfg, nil
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "windows" {
		appData := os.Getenv("AppData")
		if appData == "" {
			return filepath.Join(home, "AppData", "Roaming", "deckhand", "config.yaml")
		}
		return filepath.Join(appData, "deckhand", "config.yaml")
	}

	return filepath.Join(home, ".config", "deckhand", "config.yaml")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
