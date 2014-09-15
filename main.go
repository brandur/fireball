package main

import ()

type State int

const (
	Up = iota
	Down
	Unknown
)

func main() {
	// builds run arguments for the program from env vars
	args, err := GetArgs()
	if err != nil {
		panic(err)
	}

	conf, err := GetConf(args.ConfFile)
	if err != nil {
		panic(err)
	}

	smtpConf, err := GetSmtpConf()
	if err != nil {
		panic(err)
	}

	for _, check := range conf.Checks {
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
