package di

import "fmt"

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
