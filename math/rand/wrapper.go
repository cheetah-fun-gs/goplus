package rand

import (
	"time"
)

var (
	defaultRand = New(time.Now().UnixNano())
)

// Randint 区间 [m,n] 中随机一个值
func Randint(m, n int) (int, error) {
	return defaultRand.Randint(m, n)
}

// MustRandint 区间 [m,n] 中随机一个值
func MustRandint(m, n int) int {
	return defaultRand.MustRandint(m, n)
}

// Randints 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func Randints(m, n, k int, isDistinct bool) ([]int, error) {
	return defaultRand.Randints(m, n, k, isDistinct)
}

// MustRandints 区间 [m,n] 中随机 k 个值, isDistinct 是否允许重复
func MustRandints(m, n, k int, isDistinct bool) []int {
	return defaultRand.MustRandints(m, n, k, isDistinct)
}

// WeightSample 根据权重列表采样, weights为样本列表的权重列表
// 返回采样样本的index, index从0开始
func WeightSample(weights []int) (int, error) {
	return defaultRand.WeightSample(weights)
}

// MustWeightSample 根据权重列表采样, weights为样本列表的权重列表
// 返回采样样本的index, index从0开始
func MustWeightSample(weights []int) int {
	return defaultRand.MustWeightSample(weights)
}

// WeightSamples 根据权重列表采样, weights为样本列表的权重列表, k 为样本数
// 返回采样样本的index列表, index从0开始
func WeightSamples(weights []int, k int, isDistinct bool) ([]int, error) {
	return defaultRand.WeightSamples(weights, k, isDistinct)
}

// MustWeightSamples 根据权重列表采样, weights为样本列表的权重列表, k 为样本数
// 返回采样样本的index列表, index从0开始
func MustWeightSamples(weights []int, k int, isDistinct bool) []int {
	return defaultRand.MustWeightSamples(weights, k, isDistinct)
}

// ProbSamples 根据概率列表采样, probs为样本列表的概率列表, k 为重复次数
// 返回采样样本的index列表, index从0开始
func ProbSamples(probs []float64, k int, isDistinct bool) ([]int, error) {
	return defaultRand.ProbSamples(probs, k, isDistinct)
}

// MustProbSamples 根据概率列表采样, probs为样本列表的概率列表, k 为重复次数
// 返回采样样本的index列表, index从0开始
func MustProbSamples(probs []float64, k int, isDistinct bool) []int {
	return defaultRand.MustProbSamples(probs, k, isDistinct)
}
