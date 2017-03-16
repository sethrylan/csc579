package mm1k

import (
	"fmt"
	"math"
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

// Question1 : Let the queue capacity K = 20. Plot the CLR against the value of ρ,
// for ρ = 0.05 to ρ = 0.95, in increments of 0.10. Submit two graphs: one for
// C = 1000 and one for C = 100000.
func Question1(seed int64) {
	K := 20
	for _, C := range []int{1000, 100000} {
		for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
			fmt.Printf("%f, %d, %d | ", ρ, K, C)
			completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
			sorted := append(rejects, completes...)
			totalEvents := sorted[len(sorted)-1].ID + 1
			sort.Sort(ByID(sorted))
			// fmt.Printf("X/N = %d/%d\n", len(rejects), sorted[len(sorted)-1].ID+1)
			fmt.Printf("CLR (Empirical) = %.4f\n", EmpiricalCLR(len(rejects), totalEvents))
		}
	}
}

// Question2 : Now let us fix ρ = 0.85. Plot the CLR against the value of the queue
// capacity K, as K increases from 10 to 100 in increments of 10. Again, submit
// two graphs: one for C = 1000 and one for C = 100000.
func Question2(seed int64) {
	ρ := 0.85
	for _, C := range []int{1000, 100000} {
		for K := 10; K <= 100; K += 10 {
			fmt.Printf("%f, %d, %d | ", ρ, K, C)
			completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
			sorted := append(rejects, completes...)
			totalEvents := sorted[len(sorted)-1].ID + 1
			sort.Sort(ByID(sorted))
			fmt.Printf("CLR (Empirical) = %.4f\n", EmpiricalCLR(len(rejects), totalEvents))
		}
	}
}

// Question3 : Let K = 20. For C = 100000 and for ρ = 0.05 to 0.95 (in increments of 0.10),
// plot the simulation and analytical values of CLR on the same graph.
func Question3(seed int64) {
	K := 20
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("%f, %d, %d | ", ρ, K, C)
		completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		sorted := append(rejects, completes...)
		totalEvents := sorted[len(sorted)-1].ID + 1
		sort.Sort(ByID(sorted))
		fmt.Printf("CLR (Empirical, Analytical) =  (%.4f, %.4f\n",
			EmpiricalCLR(len(rejects), totalEvents),
			AnalyticalCLR(ρ, K))
	}
}

// Question4 : Let us set K = 100 and C = 100000. Compute the average waiting time W of the
// C customers that have received service at the end of the simulation (i.e.,
// ignore any lost customers or customers waiting in the queue when the
// simulation ends). Plot W against the value of ρ for ρ = 0.05 to 0.95.
func Question4(seed int64) {
	K := 100
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("%f, %d, %d | ", ρ, K, C)
		completes, _ := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))
	}
}

// Question5 : Let us again set K = 40 and C = 100000. Time the running time of your
// simulation for ρ = 0.05 to 0.95. Plot the running time against the value of
// ρ. Note: turn off I/O
func Question5(seed int64) {
	K := 40
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		func() {
			defer timeTrack(time.Now(), fmt.Sprintf("\nρ = %.2f", ρ))
			Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		}()
	}
	fmt.Printf("\n")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}
