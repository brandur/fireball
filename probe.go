package main

import (
	"fmt"
	"net/http"
	"time"
)

type Probe struct {
	StateChanged chan State
	Stop         chan bool

	check      *Check
	client     *http.Client
	downChecks int
	state      State
}

func NewProbe(check *Check) *Probe {
	return &Probe{
		StateChanged: make(chan State),
		Stop:         make(chan bool),

		check:      check,
		client:     &http.Client{},
		downChecks: 0,
		state:      Unknown,
	}
}

func (c *Probe) Run() {
	for {
		select {
		case <-c.Stop:
			break
		case <-time.After(time.Duration(c.check.CheckInterval) * time.Second):
			err := c.probe()
			if err != nil {
				c.handleFailure(err)
			} else {
				c.handleSuccess()
			}
		}
	}
}

func (c *Probe) probe() error {
	req, err := http.NewRequest(c.check.Method, c.check.Url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.check.StatusCode {
		return fmt.Errorf("[%v] Unexpected status code: %v (expected: %v)\n",
			c.check.Url, resp.StatusCode, c.check.StatusCode)
	}

	fmt.Printf("[%v] Check okay\n", c.check.Url)

	return nil
}

func (c *Probe) handleFailure(err error) {
	fmt.Printf("%v\n", err.Error())
	c.downChecks += 1
	if c.downChecks >= c.check.MaxDownChecks {
		if c.state != Down {
			c.state = Down
			c.StateChanged <- Down
		}
	}
}

func (c *Probe) handleSuccess() {
	if c.state != Up {
		c.state = Up
		c.StateChanged <- Up
	}
	c.downChecks = 0
}
