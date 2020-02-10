package ditest

import "github.com/goava/di"

// Baz
type Baz struct {
	foo *Foo
	bar *Bar
}

// NewBaz
func NewBaz(foo *Foo, bar *Bar) *Baz {
	return &Baz{
		foo: foo,
		bar: bar,
	}
}

// BazParameters
type BazParameters struct {
	di.Parameter // todo: remove parameters

	Foo *Foo `di:""`
	Bar *Bar `di:"optional"`
}

// NewBazFromParameters
func NewBazFromParameters(params BazParameters) *Baz {
	return &Baz{
		foo: params.Foo,
		bar: params.Bar,
	}
}

func (b *Baz) Foo() *Foo { return b.foo }
func (b *Baz) Bar() *Bar { return b.bar }
