package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli/v2"

	"geo-weather-cli/cli"
)

func main() {
	app := &cli.App{
		Name:     "geo-weather-tui",
		Version:  "1.0.0",
		Usage:    "Interactive geolocation and weather tool",
		Commands: commands.GetCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
