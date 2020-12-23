package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
)

// gh completion -s zsh > $ZSH/completions/_gh

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "init-zsh",
				Usage: "setup auto-completion and run fetch-apps",
				Action: func(c *cli.Context) error {
					if err := initZshCompletion(); err != nil {
						return err
					}
					if err := fetchAndStoreAppList(c); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:   "fetch-apps",
				Usage:  "fetches list of heroku appList for auto-completions",
				Action: fetchAndStoreAppList,
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

	appList, err := loadAppList()
	if err != nil {
		log.Fatal(err)
	}

	switch c.NArg() {
	case 0:
		fmt.Println("fetch-apps")
		fmt.Println("init")
		for _, t := range appList {
			fmt.Println(t)
		}
	case 1:
		for _, t := range []string{"logs", "releases", "config:cp"} {
			fmt.Println(t)
		}
	}

	return
}

func fetchAndStoreAppList(c *cli.Context) error {
	fmt.Println("fetching list of apps from heroku")
	herokuApps, err := fetchHerokuApps(c.Context)
	if err != nil {
		return err
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := storeHerokuAppsConfig(dir, herokuApps); err != nil {
		return err
	}
	fmt.Printf("%d app(s) cached to %s\n", len(herokuApps), dir)
	return nil
}

func loadAppList() ([]string, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	appList, err := ioutil.ReadFile(path.Join(dir, hwrapAppsListFile))
	if os.IsNotExist(err) {
		return []string{}, nil
	}

	if err != nil {
		return nil, err
	}

	return strings.Split(string(appList), "\n"), nil
}

const hwrapAppsListFile = "apps.config"
const hwrapCompletionFile = "hwrap.completion"

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "hwrap"), nil
}

const zshCompletionFile = `
#compdef hwrap

_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi

  return
}

echo "YEAH"

compdef _cli_zsh_autocomplete hwrap
`

func initZshCompletion() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	_, err = os.Stat(path.Join(home, ".oh-my-zsh"))
	switch {
	case os.IsNotExist(err):
		return errors.New("oh-my-zsh not installed, can't configure completion for zsh only")
	case err != nil:
		return err
	}

	completionPath := path.Join(home, ".oh-my-zsh", "completions")
	_, err = os.Stat(completionPath)
	switch {
	case os.IsNotExist(err):
		if err := os.Mkdir(completionPath, 0755); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	if err := ioutil.WriteFile(path.Join(completionPath, hwrapCompletionFile), []byte(zshCompletionFile), 0755); err != nil {
		return err
	}
	return nil
}

func storeHerokuAppsConfig(location string, herokuApps []string) error {
	_, err := os.Stat(location)
	switch {
	case os.IsNotExist(err):
		if err := os.Mkdir(location, 0755); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	content := strings.Join(herokuApps, "\n")
	if err := ioutil.WriteFile(path.Join(location, hwrapAppsListFile), []byte(content), 0755); err != nil {
		return err
	}
	return nil
}

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
