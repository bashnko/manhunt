package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/bashnko/manhunt/external/browsers"
	"github.com/bashnko/manhunt/external/runners"
	"github.com/bashnko/manhunt/internal/bookmarks"
	"github.com/bashnko/manhunt/internal/commands"
	"github.com/bashnko/manhunt/internal/config"
	"github.com/bashnko/manhunt/internal/search"
)

func Run(args []string) error {
	if len(args) > 0 && args[0] == "init" {
		return config.Initialize("")
	}

	cfg, configPath, err := loadConfig()
	if err != nil {
		return err
	}
	prompt := buildPrompt(cfg)
	startupItems := commands.StartupItems(cfg)
	startupRows := len(commands.Items(cfg))
	selection, err := runners.Rofi{}.SelectWithLines(prompt, startupItems, startupRows)
	if err != nil {
		return err
	}

	selection = strings.TrimSpace(selection)
	if selection == "" {
		return nil
	}

	selection, openPrivate := extractPrivateSelection(selection, cfg)

	if commands.IsInput(selection, cfg) {
		return runCommand(selection, cfg, configPath)
	}

	if bookmarks.IsLinksInput(selection) {
		return runSlashLinks(selection, cfg, openPrivate)
	}

	url, err := search.Resolve(selection, cfg)
	if err != nil {
		return err
	}

	return openURL(url, openPrivate)
}

func buildPrompt(cfg config.Config) string {
	return "manhunt search "
}

func openURL(target string, private bool) error {
	command := os.Getenv("BROWSER")
	return browsers.Open(target, command, private)
}

func runCommand(selection string, cfg config.Config, configPath string) error {
	selectedCommand := commands.Selection(selection)
	if selectedCommand == commands.Prefix(cfg) || selectedCommand == ":help" {
		return runCommandMenu(cfg, configPath)
	}
	if commands.IsLinks(selectedCommand, cfg) {
		return runLinksMode(cfg, false)
	}
	if commands.IsAddURL(selectedCommand, cfg) {
		return runAddURLMode(configPath, cfg)
	}
	return runCommandMenu(cfg, configPath)

}

func runCommandMenu(cfg config.Config, configPath string) error {
	items := commands.Items(cfg)
	selection, err := runners.Rofi{}.Select("commands", items)
	if err != nil {
		return err
	}

	selection = strings.TrimSpace(selection)
	if selection == "" {
		return nil
	}

	return runCommand(selection, cfg, configPath)
}

func runLinksMode(cfg config.Config, private bool) error {
	items := bookmarks.SlashItems(cfg)
	if len(items) == 0 {
		return fmt.Errorf("no bookmarks configured")
	}

	selection, err := runners.Rofi{}.Select("links", items)
	if err != nil {
		return err
	}

	selection = strings.TrimSpace(selection)
	if selection == "" {
		return nil
	}

	url, err := bookmarks.ResolveSelection(selection, cfg)
	if err != nil {
		return err
	}

	return openURL(url, private)

}

func runSlashLinks(selection string, cfg config.Config, private bool) error {
	trimmed := bookmarks.TrimInput(selection)
	if trimmed == "" {
		return runLinksMode(cfg, private)
	}
	url, err := bookmarks.ResolveSelection(trimmed, cfg)
	if err != nil {
		return runLinksMode(cfg, private)
	}
	return openURL(url, private)

}

func extractPrivateSelection(input string, cfg config.Config) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	specifier := strings.TrimSpace(cfg.PrivTabSpecifire)
	if specifier == "" {
		specifier = config.DefaultConfig().PrivTabSpecifire
	}

	if !strings.HasPrefix(trimmed, specifier) {
		return trimmed, false
	}

	withoutPrefix := strings.TrimSpace(strings.TrimPrefix(trimmed, specifier))
	if withoutPrefix == "" {
		return trimmed, false
	}

	return withoutPrefix, true
}

func runAddURLMode(configPath string, cfg config.Config) error {
	name, err := promptInput("bookmark name")
	if err != nil {
		return err
	}

	keyword, err := promptInput("bookmark keyword")
	if err != nil {
		return err
	}
	urlValue, err := promptInput("bookmarks url")
	if err != nil {
		return err
	}

	bookmark := config.Shortcut{
		Name:    strings.TrimSpace(name),
		Keyword: strings.TrimSpace(keyword),
		URL:     bookmarks.NormalizeURL(urlValue),
	}
	if bookmark.Name == "" || bookmark.Keyword == "" || bookmark.URL == "" {
		return fmt.Errorf("values are empty")
	}

	bookmarks.Upsert(&cfg, bookmark)
	if err := config.SaveConfig(configPath, cfg); err != nil {
		return err
	}

	return nil

}

func loadConfig() (config.Config, string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return config.Config{}, "", err
	}

	configPath := config.ConfigPath(configDir)
	cfg, err := config.LoadConfig(configPath)
	return cfg, configPath, err
}

func promptInput(prompt string) (string, error) {
	return runners.Rofi{}.Select(prompt, nil)
}
