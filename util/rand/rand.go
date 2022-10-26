package rand

import (
	"github.com/seehuhn/mt19937"
	"math/rand"
	"sync"
	"time"
)

var randPool = sync.Pool{
	New: func() interface{} {
		rng := rand.New(mt19937.New())
		rng.Seed(time.Now().UnixNano())
		return rng
	},
}

//[0,n)
func RandomInt64(n int64) int64 {
	rng := randPool.Get()
	result := rng.(*rand.Rand).Int63n(n)
	randPool.Put(rng)
	return result
}

//[0,n)
func RandomInt32(v int32) int32 {
	rng := randPool.Get()
	result := rng.(*rand.Rand).Int31n(v)
	randPool.Put(rng)
	return result
}

//[0,n)
func RandomInt(n int) int {
	rng := randPool.Get()
	result := rng.(*rand.Rand).Int()
	randPool.Put(rng)
	return result % n
}

func RandomFloat64() float64 {
	rng := randPool.Get()
	result := rng.(*rand.Rand).Float64()
	randPool.Put(rng)
	return result
}

func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}
	rng := randPool.Get().(*rand.Rand)
	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(rng.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(rng.Int31n(int32(i + 1)))
		swap(i, j)
	}
	randPool.Put(rng)
}

//[start,end)
func RandomRangeInt64(start, end int64) int64 {
	return RandomInt64(end-start) + start
}

//[start,end)
func RandomRangeInt32(start, end int32) int32 {
	return RandomInt32(end-start) + start
}

//[start,end)
func RandomRangeInt(start, end int) int {
	return RandomInt(end-start) + start
}
