package state

func Default() *Provider {
	return NewProvider(func() State {
		return NewDefaultState()
	})
}
