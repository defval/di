package di

// Invocation is a function whose signature looks like:
//
//		func StartServer(server *http.Server) error {
//			return server.ListenAndServe()
//		}
//
// Like a constructor invocation may have unlimited count of arguments and
// they will be resolved automatically. The invocation can return an optional error.
// Error will be returned as is.
type Invocation interface{}

// validateInvocation validates function.
func validateInvocation(fn function) bool {
	if fn.NumOut() == 0 {
		return true
	}
	if fn.NumOut() == 1 && isError(fn.Out(0)) {
		return true
	}
	return false
}

// parseInvocationParameters parses invocation and returns slice of nodes.
func parseInvocationParameters(fn function, s schema) (params []*node, err error) {
	for i := 0; i < fn.NumIn(); i++ {
		in := fn.Type.In(i)
		node, err := s.find(in, Tags{})
		if err != nil {
			return nil, err
		}
		params = append(params, node)
	}
	return params, nil
}
