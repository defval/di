package di

import (
	"errors"
	"fmt"
)

const (
	temporary = 1
	permanent = 2
)

// used depth-first topological sort algorithm
func checkCycles(c *Container, param parameter) error {
	var marks = map[id]int{}
	provider, exists := param.ResolveProvider(c)
	if !exists {
		return errParameterProviderNotFound{param}
	}
	if err := visit(c, provider, marks); err != nil {
		return err
	}
	return nil
}

func visit(c *Container, provider provider, marks map[id]int) error {
	id := provider.ID()
	if marks[id] == permanent {
		return nil
	}
	if marks[id] == temporary {
		return errors.New("cycle detected") // todo: improve message
	}
	marks[id] = temporary
	for _, param := range provider.ParameterList() {
		paramProvider, exists := param.ResolveProvider(c)
		if !exists && param.optional {
			continue
		}
		if !exists {
			return fmt.Errorf("%s: dependency %s not exists in container", provider.ID(), param)
		}
		if err := visit(c, paramProvider, marks); err != nil {
			return err
		}
	}
	marks[id] = permanent
	return nil
}
