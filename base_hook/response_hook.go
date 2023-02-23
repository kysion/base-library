package base_hook

type ResponseFactoryHook[TResponse any] struct {
	responseFactoryHook func() TResponse
}

func (s *ResponseFactoryHook[TResponse]) RegisterResponseFactory(f func() TResponse) {
	s.responseFactoryHook = f
}

func (s *ResponseFactoryHook[TResponse]) FactoryMakeResponseInstance() TResponse {
	var result TResponse
	if s.responseFactoryHook == nil {
		return result
	}
	return s.responseFactoryHook()
}
