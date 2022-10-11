package main

import (
	"os"

	"staff-service/internal/cli"
)

func main() {
	//os.Setenv("KV_VIPER_FILE", "config.yaml")
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
