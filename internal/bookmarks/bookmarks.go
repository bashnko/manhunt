package bookmarks

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"

	"github.com/h3yng/manhunt/internal/config"
)

func Items(cfg config.Config) []string {
	bookmarks := cfg.EffectiveBookmarks()
	items := make([]string, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		items = append(items, fmt.Sprintf("%s\t%s\t%s", bookmark.Keyword, bookmark.Name, bookmark.URL))
	}
	sort.Strings(items)
	return items
}

func ResolveSelection(input string, cfg config.Config) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("empty bookmark selection")
	}

	linksSource, _ := cfg.EffectiveSource("links")
	fuzzy := linksSource.Fuzzy

	if strings.Contains(trimmed, "\t") {
		parts := strings.Split(trimmed, "\t")
		if len(parts) > 0 {
			last := strings.TrimSpace(parts[len(parts)-1])
			if last != "" {
				return last, nil
			}
		}
	}

	for _, bookmark := range cfg.EffectiveBookmarks() {
		if strings.EqualFold(bookmark.Keyword, trimmed) || strings.EqualFold(bookmark.Name, trimmed) {
			return bookmark.URL, nil

		}
	}

	if !fuzzy {
		return "", fmt.Errorf("bookmark %q not found", trimmed)
	}

	lower := strings.ToLower(trimmed)

	matches := make([]config.Shortcut, 0)
	for _, bookmark := range cfg.EffectiveBookmarks() {
		name := strings.ToLower(bookmark.Name)
		Keyword := strings.ToLower(bookmark.Keyword)
		if strings.Contains(name, lower) || strings.Contains(Keyword, lower) {
			matches = append(matches, bookmark)
		}
	}
	if len(matches) == 1 {
		return matches[0].URL, nil
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple bookmarks match %q, be more specific", trimmed)
	}
	return "", fmt.Errorf("bookmark %q not found", trimmed)

}

func Upsert(cfg *config.Config, bookmark config.Shortcut) {
	bookmark.Keyword = strings.TrimSpace(bookmark.Keyword)
	bookmark.Name = strings.TrimSpace(bookmark.Name)
	bookmark.URL = strings.TrimSpace(bookmark.URL)

	bookmarks := cfg.EffectiveBookmarks()
	updated := make([]config.Shortcut, 0, len(bookmarks)+1)
	for _, existing := range bookmarks {
		if strings.EqualFold(existing.Keyword, bookmark.Keyword) {
			continue
		}
		updated = append(updated, existing)
	}
	if cfg.Sources == nil {
		cfg.Sources = map[string]config.SourceConfig{}
	}
	linksSource, _ := cfg.EffectiveSource("links")
	cfg.Sources["links"] = config.SourceConfig{
		Type:  "bookmarks",
		Fuzzy: linksSource.Fuzzy,
		Items: append(updated, bookmark),
	}

}

func NormalizeURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if directURL, ok := resolveDircetURL(trimmed); ok {
		return directURL
	}
	return trimmed
}

func resolveDircetURL(value string) (string, bool) {
	if parsed, err := url.ParseRequestURI(value); err == nil {
		if parsed.Scheme == "http" || parsed.Scheme == "https" {
			return value, true
		}
	}

	if strings.ContainsAny(value, " \t\n") {
		return "", false
	}

	host := value
	if strings.Contains(value, "/") {
		host = strings.SplitN(value, "/", 2)[0]
	}

	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err == nil {
			host = h
		}
	}

	if strings.Contains(host, ".") {
		return "https://" + value, true
	}

	return "", false
}
