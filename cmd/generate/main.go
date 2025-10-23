package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Run oapi-codegen with the same arguments
	cmd := exec.Command("oapi-codegen", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
