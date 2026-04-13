package runners

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Selector interface {
	Select(prompt string, items []string) (string, error)
}

func Open(command string, args []string) error {
	return exec.Command(command, args...).Start()
}

func runSelection(command string, args []string, items []string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message != "" {
			return "", fmt.Errorf("%s: %s: %w", command, message, err)
		}
		return "", fmt.Errorf("%s: %w", command, err)
	}
	return strings.TrimSpace(stdout.String()), nil
}
