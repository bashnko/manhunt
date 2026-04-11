package runners

type Rofi struct{}

func (Rofi) Select(prompt string, items []string) (string, error) {
	args := []string{"-dmenu", "-i", "-p", prompt}
	return runSelection("rofi", args, items)
}
