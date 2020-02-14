package graph

import (
	"fmt"
)

const (
	white = iota
	gray
	black
)

// ErrCycleDetected causes when graph is cyclic.
type ErrCycleDetected struct {
	sequence []ID
}

// Error is a error implementation.
func (e ErrCycleDetected) Error() string {
	return fmt.Sprintf("%s cycle detected", e.sequence)
}

// CheckCycles checks graph g to not to be cyclic.
func CheckCycles(g *Graph) error {
	color := make(map[ID]int)
	for id := range g.Nodes() {
		color[id] = white
	}
	for id := range g.Nodes() {
		if color[id] == white {
			var sequence []ID
			if !check(id, g, color, &sequence) {
				return ErrCycleDetected{trimSequence(sequence)}
			}
		}
	}
	return nil
}

func check(id ID, g *Graph, color map[ID]int, sequence *[]ID) bool {
	*sequence = append(*sequence, id)
	if color[id] == gray {
		return false
	}
	color[id] = gray
	sources, err := g.Sources(id)
	if err != nil {
		panic(err)
	}
	for tid := range sources {
		if !check(tid, g, color, sequence) {
			return false
		}
	}
	color[id] = black
	return true
}

func trimSequence(sequence []ID) []ID {
	last := len(sequence) - 1
	for i := last - 1; i >= 0; i-- {
		if sequence[i] == sequence[last] {
			return sequence[i:]
		}
	}
	panic("you found a bug")
}
