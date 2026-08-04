
func findMissingElements(nums []int) []int {
    mp := map[int]struct{}{}

	max := -1
	min := 101

	for _, num := range(nums) {
		mp[num] = struct{}{}
		if(num < min) {
			min = num
		}
		if(num > max) {
			max = num
		}
	}

	res := make([]int, 0, len(nums))
	for i := min; i <= max; i++ {
		if _, exists := mp[i]; !exists {
			res = append(res, i)
		}
	}
	return res
}
