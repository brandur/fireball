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

	state State
}

func (n *Notifier) Run() {
	for {
		select {
		case args := <-n.StateChanged:
			switch args.State {
			case Down:
				// Don't notify coming from unknown because we assume that the
				// user already knows
				if n.state == Up {
					n.notifyDown(args)
				}
			case Up:
				// Same here: don't notify coming from unknown
				if n.state == Down {
					n.notifyUp(args)
				}
			}
			n.state = args.State
		}
	}
}

func (n *Notifier) notifyDown(args StateChangedArgs) {
	subject := fmt.Sprintf("%v is DOWN", n.Check.Name)
	fmt.Printf("[%v] Mailing out \"%v\" to: %v\n",
		n.Check.Name, subject, n.Check.To)
	body := []byte(subject)
	n.sendMail(subject, body)
}

func (n *Notifier) notifyUp(args StateChangedArgs) {
	subject := fmt.Sprintf("%v is UP", n.Check.Name)
	fmt.Printf("[%v] Mailing out \"%v\" to: %v\n",
		n.Check.Name, subject, n.Check.To)
	body := []byte(subject)
	n.sendMail(subject, body)
}

func (n *Notifier) sendMail(subject string, body []byte) {
	auth := smtp.PlainAuth("", n.SmtpConf.User, n.SmtpConf.Pass, n.SmtpConf.Host)
	addr := n.SmtpConf.Host + ":" + strconv.Itoa(n.SmtpConf.Port)
	err := smtp.SendMail(addr, auth, n.From, n.Check.To, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%v] Failed to send mail: %v\n",
			n.Check.Name, err.Error())
	}
}
