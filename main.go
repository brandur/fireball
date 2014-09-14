package main

import (
	"time"
)

type State int

const (
	Up = iota
	Down
)

type Check struct {
	CheckInterval time.Duration
	MaxDownChecks int
	Method        string
	StatusCode    int
	Url           string
}

func main() {
	checks := [...]Check{
		Check{CheckInterval: 10 * time.Second, MaxDownChecks: 2, Method: "GET", StatusCode: 200, Url: "https://brandur.org"},
		Check{CheckInterval: 10 * time.Second, MaxDownChecks: 2, Method: "GET", StatusCode: 200, Url: "https://mutelight.org"},
	}
	for _, check := range checks {
		canary := NewCanary(&check)
		notifier := NewNotifier(canary.StateChanged)
		go canary.Run()
		go notifier.Run()
	}

	done := make(chan bool)
	<-done
}
