package ptr

// String returns a pointer to the provided string value.
func String(s string) *string {
	return &s
}
