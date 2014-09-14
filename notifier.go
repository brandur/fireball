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

		state: Up,
	}
}

func (n *Notifier) Run() {
	select {
	case state := <-n.StateChanged:
		if state == Down {
			n.notifyDown()
		}
		n.state = state
	}
}

func (n *Notifier) notifyDown() {
	fmt.Printf("Mailing out\n")
}
