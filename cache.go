package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const hwrapAppsListFile = "apps.config"

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "hwrap"), nil
}

func loadAppListFromCache() ([]string, error) {
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

func storeAppsToCache(location string, herokuApps []string) error {
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
