package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const DefaultConfigName = "config.json"

type Shortcut struct {
	Keyword string `json:"keyword"`
	Name    string `json:"name"`
	URL     string `json:"uRL"`
}

type Config struct {
	DefaultEngine string
	CommandPrfix  string
	LinksCommand  string
	AddURLCommand string
	SearchEngines map[string]string
	Bookmarks     []Shortcut
}

func DefaultConfig() Config {
	return Config{
		DefaultEngine: "gg",
		CommandPrfix:  ":",
		LinksCommand:  ":links",
		AddURLCommand: ":add_url",
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
