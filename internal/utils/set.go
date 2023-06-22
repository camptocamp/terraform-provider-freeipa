package utils

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func Compare[T constraints.Ordered](a, b T) int {
	if a == b {
		return 0
	}

	if a < b {
		return -1
	}

	return +1
}

func SetDiff[T constraints.Ordered](actual, desired []T) (toAdd, toRemove []T) {
	slices.Sort(actual)
	slices.Sort(desired)

	i := 0
	j := 0

	for i < len(actual) && j < len(desired) {
		switch Compare(actual[i], desired[j]) {
		case -1:
			toRemove = append(toRemove, actual[i])
			i++
		case 0:
			i++
			j++
		case 1:
			toAdd = append(toAdd, desired[j])
			j++
		}
	}

	if i < len(actual) {
		toRemove = append(toRemove, actual[i:]...)
	} else if j < len(desired) {
		toAdd = append(toAdd, desired[j:]...)
	}

	return
}
