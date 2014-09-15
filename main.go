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
	CheckInterval int      `gcfg:"check-interval"`
	MaxDownChecks int      `gcfg:"max-down-checks"`
	Method        string   `gcfg:"method"`
	Name          string   `gcfg:"name"`
	StatusCode    int      `gcfg:"status-code"`
	To            []string `gcfg:"to"`
	Url           string   `gcfg:"url"`
}

func main() {
	conf := struct {
		Common struct {
			From string `gcfg:"from"`
		} `gcfg:"common"`
		Checks map[string]*Check `gcfg:"check"`
	}{}
	err := gcfg.ReadFileInto(&conf, "./checks.ini")
	if err != nil {
		panic(err)
	}

	smtpConf, err := GetSmtpConf()
	if err != nil {
		panic(err)
	}

	for name, check := range conf.Checks {
		err := setDefaults(name, check)
		if err != nil {
			panic(err)
		}

		canary := NewProbe(check)
		notifier := &Notifier{
			Check:        check,
			From:         conf.Common.From,
			SmtpConf:     smtpConf,
			StateChanged: canary.StateChanged,
		}

		go canary.Run()
		go notifier.Run()
	}

	done := make(chan bool)
	<-done
}

func setDefaults(name string, check *Check) error {
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
