package utils

func MinIntegerOf(ints ...int) int {
	min := ints[0]
	for _, i := range ints {
		if min > i {
			min = i
		}
	}
	return min
}

func MaxIntegerOf(ints ...int) int {
	max := ints[0]
	for _, i := range ints {
		if max < i {
			max = i
		}
	}
	return max
}
