package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

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
	if err := storeAppsToCache(dir, herokuApps); err != nil {
		return err
	}
	fmt.Printf("%d app(s) cached to %s\n", len(herokuApps), dir)
	return nil
}