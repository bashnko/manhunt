package bookmarks

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"

	"github.com/bashnko/manhunt/internal/config"
)

func Items(cfg config.Config) []string {
	items := make([]string, 0, len(cfg.Bookmarks))
	for _, bookmark := range cfg.Bookmarks {
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

	if strings.Contains(trimmed, "\t") {
		parts := strings.Split(trimmed, "\t")
		if len(parts) > 0 {
			last := strings.TrimSpace(parts[len(parts)-1])
			if last != "" {
				return last, nil
			}
		}
	}

	lower := strings.ToLower(trimmed)

	for _, bookmark := range cfg.Bookmarks {
		if strings.EqualFold(bookmark.Keyword, trimmed) || strings.EqualFold(bookmark.Name, trimmed) {
			return bookmark.URL, nil

		}
	}

	matches := make([]config.Shortcut, 0)
	for _, bookmark := range cfg.Bookmarks {
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

	if strings.Contains(value, " \t\n") {
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
