package ditest

// Bar
type Bar struct {
	foo *Foo
}

// NewBar
func NewBar(foo *Foo) *Bar {
	return &Bar{
		foo: foo,
	}
}

// CreateBarConstructor
func CreateBarConstructor(bar *Bar) func(foo *Foo) *Bar {
	return func(foo *Foo) *Bar {
		bar.foo = foo
		return bar
	}
}

func (b *Bar) Foo() *Foo { return b.foo }
