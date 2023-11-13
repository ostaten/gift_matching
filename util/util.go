package util

import (
	"math/rand"
	"time"
)

//alters the random seed when executed in close conjunction back to back so a different seed is used for the same time.Now().UnixNano()
var randSeed float64 = 0

/*
 $ Gets a random number of [min, max)
*/
func GetRandom(min int, max int) int {
	source := rand.NewSource(time.Now().UnixNano() + int64(randSeed))
	rng := rand.New(source)
	
	res := rng.Intn(min + max) + min
	randSeed++
	return res
}

/*
 $ Checks if a thing e exists inside a slice s, where e is of type T and s contains items of type T
 */
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
			if v == e {
					return true
			}
	}
	return false
}

