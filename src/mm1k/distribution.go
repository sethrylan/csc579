package mm1k

import (
	"math"
	"math/rand"
)

// A Distribution type provides a single function to get the next value
type Distribution interface {
	Get() (r float64)
}

// An ExpDistribution is an Exponential distribution
type ExpDistribution struct {
	generator *rand.Rand
	λ         float64
}

// NewExpDistribution returns a pointer to a new Exponential Distribution
func NewExpDistribution(λ float64, seed int64) (e *ExpDistribution) {
	return &ExpDistribution{rand.New(rand.NewSource(seed)), λ}
}

// Get returns the next exponentially distributed IID value
func (e *ExpDistribution) Get() float64 {
	return e.generator.ExpFloat64() / e.λ
}

func rand0() float64 {
	return rand.Float64()
}

// Expdev is used only for project requirements, but not run during simulation. The
// mm1k.ExpDistribution type is used instead to maintain seed state.
func Expdev(λ float64) float64 {
	return math.Log(1-rand0()) / (-λ)
}

// ParetoDistribution is an Pareto distribution
type ParetoDistribution struct {
	generator *rand.Rand
	α         float64
	k         int
	p         int
}

// NewParetoDistribution returns a pointer to a new Exponential Distribution
func NewParetoDistribution(α float64, k int, p int, seed int64) (d *ParetoDistribution) {
	return &ParetoDistribution{rand.New(rand.NewSource(seed)), α, k, p}
}

// Get returns the next bounded pareto distributed IID value
// See http://www.csee.usf.edu/~kchriste/tools/toolpage.html
func (d *ParetoDistribution) Get() float64 {
	var z float64 // Uniform random number from 0 to 1

	for {
		z = d.generator.Float64() // Get a uniform random number (0.0 < z < 1.0)
		if z != 0.0 {
			break
		}
	}

	// Generate the bounded Pareto rv using the inversion method
	// boundedRV = -(z*math.Pow(d.p, d.α) - z*math.Pow(d.k, d.α) - math.Pow(d.p, d.α)) / (math.Pow(d.p, d.α) * math.Pow(d.k, d.α))
	// boundedRV = math.Pow(boundedRV, (-1.0 / d.α))
	return math.Pow((math.Pow(float64(d.k), d.α) / (z*math.Pow(float64(d.k/d.p), d.α) - z + 1)), (1.0 / d.α))
}
