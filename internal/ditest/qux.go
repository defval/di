package ditest

// Qux
type Qux struct {
	fooer Fooer
}

// NewQux
func NewQux(foo Fooer) *Qux {
	return &Qux{
		fooer: foo,
	}
}

func (q *Qux) Fooer() Fooer { return q.fooer }
