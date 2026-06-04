package commands

import (
	"fmt"
	"strings"

	"github.com/h3yng/manhunt/internal/bookmarks"
	"github.com/h3yng/manhunt/internal/config"
)

type Command struct {
	Value       string
	Description string
}

func Prefix(cfg config.Config) string {
	aliases := config.NamespaceAliases(cfg, "commands")
	if len(aliases) > 0 {
		return aliases[0]
	}
	return ":"
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
	commandSource, ok := cfg.Sources["commands"]
	if !ok {
		commandSource = config.DefaultConfig().Sources["commands"]
	}
	commands := []Command{
		{Value: Prefix(cfg), Description: "show available command"},
		{Value: ":help", Description: "show available commands"},
	}
	for _, item := range commandSource.Items {
		commands = append(commands, Command{Value: item.Keyword, Description: item.Name})
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
	for _, alias := range commandAliases("add_url", Prefix(cfg), "add_url") {
		if strings.EqualFold(selection, alias) {
			return true
		}
	}
	return false
}

func IsLinks(input string, cfg config.Config) bool {
	selection := Selection(strings.TrimSpace(input))
	for _, alias := range commandAliases("links", Prefix(cfg), "links") {
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
		if _, ok := seen[alias]; ok {
			continue
		}
		seen[alias] = struct{}{}
		unique = append(unique, alias)
	}
	return unique
}

func StartupItems(cfg config.Config) []string {
	items := make([]string, 0, len(Items(cfg))+len(bookmarks.SlashItems(cfg)))
	items = append(items, Items(cfg)...)
	items = append(items, bookmarks.SlashItems(cfg)...)
	return items
}

func commandName(command string, prefix string) string {
	name := strings.TrimSpace(command)
	name = strings.TrimPrefix(name, prefix)
	return name
}
