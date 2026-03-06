package exercises

import (
	"os"
)

type Exercise struct {
	Name string
	Path string
	Mode string
	Hint string
}

func (e Exercise) State() State {
	_, err := os.Stat(e.Path)
	if err != nil {
		return Pending
	}

	return Done
}

type State int

const (
	Pending State = iota + 1
	Done
)

func (s State) String() string {
	return [...]string{"Pending", "Done"}[s-1]
}
