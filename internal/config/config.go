package config

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
