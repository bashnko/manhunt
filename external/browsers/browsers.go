package browsers

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bashnko/manhunt/external/runners"
)

func Open(target string, browserCommand string, private bool) error {
	command, args := commandArgs(browserCommand, target, private)
	return runners.Open(command, args)
}

func commandArgs(browserCommand string, target string, private bool) (string, []string) {
	trimmed := strings.TrimSpace(browserCommand)
	if private {
		trimmed = resolvePrivateBrowserCommand(trimmed)
	}

	if trimmed == "" {
		return "xdg-open", []string{target}
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return "xdg-open", []string{target}
	}

	command := fields[0]
	args := make([]string, 0, len(fields))
	hasTargetPlaceholder := false
	for _, arg := range fields[1:] {
		if strings.Contains(arg, "%s") {
			hasTargetPlaceholder = true
			args = append(args, strings.ReplaceAll(arg, "%s", target))
			continue
		}
		args = append(args, arg)
	}
	if private {
		args = append(args, privateFlags(command)...)
	}
	if !hasTargetPlaceholder {
		args = append(args, target)
	}
	return command, args
}

func resolvePrivateBrowserCommand(browserCommand string) string {
	trimmed := strings.TrimSpace(browserCommand)
	if trimmed != "" && trimmed != "xdg-open" {
		return trimmed
	}

	if detected := detectDefaultBrowserCommand(); detected != "" {
		return detected
	}

	return trimmed
}

func detectDefaultBrowserCommand() string {
	output, err := exec.Command("xdg-settings", "get", "default-web-browser").Output()
	if err != nil {
		return ""
	}

	desktop := strings.ToLower(strings.TrimSpace(string(output)))
	switch desktop {
	case "firefox.desktop", "org.mozilla.firefox.desktop", "librewolf.desktop", "waterfox.desktop":
		return "firefox"
	case "google-chrome.desktop", "google-chrome-stable.desktop":
		return "google-chrome"
	case "chromium.desktop", "chromium-browser.desktop":
		return "chromium"
	case "brave-browser.desktop", "brave.desktop":
		return "brave-browser"
	case "microsoft-edge.desktop", "microsoft-edge-stable.desktop":
		return "microsoft-edge"
	case "vivaldi.desktop":
		return "vivaldi"
	case "opera.desktop":
		return "opera"
	default:
		return ""
	}
}

func privateFlags(command string) []string {
	name := strings.ToLower(filepath.Base(strings.TrimSpace(command)))

	if isFirefoxFamily(name) {
		return []string{"--private-window"}
	}

	if isChromiumFamily(name) {
		return []string{"--incognito"}
	}

	return nil
}

func isFirefoxFamily(name string) bool {
	switch name {
	case "firefox", "firefox-bin", "librewolf", "waterfox":
		return true
	default:
		return false
	}
}

func isChromiumFamily(name string) bool {
	switch name {
	case "chrome", "chromium", "chromium-browser", "google-chrome", "google-chrome-stable", "brave", "brave-browser", "microsoft-edge", "microsoft-edge-stable", "vivaldi", "opera":
		return true
	default:
		return false
	}
}
