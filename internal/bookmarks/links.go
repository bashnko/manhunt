package bookmarks

import (
	"strings"

	"github.com/bashnko/manhunt/internal/config"
)

func IsLinksInput(input string) bool {
	trimmed := strings.TrimSpace(input)
	return trimmed == "/" || strings.HasPrefix(trimmed, "/")
}

func SlashItems(cfg config.Config) []string {
	items := Items(cfg)
	for i := range items {
		parts := strings.SplitN(items[i], "\t", 2)
		if len(parts) == 2 {
			items[i] = "/" + parts[0] + "\t" + parts[1]
		}
	}
	return items
}

func TrimInput(input string) string {
	selection := strings.TrimSpace(input)
	if strings.Contains(selection, "\t") {
		selection = strings.TrimSpace(strings.SplitN(selection, "\t", 2)[0])
	}
	selection = strings.TrimSpace(selection)
	selection = strings.TrimPrefix(selection, "/")
	return strings.TrimSpace(selection)
}
