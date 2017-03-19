package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mm1k"
	"os"
	"sort"
	"strconv"
)

var λ, µ float64
var k, c, l, m int

const seed int64 = 42

const usageMsg string = "λ K C L\n" +
	"λ = distribution of interarrival times\n" +
	"K = customers that the queue may hold\n" +
	"C = customers served before the program terminates\n" +
	"L = customers of interest; 1 < L < C\n" +
	"M = 1–FCFS, 2–LCFS-NP, 3–SJF-NP, 4–Prio-NP, 5–Prio-P"

func init() {

	if len(os.Args) < 6 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}

	debugPtr := flag.Bool("debug", false, "a bool")
	flag.Parse()
	args := flag.Args()

	λ, _ = strconv.ParseFloat(args[0], 64)
	k, _ = strconv.Atoi(args[1])
	c, _ = strconv.Atoi(args[2])
	l, _ = strconv.Atoi(args[3])
	m, _ = strconv.Atoi(args[4])

	µ = 1.0

	log.SetFlags(log.Lshortfile)
	if *debugPtr {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	fmt.Printf("λ =    %.3f\n", λ)
	fmt.Printf("µ =    %.3f\n", µ)
	fmt.Printf("K =    %d\n", k)
	fmt.Printf("C =    %d\n", c)
	fmt.Printf("L =    %d\n", l)

	var queue mm1k.Queue
	switch m {
	case 1:
		queue = mm1k.NewFIFO(k)
	case 2:
		queue = mm1k.NewLIFO(k)
	case 3:
		queue = mm1k.NewSJF(k)
	case 4:
		queue = mm1k.NewPriority(k, 4, false)
	case 5:
		queue = mm1k.NewPriority(k, 4, true)
	default:
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}
	completes, rejects := mm1k.Simulate(λ, µ, queue, c, seed)

	sorted := append(rejects, completes...)
	sort.Sort(mm1k.ByID(sorted))
	totalEvents := sorted[len(sorted)-1].ID + 1

	fmt.Printf("Master clock =                   %.3f\n", completes[len(completes)-1].Departure)
	fmt.Printf("CLR (Analytical) =               %.3f\n", mm1k.AnalyticalCLR(λ, k))
	fmt.Printf("CLR (Empirical; X/N = %d/%d) =   %.3f\n", len(rejects), totalEvents, mm1k.EmpiricalCLR(len(rejects), totalEvents))
	fmt.Printf("Mean Service Time (S̄) =          %.3f\n", mm1k.Mean(completes, mm1k.Service))
	fmt.Printf("Mean Wait Time (W̄) =             %.3f\n", mm1k.Mean(completes, mm1k.Wait))

	for _, customer := range sorted {
		// L, L + 1, L + 10, and L + 11
		if customer.ID == l || customer.ID == l+1 || customer.ID == l+10 || customer.ID == l+11 {
			mm1k.PrintCustomer(customer)
		}
	}

	// mm1k.P2Question1(seed)
}
