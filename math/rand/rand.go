// Package rand 随机方法
package rand

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cheetah-fun-gs/goplus/number"
)

func randintInter(s rand.Source, m, n int) int {
	if m == n {
		return m
	}
	return rand.New(s).Intn(n+1-m) + m
}

// Randint 区间 [m,n] 中随机一个值
func Randint(m, n int) (int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	return RandintWithSource(s, m, n)
}

// RandintWithSource 区间 [m,n] 中随机一个值, 给定种子
func RandintWithSource(s rand.Source, m, n int) (int, error) {
	if m > n {
		return 0, fmt.Errorf("m is more than n")
	}
	return randintInter(s, m, n), nil
}

// RandintsWithSource 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func RandintsWithSource(s rand.Source, m, n, k int, isDistinct bool) ([]int, error) {
	if m > n {
		return nil, fmt.Errorf("m is more than n")
	}
	if k <= 0 {
		return nil, fmt.Errorf("k is less than zero")
	}
	if k > (n-m+1) && isDistinct {
		return nil, fmt.Errorf("k is less than n - m")
	}
	if k == (n-m+1) && isDistinct {
		return number.Xrange(m, n+1, 1), nil
	}

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
	return result, nil
}

// Randints 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func Randints(m, n, k int, isDistinct bool) ([]int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	return RandintsWithSource(s, m, n, k, isDistinct)
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

// WeightSampleWithSource 根据权重列表采样, 返回index 给定种子
func WeightSampleWithSource(s rand.Source, weights []int) (int, error) {
	if len(weights) == 0 {
		return 0, fmt.Errorf("weights is blank")
	}
	totalWeight := number.Sum(weights)
	if totalWeight == 0 {
		return 0, fmt.Errorf("totalWeight is zero")
	}
	return weightSampleInter(s, weights, totalWeight), nil
}

// WeightSample 根据权重列表采样, 返回index
func WeightSample(weights []int) (int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	return WeightSampleWithSource(s, weights)
}

// WeightSamplesWithSource 根据权重列表批量采样, 返回index列表
func WeightSamplesWithSource(s rand.Source, weights []int, k int, isDistinct bool) ([]int, error) {
	if len(weights) == 0 {
		return nil, fmt.Errorf("weights is blank")
	}
	if k > len(weights) && isDistinct {
		return nil, fmt.Errorf("k is more than weights length")
	}
	totalWeight := number.Sum(weights)
	if totalWeight == 0 {
		return nil, fmt.Errorf("totalWeight is zero")
	}
	if k == len(weights) && isDistinct {
		return number.Xrange(0, len(weights), 1), nil
	}

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
	return result, nil
}

// WeightSamples 根据权重列表批量采样, 返回index列表
func WeightSamples(weights []int, k int, isDistinct bool) ([]int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	return WeightSamplesWithSource(s, weights, k, isDistinct)
}

// ProbSamplesWithSource 根据概率列表采样,k 为重复次数 返回index列表
func ProbSamplesWithSource(s rand.Source, probs []float64, k int, isDistinct bool) ([]int, error) {
	if len(probs) == 0 {
		return nil, fmt.Errorf("probs is blank")
	}
	if k > len(probs) && isDistinct {
		return nil, fmt.Errorf("k is more than probs length")
	}

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
	return result, nil
}

// ProbSamples 根据概率列表采样,k 为重复次数 返回index列表
func ProbSamples(probs []float64, k int, isDistinct bool) ([]int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	return ProbSamplesWithSource(s, probs, k, isDistinct)
}
