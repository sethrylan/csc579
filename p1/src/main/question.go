package main

import (
	"math"
	"mm1k"
)

// Returns the Customer Loss Rate as a function of ρ and K
// ρ = λ/µ < 1
// ∴ µ = 1
// ∴ ρ = λ
// CLR = ((1-ρ)*ρ^K)/(1-ρ^(K+1))
func AnalyticalCLR(ρ float64, K int) float64 {
	k := float64(K)
	return ((1 - ρ) * math.Pow(ρ, k)) / (1 - math.Pow(ρ, (k+1)))
}

// Let N be the total number of customers that arrived to the system at the time
// the simulation ends (i.e., after the C-th customer completes service, C ≤ N).
// Let X be the number of customers denied service (lost) at the time the
// simulation ends. Let us define the customer loss rate (CLR) as:
// CLR = X/N
func EmpiricalCLR(x int, n int) float64 {
	return float64(x) / float64(n)
}

type field func(c mm1k.Customer) float64

func Service(c mm1k.Customer) float64 {
	return c.Service
}

func Wait(c mm1k.Customer) float64 {
	return c.Departure - c.Arrival
}

func mean(customers []mm1k.Customer, fn field) float64 {
	total := 0.0
	for _, c := range customers {
		total += fn(c)
	}
	return total / float64(len(customers))
}

// ByID implements sort.Interface for []Customer
type ByID []mm1k.Customer

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Let the queue capacity K = 20. Plot the CLR against the value of ρ,
// for ρ = 0.05 to ρ = 0.95, in increments of 0.10. Submit two graphs: one for
// C = 1000 and one for C = 100000.
func Question1() {

}
