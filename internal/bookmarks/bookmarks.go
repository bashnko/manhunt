package bookmarks

import (
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
