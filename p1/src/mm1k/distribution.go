package mm1k

import (
	"math"
	"math/rand"
)

type Distribution interface {
	Get() (r float64)
}

type ExpDistribution struct {
	generator *rand.Rand
	λ         float64
}

func NewExpDistribution(λ float64, seed int64) (e *ExpDistribution) {
	return &ExpDistribution{rand.New(rand.NewSource(seed)), λ}
}

func (e *ExpDistribution) Get() float64 {
	return e.generator.ExpFloat64() / e.λ
}

func rand0() float64 {
	return rand.Float64()
}

/*
 *  Note: this function is used for project requirements, but not run during simulation. The
 *  mm1k.ExpDistribution type is used instead to maintain seed state.
 */
func Expdev(λ float64) float64 {
	return math.Log(1-rand0()) / (-λ)
}
