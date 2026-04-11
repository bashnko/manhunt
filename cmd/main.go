package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd, err := exec.Command("whomai").Output()
	if err != nil {
		log.Fatal()
	}
	fmt.Println(string(cmd))

}
