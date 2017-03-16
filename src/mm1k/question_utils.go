package mm1k

import (
	"fmt"
	"math"
	"time"
)

// AnalyticalCLR returns the Customer Loss Rate as a function of ρ and K
// ρ = λ/µ < 1
// ∴ µ = 1
// ∴ ρ = λ
// CLR = ((1-ρ)*ρ^K)/(1-ρ^(K+1))
func AnalyticalCLR(ρ float64, K int) float64 {
	k := float64(K)
	return ((1 - ρ) * math.Pow(ρ, k)) / (1 - math.Pow(ρ, (k+1)))
}

// EmpiricalCLR returns the Customer Loss Rate:
// Let N be the total number of customers that arrived to the system at the time
// the simulation ends (i.e., after the C-th customer completes service, C ≤ N).
// Let X be the number of customers denied service (lost) at the time the
// simulation ends. Let us define the customer loss rate (CLR) as:
// CLR = X/N
func EmpiricalCLR(x int, n int) float64 {
	return float64(x) / float64(n)
}

type field func(c Customer) float64

// Service is a field for sorting customer.
func Service(c Customer) float64 {
	return c.Service
}

// Wait is a field for sorting customer.
func Wait(c Customer) float64 {
	return c.Departure - c.Arrival
}

// Mean calculates the mean for field fn in a list of customers
func Mean(customers []Customer, fn field) float64 {
	total := 0.0
	for _, c := range customers {
		total += fn(c)
	}
	return total / float64(len(customers))
}

// ByID implements sort.Interface for []Customer
type ByID []Customer

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}
