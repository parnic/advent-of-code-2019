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

// CopyMap returns a copy of the passed-in map. Note: currently only works if [U]
// is not a map or slice.
func CopyMap[T comparable, U any](m map[T]U) map[T]U {
	r := make(map[T]U)
	for k, v := range m {
		r[k] = v
	}
	return r
}
