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
	for {
		select {
		case state := <-n.StateChanged:
			switch state {
			case Down:
				n.notifyDown()
			case Up:
				n.notifyUp()
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
