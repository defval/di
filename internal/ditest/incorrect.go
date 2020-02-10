package ditest

// ConstructorWithoutResult
func ConstructorWithoutResult() {

}

// ConstructorWithManyResults
func ConstructorWithManyResults() (*Foo, *Bar, error) {
	return &Foo{}, &Bar{}, nil
}

// ConstructorWithIncorrectResultError
func ConstructorWithIncorrectResultError() (*Foo, *Bar) {
	return &Foo{}, &Bar{}
}
