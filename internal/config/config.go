package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const DefaultConfigName = "config.json"

type Shortcut struct {
	Keyword string `json:"keyword"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

type Config struct {
	DefaultEngine    string
	CommandPrefix    string
	LinksCommand     string
	AddURLCommand    string
	PrivTabSpecifire string `json:"priv_tab_specifire"`
	SearchEngines    map[string]string
	Bookmarks        []Shortcut
}

func DefaultConfig() Config {
	return Config{
		DefaultEngine:    "gg",
		CommandPrefix:    ":",
		LinksCommand:     ":links",
		AddURLCommand:    ":add_url",
		PrivTabSpecifire: "!",
		SearchEngines: map[string]string{
			"gg": "https://www.google.com/search?q=%s",
			"yt": "https://www.youtube.com/results?search_query=%s",
			"rd": "https://www.reddit.com/search/?q=%s",
			"so": "https://stackoverflow.com/search?q=%s",
		},
	}
}

func SaveConfig(path string, config Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ConfigPath(configDir string) string {
	return filepath.Join(configDir, "manhunt", DefaultConfigName)
}

func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()

	if strings.TrimSpace(path) == "" {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if saveErr := SaveConfig(path, config); saveErr != nil {
				return Config{}, saveErr
			}
			return config, nil
		}
		return Config{}, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	config.applyDefaults()
	return config, nil
}

func Initialize(configDir string) error {
	if strings.TrimSpace(configDir) == "" {
		var err error
		configDir, err = os.UserConfigDir()
		if err != nil {
			return err
		}
	}

	configPath := ConfigPath(configDir)
	if err := SaveConfig(configPath, DefaultConfig()); err != nil {
		return err
	}
	return nil
}

func SearchEnginesKeys(config Config) []string {
	keys := make([]string, 0, len(config.SearchEngines))
	for key := range config.SearchEngines {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func BookmarkKeys(config Config) []string {
	keys := make([]string, 0, len(config.Bookmarks))
	for _, bookmark := range config.Bookmarks {
		keys = append(keys, bookmark.Keyword)
	}
	sort.Strings(keys)
	return keys
}

func (config *Config) applyDefaults() {
	defaults := DefaultConfig()

	// aaaaahhh! only if someone could do this for me. repeatativeeeeeeeeee
	if config.SearchEngines == nil {
		config.SearchEngines = defaults.SearchEngines
	}
	if config.Bookmarks == nil {
		config.Bookmarks = defaults.Bookmarks
	}
	if strings.TrimSpace(config.DefaultEngine) == "" {
		config.DefaultEngine = defaults.DefaultEngine
	}
	if strings.TrimSpace(config.CommandPrefix) == "" {
		config.CommandPrefix = defaults.CommandPrefix
	}
	if strings.TrimSpace(config.LinksCommand) == "" {
		config.LinksCommand = defaults.LinksCommand
	}
	if strings.TrimSpace(config.AddURLCommand) == "" {
		config.AddURLCommand = defaults.AddURLCommand
	}
	if strings.TrimSpace(config.PrivTabSpecifire) == "" {
		config.PrivTabSpecifire = defaults.PrivTabSpecifire
	}

}
