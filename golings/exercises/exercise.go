package exercises

type Exercise struct {
	Name string
	Path string
	Mode string
	Hint string
}

func (e Exercise) State() State {
	solved := GetSolved()
	if solved[e.Name] {
		return Done
	}
	return Pending
}

type State int

const (
	Pending State = iota + 1
	Done
)

func (s State) String() string {
	return [...]string{"Pending", "Done"}[s-1]
}
