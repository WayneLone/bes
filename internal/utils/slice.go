package utils

func MoveSliceElemTo[T comparable](s []T, pos int, findElemFunc func(T) bool) []T {
	if len(s) < 2 {
		return s
	}
	index := -1
	for i, item := range s {
		if findElemFunc(item) {
			index = i
			break
		}
	}
	if index == -1 {
		return s
	}
	val := s[index]
	s = append(s[:index], s[index+1:]...)
	ns := make([]T, pos+1)
	copy(ns, s[:pos])
	ns[pos] = val
	s = append(ns, s[pos:]...)
	return s
}

func EqualsSlice[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
