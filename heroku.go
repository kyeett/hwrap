package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

func handleHerokuCommand(c *cli.Context) error {
	// Get all args except the base command
	args := c.Args().Slice()
	if len(args) < 2 {
		fmt.Println(args)
		return errors.New("please specify app and command")
	}
	app := args[0]
	herokuCmd := args[1]
	args2 := []string{herokuCmd, "-a", app}
	args2 = append(args2, args[2:]...)
	cmd := exec.CommandContext(c.Context, "heroku", args2...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func fetchHerokuApps(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "heroku", "apps", "--all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var apps []string
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "===") {
			continue
		}
		line = strings.Split(line, " ")[0]
		apps = append(apps, line)
	}
	return apps, nil
}
