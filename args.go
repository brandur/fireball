package main

import (
	"os"
)

const (
	ConfFile = "conf.ini"
)

// run arguments for the program
type Args struct {
	ConfFile           string
	DropboxAccessToken string
}

func GetArgs() (*Args, error) {
	args := &Args{
		ConfFile:           os.Getenv("CONF_FILE"),
		DropboxAccessToken: os.Getenv("DROPBOX_ACCESS_TOKEN"),
	}

	if args.ConfFile == "" {
		args.ConfFile = ConfFile
	}

	return args, nil
}
