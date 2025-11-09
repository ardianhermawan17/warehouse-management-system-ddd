package pkg

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}

// Val returns the value of a pointer, or the zero value if nil
func Val[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
