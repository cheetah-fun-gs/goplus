package rand

import "testing"

func Test_Randint(t *testing.T) {
	Randint(-100, -20)
	Randint(15, 20)
	Randint(-10, 20)
}

func Test_Randints(t *testing.T) {
	Randints(-100, -20, 5, false)
	Randints(15, 20, 1, false)
	Randints(-10, 20, 15, true)
}

func Test_WeightSample(t *testing.T) {
	WeightSample([]int{15, 23, 44, 5})
}

func Test_WeightSamples(t *testing.T) {
	WeightSamples([]int{15, 23, 44, 5}, 3, false)
	WeightSamples([]int{15, 23, 44, 5}, 3, true)
	WeightSamples([]int{15, 23, 44, 5}, 1, false)
}
