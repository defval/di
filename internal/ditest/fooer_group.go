package ditest

// FooerGroup
type FooerGroup struct {
	fooers []Fooer
}

// NewFooerGroup
func NewFooerGroup(fooers []Fooer) *FooerGroup {
	return &FooerGroup{fooers: fooers}
}

func (g *FooerGroup) Fooers() []Fooer {
	return g.fooers
}
