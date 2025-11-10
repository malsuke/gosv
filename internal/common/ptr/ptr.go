package ptr

// String returns a pointer to the provided string value.
func Ptr[T any](v T) *T {
	return &v
}
