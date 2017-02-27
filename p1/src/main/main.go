package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	. "mm1k"
	"os"
	"sort"
	"strconv"
)

var λ float64
var K, C, L int

const seed int64 = 42

const usageMsg string = "λ K C L\n" +
	"λ = distribution of interarrival times\n" +
	"K = customers that the queue may hold\n" +
	"C = customers served before the program terminates\n" +
	"L = customers of interest; 1 < L < C\n"

func init() {

	if len(os.Args) < 5 {
		fmt.Printf("usage: %s %s", os.Args[0], usageMsg)
		os.Exit(1)
	}

	debugPtr := flag.Bool("debug", false, "a bool")
	flag.Parse()
	args := flag.Args()

	λ, _ = strconv.ParseFloat(args[0], 64)
	K, _ = strconv.Atoi(args[1])
	C, _ = strconv.Atoi(args[2])
	L, _ = strconv.Atoi(args[3])

	log.SetFlags(log.Lshortfile)
	if !*debugPtr {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	fmt.Printf("λ =    %.3f\n", λ)
	fmt.Printf("K =    %d\n", K)
	fmt.Printf("C =    %d\n", C)
	fmt.Printf("L =    %d\n", L)
	fmt.Printf("µ =    1\n")

	completes, rejects := Simulate(λ, K, C, seed)

	sorted := append(rejects, completes...)
	sort.Sort(ByID(sorted))
	totalEvents := sorted[len(sorted)-1].ID + 1

	fmt.Printf("Master clock =          %.2f\n", completes[len(completes)-1].Departure)
	fmt.Printf("CLR (Analytical) =      %.2f\n", AnalyticalCLR(λ, K))
	fmt.Printf("CLR (Empirical) =       %.2f\n", EmpiricalCLR(len(rejects), totalEvents))
	fmt.Printf("Mean Service Time (S̄) = %.2f\n", Mean(completes, Service))
	fmt.Printf("Mean Wait Time (W̄) =    %.2f\n", Mean(completes, Wait))

	for _, c := range sorted {
		// L, L + 1, L + 10, and L + 11
		if c.ID == L || c.ID == L+1 || c.ID == L+10 || c.ID == L+11 {
			PrintCustomer(c)
		}
	}

	Question1(seed)
}
