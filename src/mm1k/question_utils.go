package mm1k

import (
	"fmt"
	"math"
	"reflect"
	"sort"
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

// AnalyticalWaitTime returns a calculated estimate of W̄
func AnalyticalWaitTime(ρ float64, K int) float64 {
	k := float64(K)
	p0 := (1 - ρ) / (1 - math.Pow(ρ, k+1))
	pK := math.Pow(ρ, k) * p0
	effLambda := (1 - pK) * ρ
	L := ρ * (1 - (k+1)*math.Pow(ρ, k) + k*math.Pow(ρ, k+1)) * p0 / math.Pow((1-ρ), 2)
	W := L / effLambda
	return W
}

type field func(c Customer) float64

// Service is a field for sorting customer.
func Service(c Customer) float64 {
	return c.Service
}

// Wait is a field for sorting customer.
func Wait(c Customer) float64 {
	return c.Start - c.Arrival
}

// Departure is a field for sorting customer.
func Departure(c Customer) float64 {
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

// MeanFloats calculates the mean for a list of floats
func MeanFloats(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))
}

// StdDev implementation
// See examples on https://github.com/ae6rt/golang-examples/blob/master/goeg/src/statistics_ans/statistics.go
func StdDev(numbers []float64, mean float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += math.Pow(number-mean, 2)
	}
	variance := total / float64(len(numbers)-1)
	return math.Sqrt(variance)
}

// RemoveFirstNByDeparture removes the first n deparatures in a list of customers
func RemoveFirstNByDeparture(customers []Customer, n int) []Customer {
	sort.Sort(byDeparture(customers))
	return customers[n:]
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func getType(myvar interface{}) (s string) {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		s = "*" + t.Elem().Name()
	} else {
		s = t.Name()
	}
	return
}
