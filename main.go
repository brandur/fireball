package main

import (
	"code.google.com/p/gcfg"
	"fmt"
)

type State int

const (
	Up = iota
	Down
	Unknown
)

type Check struct {
	CheckInterval int    `gcfg:"check-interval"`
	MaxDownChecks int    `gcfg:"max-down-checks"`
	Method        string `gcfg:"method"`
	StatusCode    int    `gcfg:"status-code"`
	Url           string `gcfg:"url"`
}

func main() {
	conf := struct {
		Check map[string]*Check
	}{}
	err := gcfg.ReadFileInto(&conf, "./checks.ini")
	if err != nil {
		panic(err)
	}

	for _, check := range conf.Check {
		err := setDefaults(check)
		if err != nil {
			panic(err)
		}

		canary := NewProbe(check)
		notifier := NewNotifier(canary.StateChanged)

		go canary.Run()
		go notifier.Run()
	}

	done := make(chan bool)
	<-done
}

func setDefaults(check *Check) error {
	if check.CheckInterval == 0 {
		check.CheckInterval = 60
	}

	if check.MaxDownChecks == 0 {
		check.MaxDownChecks = 2
	}

	if check.Method == "" {
		check.Method = "GET"
	}

	if check.StatusCode == 0 {
		check.StatusCode = 200
	}

	if check.Url == "" {
		return fmt.Errorf("Field `url` is required for a check")
	}

	return nil
}
