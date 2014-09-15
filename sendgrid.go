package main

import (
	"fmt"
	"os"
)

func TrySendGridConf() (*SmtpConf, error) {
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
