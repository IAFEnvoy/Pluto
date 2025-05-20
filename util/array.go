package util

func Contains[T comparable](slice []T, target T) bool {
	set := make(map[T]struct{})
	for _, v := range slice {
		set[v] = struct{}{}
	}
	_, ok := set[target]
	return ok
}

func ConcatMultipleSlices[T any](slices [][]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]T, totalLen)

	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}
