package main

import (
	"fmt"
	"net/http"
	"time"
)

type Probe struct {
	StateChanged chan StateChangedArgs
	Stop         chan bool

	check      *Check
	client     *http.Client
	downChecks int
	state      State
}

type StateChangedArgs struct {
	Error error
	State State
}

func NewProbe(check *Check, stop chan bool) *Probe {
	return &Probe{
		StateChanged: make(chan StateChangedArgs),
		Stop:         stop,

		check:      check,
		client:     &http.Client{},
		downChecks: 0,
		state:      Unknown,
	}
}

func (c *Probe) Run() {
loop:
	for {
		select {
		case <-c.Stop:
			fmt.Printf("[%v] Stopping probe\n", c.check.Name)
			// pass the value back in for the next listener
			c.Stop <- true
			break loop
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
		return fmt.Errorf("Unexpected status code: %v (expected: %v)",
			resp.StatusCode, c.check.StatusCode)
	}

	fmt.Printf("[%v] Check okay\n", c.check.Name)

	return nil
}

func (c *Probe) handleFailure(err error) {
	fmt.Printf("[%v] %v\n", c.check.Name, err.Error())
	c.downChecks += 1
	if c.downChecks >= c.check.MaxDownChecks {
		if c.state != Down {
			c.state = Down
			c.StateChanged <- StateChangedArgs{
				Error: err,
				State: Down,
			}
		}
	}
}

func (c *Probe) handleSuccess() {
	if c.state != Up {
		c.state = Up
		c.StateChanged <- StateChangedArgs{
			Error: nil,
			State: Up,
		}
	}
	c.downChecks = 0
}
