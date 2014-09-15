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

	stop := make(chan bool)
	run(conf, smtpConf, stop)

	done := make(chan bool)
	<-done
}

func run(conf *Conf, smtpConf *SmtpConf, stop chan bool) {
	for _, check := range conf.Checks {
		probe := NewProbe(check, stop)
		notifier := &Notifier{
			Check:        check,
			From:         conf.Common.From,
			SmtpConf:     smtpConf,
			StateChanged: probe.StateChanged,
			Stop:         stop,
		}

		go probe.Run()
		go notifier.Run()
	}
}
