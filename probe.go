package main

import (
	"fmt"
	"net/http"
	"time"
)

type Probe struct {
	CheckInterval time.Duration
	MaxDownChecks int
	Method        string
	StateChanged  chan State
	StatusCode    int
	Stop          chan bool
	Url           string

	client     *http.Client
	downChecks int
	state      State
}

func NewProbe(check *Check) *Probe {
	return &Probe{
		CheckInterval: check.CheckInterval,
		StateChanged:  make(chan State),
		MaxDownChecks: check.MaxDownChecks,
		Method:        check.Method,
		StatusCode:    check.StatusCode,
		Stop:          make(chan bool),
		Url:           check.Url,

		client:     &http.Client{},
		downChecks: 0,
		state:      Up,
	}
}

func (c *Probe) Run() {
	for {
		select {
		case <-c.Stop:
			break
		case <-time.After(c.CheckInterval):
			err := c.check()
			if err != nil {
				c.handleFailure(err)
			} else {
				c.handleSuccess()
			}
		}
	}
}

func (c *Probe) check() error {
	req, err := http.NewRequest(c.Method, c.Url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.StatusCode {
		return fmt.Errorf("[%v] Unexpected status code: %v (expected: %v)\n",
			c.Url, resp.StatusCode, c.StatusCode)
	}

	fmt.Printf("[%v] Check okay\n", c.Url)

	return nil
}

func (c *Probe) handleFailure(err error) {
	fmt.Printf("%v\n", err.Error())
	c.downChecks += 1
	if c.downChecks >= c.MaxDownChecks {
		if c.state == Up {
			c.state = Down
			c.StateChanged <- Down
		}
	}
}

func (c *Probe) handleSuccess() {
	if c.state == Down {
		c.state = Up
		c.StateChanged <- Up
	}
	c.downChecks = 0
}
