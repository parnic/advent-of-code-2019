package utilities

import "math"

func GCD[T Integer](a, b T) T {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func LCM[T Integer](nums ...T) uint64 {
	num := len(nums)
	if num == 0 {
		return 0
	} else if num == 1 {
		return uint64(nums[0])
	}

	ret := lcm(nums[0], nums[1])
	for i := 2; i < len(nums); i++ {
		ret = lcm(uint64(nums[i]), ret)
	}
	return ret
}

func lcm[T Integer](a, b T) uint64 {
	return uint64(a*b) / uint64(GCD(a, b))
}

func Min[T Number](nums ...T) T {
	numNums := len(nums)
	if numNums == 2 {
		return T(math.Min(float64(nums[0]), float64(nums[1])))
	}

	if numNums == 0 {
		return 0
	}

	least := nums[0]
	for i := 1; i < numNums; i++ {
		if nums[i] < least {
			least = nums[i]
		}
	}

	return least
}
