package rand

import (
	"math/rand"
	"time"

	"gitlab.liebaopay.com/mikezhang/goplus/number"
)

// Randint 区间 [m,n] 中随机一个值
func Randint(m, n int) int {
	if m > n {
		panic("m is more than n")
	}
	if m == n {
		return m
	}
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n+1-m) + m
}

// Randints 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func Randints(m, n, k int, isDistinct bool) []int {
	if m > n {
		panic("m is more than n")
	}
	if k <= 0 {
		panic("k is less than zero")
	}
	if k > (n - m + 1) {
		panic("k is less than n - m")
	}
	if k == (n-m+1) && isDistinct {
		return number.Xrange(m, n+1, 1)
	}

	result := []int{}
	selected := map[int]bool{}
	for {
		if len(result) >= k {
			break
		}
		randomNum := Randint(m, n)
		if isDistinct {
			if _, ok := selected[randomNum]; ok {
				continue
			}
		}
		result = append(result, randomNum)
		if isDistinct {
			selected[randomNum] = true
		}
	}
	return result
}

// WeightSample 根据权重列表采样, 返回index
func WeightSample(weights []int) int {
	if len(weights) == 0 {
		panic("weights is blank")
	}
	totalWeight := number.Sum(weights)
	randNum := Randint(1, totalWeight)
	startNum := 0
	index := 0
	for i, w := range weights {
		if randNum <= startNum+w {
			index = i
			break
		}
		startNum += w
	}
	return index
}

// WeightSamples 根据权重列表批量采样, 返回index列表
func WeightSamples(weights []int, k int, isDistinct bool) []int {
	if len(weights) == 0 {
		panic("weights is blank")
	}
	if k > len(weights) {
		panic("k is more than weights length")
	}
	if k == len(weights) && isDistinct {
		return number.Xrange(0, len(weights), 1)
	}

	result := []int{}
	selected := map[int]bool{}
	for {
		if len(result) >= k {
			break
		}
		sampleIndex := WeightSample(weights)
		if isDistinct {
			if _, ok := selected[sampleIndex]; ok {
				continue
			}
		}
		result = append(result, sampleIndex)
		if isDistinct {
			selected[sampleIndex] = true
		}
	}
	return result
}
