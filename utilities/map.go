package utilities

func MapKeys[T comparable, U any](m map[T]U) []T {
	r := make([]T, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func MapValues[T comparable, U any](m map[T]U) []U {
	r := make([]U, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
