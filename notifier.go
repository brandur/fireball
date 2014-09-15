package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type Notifier struct {
	Check        *Check
	From         string
	SmtpConf     *SmtpConf
	StateChanged chan StateChangedArgs
	Stop         chan bool

	state State
}

func (n *Notifier) Run() {
loop:
	for {
		select {
		case <-n.Stop:
			fmt.Printf("[%v] Stopping notifier\n", n.Check.Name)
			// pass the value back in for the next listener
			n.Stop <- true
			break loop
		case args := <-n.StateChanged:
			switch args.State {
			case Down:
				// Don't notify coming from unknown because we assume that the
				// user already knows
				if n.state == Up {
					go n.notifyDown(args)
				}
			case Up:
				// Same here: don't notify coming from unknown
				if n.state == Down {
					go n.notifyUp(args)
				}
			}
			n.state = args.State
		}
	}
}

func (n *Notifier) mailTemplateDown(name string, url string, err error) string {
	return fmt.Sprintf(`Subject: %v is DOWN

Please note that according to an HTTP test, %v appears to be DOWN.

The following error was encountered while trying to probe:

    %v

The URL being checked is:

    %v

This message was generated automatically by an installation of Fireball:

    https://github.com/brandur/fireball`,
		name, name, err.Error(), url)
}

func (n *Notifier) mailTemplateUp(name string, url string) string {
	return fmt.Sprintf(`Subject: %v is UP

Please note that according to an HTTP test, %v appears to be back UP.

The URL being checked is:

    %v

This message was generated automatically by an installation of Fireball:

    https://github.com/brandur/fireball`,
		name, name, url)
}

func (n *Notifier) notifyDown(args StateChangedArgs) {
	fmt.Printf("[%v] Mailing DOWN to: %v\n", n.Check.Name, n.Check.To)
	body := []byte(n.mailTemplateDown(n.Check.Name, n.Check.Url, args.Error))
	n.sendMail(body)
}

func (n *Notifier) notifyUp(args StateChangedArgs) {
	fmt.Printf("[%v] Mailing UP to: %v\n", n.Check.Name, n.Check.To)
	body := []byte(n.mailTemplateUp(n.Check.Name, n.Check.Url))
	n.sendMail(body)
}

func (n *Notifier) sendMail(body []byte) {
	auth := smtp.PlainAuth("", n.SmtpConf.User, n.SmtpConf.Pass, n.SmtpConf.Host)
	addr := n.SmtpConf.Host + ":" + strconv.Itoa(n.SmtpConf.Port)
	err := smtp.SendMail(addr, auth, n.From, n.Check.To, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%v] Failed to send mail: %v\n",
			n.Check.Name, err.Error())
	}
}
