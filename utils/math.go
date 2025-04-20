package utils

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}

func Min2[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	}
	return b
}
