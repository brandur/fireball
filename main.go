package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

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

	smtpConf, err := GetSmtpConf()
	if err != nil {
		panic(err)
	}

	conf, err := readConf(args)
	if err != nil {
		panic(err)
	}

	stop := make(chan bool)
	run(conf, smtpConf, stop)

	for {
		select {
		case <-time.After(30 * time.Minute):
			fmt.Printf("Reloading configuration\n")
			stop <- true

			// note that exit if there's a config problem when starting up, but
			// on reloads, just log it since we're already doing work
			newConf, err := readConf(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reloading configuration: %v\n",
					err.Error())
			} else {
				conf = newConf
			}

			// make a new channel to drop the leftover value in this one
			stop = make(chan bool)
			run(conf, smtpConf, stop)
		}
	}
}

func readConf(args *Args) (*Conf, error) {
	var conf *Conf
	if args.DropboxAccessToken != "" {
		fmt.Printf("Reading configuration from Dropbox\n")

		api := &Dropbox{
			AccessToken: args.DropboxAccessToken,
		}

		contents, err := api.GetContents(args.ConfFile)
		if err != nil {
			return nil, err
		}

		conf, err = GetConf(contents)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Printf("Reading configuration from file\n")

		contents, err := ioutil.ReadFile(args.ConfFile)
		if err != nil {
			return nil, err
		}

		conf, err = GetConf(string(contents))
		if err != nil {
			return nil, err
		}
	}

	return conf, nil
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
