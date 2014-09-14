package main

import (
	"fmt"
)

type Notifier struct {
	StateChanged chan State

	state State
}

func NewNotifier(stateChanged chan State) *Notifier {
	return &Notifier{
		StateChanged: stateChanged,

		state: Unknown,
	}
}

func (n *Notifier) Run() {
	for {
		select {
		case state := <-n.StateChanged:
			switch state {
			case Down:
				// Don't notify coming from unknown because we assume that the
				// user already knows
				if n.state == Up {
					n.notifyDown()
				}
			case Up:
				// Same here: don't notify coming from unknown
				if n.state == Down {
					n.notifyUp()
				}
			}
			n.state = state
		}
	}
}

func (n *Notifier) notifyDown() {
	fmt.Printf("Mailing out: service is DOWN\n")
}

func (n *Notifier) notifyUp() {
	fmt.Printf("Mailing out: service is UP\n")
}
