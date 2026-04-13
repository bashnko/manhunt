package search

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/bashnko/manhunt/internal/config"
)

func Resolve(input string, cfg config.Config) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("empty query")
	}

	if directURL, ok := resolevDirectUrl(trimmed); ok {
		return directURL, nil
	}

	parts := strings.Fields(trimmed)
	keyword := parts[0]
	queryText := strings.TrimSpace(strings.TrimPrefix(trimmed, keyword))

	if template, ok := cfg.SearchEngines[keyword]; ok {
		return formatTemplate(template, queryText)
	}

	for _, bookmark := range cfg.Bookmarks {
		if bookmark.Keyword == keyword {
			return formatTemplate(bookmark.URL, queryText)
		}
	}

	defaultTemplate, ok := cfg.SearchEngines[cfg.DefaultEngine]
	if !ok {
		return "", fmt.Errorf("default search engine %q not found", cfg.DefaultEngine)
	}
	return formatTemplate(defaultTemplate, trimmed)
}

func formatTemplate(template string, queryText string) (string, error) {
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, url.QueryEscape(queryText)), nil
	}
	return template, nil
}

func resolevDirectUrl(value string) (string, bool) {
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
