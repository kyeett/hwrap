package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:   "fetch-apps",
				Usage:  "fetches list of heroku appList for auto-completions",
				Action: fetchAndStoreAppList,
			},
			{
				Name:    "completion",
				Aliases: nil,
				Usage:   `outputs file for shell completion`,
				UsageText: "" +
					"hwrap completion -s bash > ~/.local/share/bash-completion/hwrap.sh  # for bash\n   " +
					"hwrap completion -s zsh > $ZSH/completions/_hwrap                   # for zsh",
				Action: outputCompletionFile,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "shell",
						Aliases:  []string{"s"},
						Usage:    "shell type: {zsh|bash}. note: bash might not work :-)",
						Required: true,
					},
				},
			},
		},
		Action:       handleHerokuCommand,
		BashComplete: handleBashCompletion,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func handleBashCompletion(c *cli.Context) {
	appList, err := loadAppListFromCache()
	if err != nil {
		log.Fatal(err)
	}

	switch c.NArg() {
	case 0:
		fmt.Println("completion")
		fmt.Println("fetch-apps")
		for _, t := range appList {
			fmt.Println(t)
		}
	case 1:
		for _, t := range []string{"logs", "releases", "config"} {
			fmt.Println(t)
		}
	}

	return
}
