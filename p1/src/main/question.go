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
func CLR(ρ float64, K int) float64 {
	k := float64(K)
	return ((1 - ρ) * math.Pow(ρ, k)) / (1 - math.Pow(ρ, (k+1)))
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
