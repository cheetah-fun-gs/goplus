// Package number 浮点, 整形方法
package number

// MinIndex 最小值 和 位置
func MinIndex(queue []int) (int, int) {
	min := queue[0]
	index := 0
	for i, j := range queue[1:] {
		if j < min {
			min = j
			index = i + 1
		}
	}
	return min, index
}

// Min 最小值
func Min(queue []int) int {
	min, _ := MinIndex(queue)
	return min
}

// MaxIndex 最大值 和 位置
func MaxIndex(queue []int) (int, int) {
	max := queue[0]
	index := 0
	for i, j := range queue[1:] {
		if j > max {
			max = j
			index = i + 1
		}
	}
	return max, index
}

// Max 最小值
func Max(queue []int) int {
	max, _ := MaxIndex(queue)
	return max
}

// Sum 求和
func Sum(queue []int) int {
	sum := queue[0]
	for _, i := range queue[1:] {
		sum += i
	}
	return sum
}

// Xrange 类 python xrange [start, end)
func Xrange(start, end, step int) []int {
	r := []int{}
	for i := start; i < end; i += step {
		r = append(r, i)
	}
	return r
}

// Pow x 的 y 次方, math.Pow 存在精度问题
func Pow(x, y int) int {
	z := 1
	for i := 0; i < y; i++ {
		z = z * x
	}
	return z
}
