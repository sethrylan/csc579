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
var kcpu, kio, c, l, m int

const seed int64 = 42

const usageMsg string = "λ K C L\n" +
	"λ = distribution of interarrival times\n" +
	"Kcpu = customers that the CPU queue may hold\n" +
	"Kio = customers that the IO queue may hold\n" +
	"C = customers served before the program terminates\n" +
	"L = 0–M/M/1 system, 1–CPU with I/O disks\n" +
	"M = 1–FCFS, 2–LCFS-NP, 3–SJF-NP, 4–Prio-NP, 5–Prio-P"

func init() {

	if len(os.Args) < 7 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}

	debugPtr := flag.Bool("debug", false, "a bool")
	flag.Parse()
	args := flag.Args()

	λ, _ = strconv.ParseFloat(args[0], 64)
	kcpu, _ = strconv.Atoi(args[1])
	kio, _ = strconv.Atoi(args[2])
	c, _ = strconv.Atoi(args[3])
	l, _ = strconv.Atoi(args[4])
	m, _ = strconv.Atoi(args[5])

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
	fmt.Printf("Kcpu =    %d\n", kcpu)
	fmt.Printf("Kcpu =    %d\n", kio)
	fmt.Printf("C =    %d\n", c)
	fmt.Printf("L =    %d\n", l)

	if m > 5 || m < 1 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}

	switch l {
	case 1:
		mm1kSimulation()
	case 2:
		// TODO:
		os.Exit(1)
	default:
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}

	mm1k.P2Question1(seed)
}
func mm1kSimulation() {
	completes, rejects := mm1k.Simulate(λ, µ, mm1k.QueueMakers[m-1](kcpu), c, seed)
	sorted := append(rejects, completes...)
	sort.Sort(mm1k.ByID(sorted))
	totalEvents := sorted[len(sorted)-1].ID + 1

	sampleServiceTimeMean := mm1k.Mean(completes, mm1k.Service)
	//TODO: run simpulation 30 to compute confidence intervals:
	//  sampleServiceTimeVariance = mm1k.StdDev(completes, sampleServiceTimeMean)

	fmt.Printf("Master clock =                   %.3f\n", completes[len(completes)-1].Departure)
	fmt.Printf("CLR (Analytical) =               %.3f\n", mm1k.AnalyticalCLR(λ, kcpu))
	fmt.Printf("CLR (Empirical; X/N = %d/%d) =   %.3f\n", len(rejects), totalEvents, mm1k.EmpiricalCLR(len(rejects), totalEvents))
	fmt.Printf("Mean Service Time (S̄) =          %.3f\n", sampleServiceTimeMean)
	fmt.Printf("Mean Wait Time (W̄) =             %.3f\n", mm1k.Mean(completes, mm1k.Wait))
	fmt.Printf("Analytical Wait Time (W̄) =       %.3f\n", mm1k.AnalyticalWaitTime(λ, kcpu))

	for _, customer := range sorted {
		// L, L + 1, L + 10, and L + 11
		if customer.ID == l || customer.ID == l+1 || customer.ID == l+10 || customer.ID == l+11 {
			mm1k.PrintCustomer(customer)
		}
	}

}
