package nats

import (
	"math/rand"
	"time"
)

func generateRandomNumbers() []int32 {
	rand.Seed(time.Now().UnixNano())
	numSet := make(map[int32]struct{})
	for len(numSet) < 5 {
		n := int32(rand.Intn(49) + 1)
		numSet[n] = struct{}{}
	}
	numbers := make([]int32, 0, 5)
	for n := range numSet {
		numbers = append(numbers, n)
	}
	return numbers
}
