package utilities

// ArrayContains returns whether the specified array contains the specified value
func ArrayContains[T comparable](array []T, val T) bool {
	for _, v := range array {
		if v == val {
			return true
		}
	}

	return false
}

func AddToArray[V comparable, T ~[]V](arr *T, val V) bool {
	for _, v := range *arr {
		if v == val {
			return false
		}
	}
	*arr = append(*arr, val)
	return true
}
