package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bashnko/manhunt/internal/config"
)

func main() {
	cmd, err := exec.Command("whomai").Output()
	if err != nil {
		log.Fatal()
	}
	fmt.Println(string(cmd))
	configPath, err := os.UserCacheDir()
	if err != nil {
		log.Fatal()
	}
	config.SaveConfig(configPath, config.DefaultConfig())

}
