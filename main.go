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

type SmtpConf struct {
	Host string
	Pass string
	Port int
	User string
}

func main() {
	conf := struct {
		Check map[string]*Check
	}{}
	err := gcfg.ReadFileInto(&conf, "./checks.ini")
	if err != nil {
		panic(err)
	}

	var smtpConf *SmtpConf
	smtpConf, err = TryMailgunConf()
	if err != nil {
		panic(err)
	}
	if smtpConf == nil {
		panic(fmt.Errorf("No SMTP credentials were found in environment"))
	}

	for name, check := range conf.Check {
		err := setDefaults(name, check)
		if err != nil {
			panic(err)
		}

		canary := NewProbe(check)
		notifier := &Notifier{
			SmtpConf:     smtpConf,
			To:           check.To,
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
