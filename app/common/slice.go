package common

func SliceMax(numbers ...uint64) uint64 {
	if len(numbers) == 0 {
		return 0
	}
	if len(numbers) == 1 {
		return numbers[0]
	}
	if len(numbers) == 2 {
		if numbers[0] < numbers[1] {
			return numbers[1]
		}
		return numbers[0]
	}
	max := numbers[0]
	for _, r := range numbers {
		if r > max {
			max = r
		}
	}
	return max
}
func SliceMin(numbers ...uint64) uint64 {
	if len(numbers) == 0 {
		return 0
	}
	if len(numbers) == 1 {
		return numbers[0]
	}
	if len(numbers) == 2 {
		if numbers[0] < numbers[1] {
			return numbers[0]
		}
		return numbers[1]
	}
	min := numbers[0]
	for _, r := range numbers {
		if r < min {
			min = r
		}
	}
	return min
}
