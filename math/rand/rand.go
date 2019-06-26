package rand

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.liebaopay.com/mikezhang/goplus/number"
)

func randintInter(s rand.Source, m, n int) int {
	if m == n {
		return m
	}
	return rand.New(s).Intn(n+1-m) + m
}

// Randint2 返回异常
func Randint2(m, n int) (int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	if err != nil {
		return 0, err
	}
	return Randint(m, n), nil
}

// Randint 区间 [m,n] 中随机一个值
func Randint(m, n int) int {
	if m > n {
		panic("m is more than n")
	}
	s := rand.NewSource(time.Now().UnixNano())
	return randintInter(s, m, n)
}

// Randints2 返回异常版本
func Randints2(m, n, k int, isDistinct bool) ([]int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	if err != nil {
		return nil, err
	}
	return Randints(m, n, k, isDistinct), nil
}

// Randints 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func Randints(m, n, k int, isDistinct bool) []int {
	if m > n {
		panic("m is more than n")
	}
	if k <= 0 {
		panic("k is less than zero")
	}
	if k > (n-m+1) && isDistinct {
		panic("k is less than n - m")
	}
	if k == (n-m+1) && isDistinct {
		return number.Xrange(m, n+1, 1)
	}

	s := rand.NewSource(time.Now().UnixNano())

	result := []int{}
	selected := map[int]bool{}
	for {
		if len(result) >= k {
			break
		}
		randomNum := randintInter(s, m, n)
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

func weightSampleInter(s rand.Source, weights []int, totalWeight int) int {
	randNum := randintInter(s, 1, totalWeight)
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

// WeightSample2 返回异常版本
func WeightSample2(weights []int) (int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	if err != nil {
		return 0, err
	}
	return WeightSample(weights), nil
}

// WeightSample 根据权重列表采样, 返回index
func WeightSample(weights []int) int {
	if len(weights) == 0 {
		panic("weights is blank")
	}
	totalWeight := number.Sum(weights)
	if totalWeight == 0 {
		panic("totalWeight is zero")
	}
	s := rand.NewSource(time.Now().UnixNano())
	return weightSampleInter(s, weights, totalWeight)
}

// WeightSamples2 返回异常版本
func WeightSamples2(weights []int, k int, isDistinct bool) ([]int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	if err != nil {
		return nil, err
	}
	return WeightSamples(weights, k, isDistinct), nil
}

// WeightSamples 根据权重列表批量采样, 返回index列表
func WeightSamples(weights []int, k int, isDistinct bool) []int {
	if len(weights) == 0 {
		panic("weights is blank")
	}
	if k > len(weights) && isDistinct {
		panic("k is more than weights length")
	}
	totalWeight := number.Sum(weights)
	if totalWeight == 0 {
		panic("totalWeight is zero")
	}
	if k == len(weights) && isDistinct {
		return number.Xrange(0, len(weights), 1)
	}

	s := rand.NewSource(time.Now().UnixNano())

	result := []int{}
	selected := map[int]bool{}
	for {
		if len(result) >= k {
			break
		}
		sampleIndex := weightSampleInter(s, weights, totalWeight)
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

// ProbSamples2 返回异常版本
func ProbSamples2(probs []float64, k int, isDistinct bool) ([]int, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()
	if err != nil {
		return nil, err
	}
	return ProbSamples(probs, k, isDistinct), nil
}

// ProbSamples 根据概率列表采样,k 为重复次数 返回index列表
func ProbSamples(probs []float64, k int, isDistinct bool) []int {
	if len(probs) == 0 {
		panic("probs is blank")
	}
	if k > len(probs) && isDistinct {
		panic("k is more than probs length")
	}

	s := rand.NewSource(time.Now().UnixNano())

	result := []int{}
	selected := map[int]bool{}
	for i := 0; i < k; i++ {
		for sampleIndex, prob := range probs {
			if isDistinct {
				if _, ok := selected[sampleIndex]; ok {
					continue
				}
			}
			if rand.New(s).Float64() > prob {
				continue
			}
			result = append(result, sampleIndex)
			if isDistinct {
				selected[sampleIndex] = true
			}
		}
	}
	return result
}
