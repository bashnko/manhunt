package app

import (
	"os"

	"github.com/bashnko/manhunt/internal/config"
)

func Run(args []string) error {
	if len(args) > 0 && args[0] == "init" {
		return config.Initialize("")
	}

	cfg, configPath, err := loadConfig()
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() (config.Config, string, error) {
	configDir, err := os.UserCacheDir()
	if err != nil {
		return config.Config{}, "", err
	}

	configPath := config.ConfigPath(configDir)
	cfg, err := config.LoadConfig(configPath)
	return cfg, configPath, err
}
