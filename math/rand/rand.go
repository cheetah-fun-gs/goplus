// Package rand 随机库加强
// 权重采样和概率采样说明: n = 样本集大小, k = 参数k
//     权重采样: 样本集按权重挑选一个样本
//     权重多重采样(不去重): 权重采样重复 k 次, 得到 k 个样本
//     权重多重采样(去重): 权重采样重复无限次, 直到得到 k 个不重复的样本
//     概率采样: 样本集每个样本按各自概率判断是否被选中, 得到 [0, n] 个样本
//     概率多重采样(不去重): 概率采样重复 k 次, 得到 [0, n*k] 个样本
//     概率多重采样(去重): 概率采样重复 k 次, 得到 [0, n*k] 个样本, 再去重
package rand

import (
	"fmt"
	"math/rand"

	"github.com/cheetah-fun-gs/goplus/number"
)

// Rand ...
type Rand struct {
	src rand.Source
	*rand.Rand
}

// New 指定种子
func New(seed int64) *Rand {
	src := rand.NewSource(seed)
	return &Rand{
		src:  src,
		Rand: rand.New(src),
	}
}

func randintInter(s rand.Source, m, n int) int {
	if m == n {
		return m
	}
	return rand.New(s).Intn(n+1-m) + m
}

// Randint 同wrapper
func (r *Rand) Randint(m, n int) (int, error) {
	if m > n {
		return 0, fmt.Errorf("m is more than n")
	}
	return randintInter(r.src, m, n), nil
}

// MustRandint 同wrapper
func (r *Rand) MustRandint(m, n int) int {
	v, err := r.Randint(m, n)
	if err != nil {
		panic(err)
	}
	return v
}

// Randints 同wrapper
func (r *Rand) Randints(m, n, k int, isDistinct bool) ([]int, error) {
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
		randomNum := randintInter(r.src, m, n)
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

// MustRandints 同wrapper
func (r *Rand) MustRandints(m, n, k int, isDistinct bool) []int {
	v, err := r.Randints(m, n, k, isDistinct)
	if err != nil {
		panic(err)
	}
	return v
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

// WeightSample 同wrapper
func (r *Rand) WeightSample(weights []int) (int, error) {
	if len(weights) == 0 {
		return 0, fmt.Errorf("weights is blank")
	}
	totalWeight := number.Sum(weights)
	if totalWeight == 0 {
		return 0, fmt.Errorf("totalWeight is zero")
	}
	return weightSampleInter(r.src, weights, totalWeight), nil
}

// MustWeightSample 同wrapper
func (r *Rand) MustWeightSample(weights []int) int {
	v, err := r.WeightSample(weights)
	if err != nil {
		panic(err)
	}
	return v
}

// WeightSamples 同wrapper
func (r *Rand) WeightSamples(weights []int, k int, isDistinct bool) ([]int, error) {
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
		sampleIndex := weightSampleInter(r.src, weights, totalWeight)
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

// MustWeightSamples 同wrapper
func (r *Rand) MustWeightSamples(weights []int, k int, isDistinct bool) []int {
	v, err := r.WeightSamples(weights, k, isDistinct)
	if err != nil {
		panic(err)
	}
	return v
}

// ProbSamples 同wrapper
func (r *Rand) ProbSamples(probs []float64, k int, isDistinct bool) ([]int, error) {
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
			if rand.New(r.src).Float64() > prob {
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

// MustProbSamples 同wrapper
func (r *Rand) MustProbSamples(probs []float64, k int, isDistinct bool) []int {
	v, err := r.ProbSamples(probs, k, isDistinct)
	if err != nil {
		panic(err)
	}
	return v
}
