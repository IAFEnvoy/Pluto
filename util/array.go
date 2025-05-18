package util

func Contains[T comparable](slice []T, target T) bool {
	set := make(map[T]struct{})
	for _, v := range slice {
		set[v] = struct{}{}
	}
	_, ok := set[target]
	return ok
}
