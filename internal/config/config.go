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

type SearchConfig struct {
	DefaultEngine string            `json:"default_engine,omitempty"`
	Engines       map[string]string `json:"engines,omitempty"`
}

type NamespaceConfig struct {
	Aliases     []string `json:"aliases,omitempty"`
	Routes      []string `json:"routes,omitempty"`
	Fuzzy       bool     `json:"fuzzy,omitempty"`
	Description string   `json:"description,omitempty"`
}

type RouterConfig struct {
	PrivatePrefix string                     `json:"private_prefix,omitempty"`
	Namespaces    map[string]NamespaceConfig `json:"namespaces,omitempty"`
}

type SourceConfig struct {
	Type  string     `json:"type,omitempty"`
	Fuzzy bool       `json:"fuzzy,omitempty"`
	Items []Shortcut `json:"items,omitempty"`
	Root  string     `json:"root,omitempty"`
}

type Config struct {
	Router   RouterConfig            `json:"router,omitempty"`
	Search   SearchConfig            `json:"search,omitempty"`
	Sources  map[string]SourceConfig `json:"sources,omitempty"`
}

func DefaultConfig() Config {
	searchEngines := defaultSearchEngines()
	demoLinks := []Shortcut{
		{
			Keyword: "demo",
			Name:    "Manhunt Demo",
			URL:     "https://github.com/h3yng/manhunt",
		},
	}

	return Config{
		Router: RouterConfig{
			PrivatePrefix: "!",
			Namespaces: map[string]NamespaceConfig{
				"/": {
					Aliases: []string{"/"},
					Fuzzy:  true,
					Routes: []string{"links"},
				},
				":": {
					Aliases: []string{":"},
					Fuzzy:  true,
					Routes: []string{"commands"},
				},
			},
		},
		Search: SearchConfig{
			DefaultEngine: "gg",
			Engines:       searchEngines,
		},
		Sources: map[string]SourceConfig{
			"links": {
				Type:  "bookmarks",
				Fuzzy: true,
				Items: demoLinks,
			},
			"commands": {
				Type:  "actions",
				Fuzzy: true,
				Items: []Shortcut{
					{Keyword: "help", Name: "show available commands"},
					{Keyword: "links", Name: "browse saved links"},
					{Keyword: "add_url", Name: "add a saved link"},
				},
			},
		},
	}
}

func defaultSearchEngines() map[string]string {
	return map[string]string{
		"gg": "https://www.google.com/search?q=%s",
		"yt": "https://www.youtube.com/results?search_query=%s",
		"rd": "https://www.reddit.com/search/?q=%s",
		"so": "https://stackoverflow.com/search?q=%s",
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
	engines := config.EffectiveSearchEngines()
	keys := make([]string, 0, len(engines))
	for key := range engines {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func BookmarkKeys(config Config) []string {
	bookmarks := config.EffectiveBookmarks()
	keys := make([]string, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		keys = append(keys, bookmark.Keyword)
	}
	sort.Strings(keys)
	return keys
}

func NamespaceKeys(config Config) []string {
	namespaces := config.Router.Namespaces
	if len(namespaces) == 0 {
		namespaces = DefaultConfig().Router.Namespaces
	}
	keys := make([]string, 0, len(namespaces))
	for key := range namespaces {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func NamespaceAliases(config Config, route string) []string {
	namespaces := config.Router.Namespaces
	if len(namespaces) == 0 {
		namespaces = DefaultConfig().Router.Namespaces
	}

	aliases := make([]string, 0)
	seen := map[string]struct{}{}
	for _, namespace := range namespaces {
		if !namespaceHasRoute(namespace, route) {
			continue
		}
		for _, alias := range namespace.Aliases {
			alias = strings.TrimSpace(alias)
			if alias == "" {
				continue
			}
			if _, ok := seen[alias]; ok {
				continue
			}
			seen[alias] = struct{}{}
			aliases = append(aliases, alias)
		}
	}
	if len(aliases) == 0 {
		if route == "commands" {
			return []string{":"}
		}
		if route == "links" {
			return []string{"/"}
		}
	}
	return aliases
}

func (config Config) EffectiveSource(name string) (SourceConfig, bool) {
	if config.Sources != nil {
		if source, ok := config.Sources[name]; ok {
			return source, true
		}
	}
	defaults := DefaultConfig().Sources
	source, ok := defaults[name]
	return source, ok
}

func (config Config) EffectiveSearchEngines() map[string]string {
	if len(config.Search.Engines) > 0 {
		return config.Search.Engines
	}
	return DefaultConfig().Search.Engines
}

func (config Config) EffectiveDefaultEngine() string {
	if strings.TrimSpace(config.Search.DefaultEngine) != "" {
		return strings.TrimSpace(config.Search.DefaultEngine)
	}
	return DefaultConfig().Search.DefaultEngine
}

func (config Config) EffectiveBookmarks() []Shortcut {
	if links, ok := config.Sources["links"]; ok && len(links.Items) > 0 {
		return links.Items
	}
	return DefaultConfig().Sources["links"].Items
}

func (config *Config) applyDefaults() {
	defaults := DefaultConfig()

	if strings.TrimSpace(config.Search.DefaultEngine) == "" {
		config.Search.DefaultEngine = defaults.Search.DefaultEngine
	}
	if len(config.Search.Engines) == 0 {
		config.Search.Engines = defaults.Search.Engines
	}
	if strings.TrimSpace(config.Router.PrivatePrefix) == "" {
		config.Router.PrivatePrefix = defaults.Router.PrivatePrefix
	}
	if len(config.Router.Namespaces) == 0 {
		config.Router.Namespaces = defaults.Router.Namespaces
	}
	if len(config.Sources) == 0 {
		config.Sources = defaults.Sources
	}
	if links, ok := config.Sources["links"]; ok && len(links.Items) == 0 {
		links.Items = defaults.Sources["links"].Items
		config.Sources["links"] = links
	}
	if commands, ok := config.Sources["commands"]; ok && len(commands.Items) == 0 {
		commands.Items = defaults.Sources["commands"].Items
		config.Sources["commands"] = commands
	}
}

func namespaceHasRoute(namespace NamespaceConfig, route string) bool {
	for _, candidate := range namespace.Routes {
		if candidate == route {
			return true
		}
	}
	return false
}
