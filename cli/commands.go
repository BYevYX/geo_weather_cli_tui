package commands

import (
	"github.com/urfave/cli/v2"
	"geo-weather-cli/tui"
)

func GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "coordinates",
			Aliases: []string{"coord"},
			Usage:   "Get coordinates by address",
			Action:  tui.RunBubbleTea(tui.GetCoordinatesModel),
		},
		{
			Name:    "weather",
			Aliases: []string{"w"},
			Usage:   "Get weather information",
			Subcommands: []*cli.Command{
				{
					Name:    "current",
					Usage:   "Current weather",
					Action:  tui.RunBubbleTea(tui.GetWeatherCurrentModel),
				},
				{
					Name:    "forecast",
					Usage:   "Weather forecast",
					Action:  tui.RunBubbleTea(tui.GetWeatherForecastModel),
				},
			},
		},
	}
}
