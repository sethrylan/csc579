package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	. "mm1k"
	"os"
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

	// Your simulation program will terminate once C customers have completed
	// service, where C is an input parameter. For initial conditions, assume that
	// at time t = 0 the system is empty. Draw a random number to decide when the
	// first arrival will occur.

	var customer Customer
	var rejected, completed <-chan Customer
	var rejects, completes []Customer
	rejected, completed = Run(
		NewExpDistribution(λ, seed),
		NewFIFOQueue(K),
		NewExpDistribution(λ, seed+1),
	)
	for len(completes) < C {
		select {
		case customer = <-rejected:
			rejects = append(rejects, customer)
			PrintCustomer("rejected ", customer)
		case customer = <-completed:
			completes = append(completes, customer)
			PrintCustomer("", customer)
		}
	}

	fmt.Printf("Master clock =         %.2f\n", customer.Departure)
	fmt.Printf("CLR =                  %.2f\n", CLR(λ, K))
	fmt.Printf("Average Service Time = %.2f\n", mean(completes, Service))
	fmt.Printf("Average waiting time = %.2f\n", mean(completes, Wait))

	// L, L + 1, L + 10, and L + 11
}

// func PrintCustomer(msg string, c Customer) {
// 	log.Printf("%sCustomer %d (%d)\n", msg, c.ID, c.Position)
// 	log.Printf("Arrival=%f\n", c.Arrival)
// 	log.Printf("Service=%f\n", c.Service)
// 	log.Printf("Start=%f\n", c.Start)
// 	log.Printf("Departure=%f\n", c.Departure)
// }

// Print the arrival time, service time, time of departure of customers, as well
// as the number of customers in the system immediately after the departure of
// each of these customers
func PrintCustomer(msg string, c Customer) {
	log.Printf("Customer %d %s (%d)\n", c.ID, msg, c.Position)
	log.Printf("Arrival =   %.3f\n", c.Arrival)
	log.Printf("Service =   %.3f\n", c.Service)
	log.Printf("Start =     %.3f\n", c.Start)
	log.Printf("Departure = %.3f\n", c.Departure)
}