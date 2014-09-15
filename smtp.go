package main

import (
	"fmt"
	"os"
	"strconv"
)

type SmtpConf struct {
	Host string
	Pass string
	Port int
	User string
}

func GetSmtpConf() (*SmtpConf, error) {
	configurers := []func() (*SmtpConf, error){
		tryMailgunConf,
		trySendGridConf,
	}

	for _, configurer := range configurers {
		smtpConf, err := configurer()
		if err != nil {
			return nil, err
		}
		if smtpConf != nil {
			return smtpConf, nil
		}
	}

	return nil, fmt.Errorf("No SMTP credentials were found in environment")
}

func tryMailgunConf() (*SmtpConf, error) {
	host := os.Getenv("MAILGUN_SMTP_SERVER")
	pass := os.Getenv("MAILGUN_SMTP_PASSWORD")
	port := os.Getenv("MAILGUN_SMTP_PORT")
	user := os.Getenv("MAILGUN_SMTP_LOGIN")

	if host == "" || pass == "" || port == "" || user == "" {
		return nil, nil
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found Mailgun credentials, user: %v\n", user)
	return &SmtpConf{Host: host, Pass: pass, Port: portInt, User: user}, nil
}

func trySendGridConf() (*SmtpConf, error) {
	host := "smtp.sendgrid.net"
	pass := os.Getenv("SENDGRID_PASSWORD")
	port := 587
	user := os.Getenv("SENDGRID_USERNAME")

	if pass == "" || user == "" {
		return nil, nil
	}

	fmt.Printf("Found SendGrid credentials, user: %v\n", user)
	return &SmtpConf{Host: host, Pass: pass, Port: port, User: user}, nil
}
