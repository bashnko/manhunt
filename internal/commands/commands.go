package commands

import (
	"fmt"
	"strings"

	"github.com/bashnko/manhunt/internal/config"
)

type Command struct {
	Value       string
	Description string
}

func Prefix(cfg config.Config) string {
	if strings.TrimSpace(cfg.CommandPrefix) == "" {
		return config.DefaultConfig().CommandPrefix
	}
	return strings.TrimSpace(cfg.CommandPrefix)
}

func IsInput(input string, cfg config.Config) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false
	}
	prefix := Prefix(cfg)
	return trimmed == prefix || strings.HasPrefix(trimmed, prefix)
}

func Items(cfg config.Config) []string {
	commands := []Command{
		{Value: Prefix(cfg), Description: "show available command"},
		{Value: ":help", Description: "show available commands"},
		{Value: cfg.LinksCommand, Description: "browse saved links"},
		{Value: cfg.AddURLCommand, Description: "add a saved link"},
	}

	items := make([]string, 0, len(commands))
	for _, command := range commands {
		items = append(items, fmt.Sprintf("%s\t%s", command.Value, command.Description))
	}
	return items
}

func Selection(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}
	if strings.Contains(trimmed, "\t") {
		return strings.TrimSpace(strings.SplitN(trimmed, "\t", 2)[0])
	}
	return trimmed
}

func IsAddURL(input string, cfg config.Config) bool {
	selection := Selection(input)
	for _, alias := range commandAliases(cfg.AddURLCommand, Prefix(cfg), "add_url") {
		if strings.EqualFold(selection, alias) {
			return true
		}
	}
	return false
}

func IsLinks(input string, cfg config.Config) bool {
	selection := Selection(strings.TrimSpace(input))
	for _, alias := range commandAliases(cfg.LinksCommand, Prefix(cfg), "links") {
		if strings.EqualFold(selection, alias) {
			return true
		}
	}
	return false
}

func commandAliases(command string, prefix string, fallbackName string) []string {
	cleanCommand := strings.TrimSpace(command)
	if cleanCommand == "" {
		cleanCommand = prefix + fallbackName
	}
	name := commandName(cleanCommand, prefix)
	alises := []string{cleanCommand, prefix + name, name}

	seen := map[string]struct{}{}
	unique := make([]string, 0, len(alises))
	for _, alias := range alises {
		alias = strings.TrimSpace(alias)
		if alias == "" {
			continue
		}
		seen[alias] = struct{}{}
		unique = append(unique, alias)
	}
	return unique

}

func commandName(command string, prefix string) string {
	name := strings.TrimSpace(command)
	name = strings.TrimPrefix(name, prefix)
	name = strings.TrimPrefix(name, "/")
	return name
}
