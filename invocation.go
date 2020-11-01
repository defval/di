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

func validateInvocation(fn function) bool {
	if fn.NumOut() == 0 {
		return true
	}
	if fn.NumOut() == 1 && isError(fn.Out(0)) {
		return true
	}
	return false
}
