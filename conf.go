package main

import (
	"code.google.com/p/gcfg"
	"fmt"
)

type Conf struct {
	Common Common            `gcfg:"common"`
	Checks map[string]*Check `gcfg:"check"`
}

type Common struct {
	From string `gcfg:"from"`
}

type Check struct {
	CheckInterval int      `gcfg:"check-interval"`
	MaxDownChecks int      `gcfg:"max-down-checks"`
	Method        string   `gcfg:"method"`
	Name          string   `gcfg:"name"`
	StatusCode    int      `gcfg:"status-code"`
	To            []string `gcfg:"to"`
	Url           string   `gcfg:"url"`
}

func GetConf(confFile string) (*Conf, error) {
	conf := &Conf{}
	err := gcfg.ReadFileInto(conf, confFile)
	if err != nil {
		return nil, err
	}

	for name, check := range conf.Checks {
		err := initCheck(name, check)
		if err != nil {
			return nil, err
		}
	}

	return conf, nil
}

func initCheck(name string, check *Check) error {
	if check.CheckInterval == 0 {
		check.CheckInterval = 60
	}

	if check.MaxDownChecks == 0 {
		check.MaxDownChecks = 2
	}

	if check.Method == "" {
		check.Method = "GET"
	}

	if check.Name == "" {
		check.Name = name
	}

	if check.StatusCode == 0 {
		check.StatusCode = 200
	}

	if len(check.To) < 1 {
		return fmt.Errorf("At least one `to` field is required for a check")
	}

	if check.Url == "" {
		return fmt.Errorf("`url` field is required for a check")
	}

	return nil
}
