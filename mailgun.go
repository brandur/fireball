package main

import (
	"fmt"
	"os"
	"strconv"
)

func TryMailgunConf() (*SmtpConf, error) {
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
