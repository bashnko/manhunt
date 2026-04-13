package runners

import "strconv"

type Rofi struct{}

func (Rofi) Select(prompt string, items []string) (string, error) {
	args := []string{"-dmenu", "-i", "-p", prompt}
	return runSelection("rofi", args, items)
}

func (Rofi) SelectWithLines(prompt string, items []string, lines int) (string, error) {
	args := []string{"-dmenu", "-i", "-p", prompt}
	if lines > 0 {
		args = append(args, "-l", strconv.Itoa(lines))
	}
	return runSelection("rofi", args, items)
}
