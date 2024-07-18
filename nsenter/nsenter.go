package nsenter

import (
	"sync"
)

type Nsenter struct {
	name string
	cmd  []string

	mu sync.Mutex
}

func New() *Nsenter {
	return &Nsenter{
		name: "nsenter",
		cmd:  []string{"-m", "-u", "-i", "-n", "-p", "-t1"},
		mu:   sync.Mutex{},
	}
}
