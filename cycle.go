package di

import (
	"errors"
)

const (
	temporary = 1
	permanent = 2
)

// used depth-first topological sort algorithm
func checkCycles(c *Container, param parameter) error {
	var marks = map[key]int{}
	provider, err := param.ResolveProvider(c)
	if err != nil {
		return err
	}
	if err := visit(c, provider, marks); err != nil {
		return err
	}
	return nil
}

func visit(c *Container, provider provider, marks map[key]int) error {
	pid := key{provider.Type(), provider.Name()}
	if marks[pid] == permanent {
		return nil
	}
	if marks[pid] == temporary {
		return errors.New("cycle detected") // todo: improve message
	}
	marks[pid] = temporary
	for _, param := range provider.ParameterList() {
		paramProvider, err := param.ResolveProvider(c)
		if _, ok := err.(errParameterProviderNotFound); ok && param.optional {
			continue
		}
		if err != nil {
			switch err.(type) {
			case errParameterProviderNotFound:
				return errDependencyNotFound{pid, param.Key()}
			default:
				return err
			}
		}
		if err := visit(c, paramProvider, marks); err != nil {
			return err
		}
	}
	marks[pid] = permanent
	return nil
}
