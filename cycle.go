package di

import (
	"fmt"
)

const (
	temporary = 1
	permanent = 2
)

func visit(s schema, node *node, marks map[*node]int) error {
	if marks[node] == permanent {
		return nil
	}
	if marks[node] == temporary {
		return errCycleDetected // todo: improve message
	}
	marks[node] = temporary
	params, err := node.deps(s)
	if err != nil {
		return fmt.Errorf("%s: %s", node, err)
	}
	for _, param := range params {
		if err := visit(s, param, marks); err != nil {
			return err
		}
	}
	for _, field := range node.fields() {
		n, err := s.find(field.rt, field.tags)
		if err != nil && field.optional {
			continue
		}
		if err != nil {
			return fmt.Errorf("%s: %s", node, err)
		}
		if err := visit(s, n, marks); err != nil {
			return err
		}
	}
	marks[node] = permanent
	return nil
}
