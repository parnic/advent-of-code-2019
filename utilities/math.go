package utilities

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
