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
	if *debugPtr {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {

	fmt.Printf("λ = %.3f\n", λ)
	fmt.Printf("K = %d\n", K)
	fmt.Printf("C = %d\n", C)
	fmt.Printf("L = %d\n", L)

	// Your simulation program will terminate once C customers have completed service, where C is an input
	// parameter. For initial conditions, assume that at time t = 0 the system is empty. Draw a random number
	// to decide when the first arrival will occur, and then start your simulation by locating the first event, etc., as
	// we discussed in class.
	var rejected, completed <-chan Customer
	rejected, completed = Run(
		NewExpDistribution(λ, seed),
		NewFIFOQueue(K),
		NewExpDistribution(λ, seed+1),
	)
	var cus Customer
	serviced := 0
	for serviced < C {
		select {
		case cus = <-rejected:
			PrintCustomer("rejected ", cus)
		case cus = <-completed:
			serviced += 1
			PrintCustomer("", cus)
		}
	}

	fmt.Printf("Master clock = TODO\n")
	fmt.Printf("CLR = TODO\n")
	fmt.Printf("Average Service Time = TODO\n")
	fmt.Printf("Average waiting time = TODO\n")

	// the arrival time, service time, time of departure of customers L, L + 1, L + 10, and L + 11, as well as
	// the number of customers in the system immediately after the departure of each of these customers; if
	// any of these customers was not accepted for service (lost), set its departure time to its arrival time

}

func PrintCustomer(msg string, c Customer) {
	fmt.Printf("%sCustomer %d (%d)\n", msg, c.ID, c.Position)
	fmt.Printf("Arrival=%f\n", c.Arrival)
	fmt.Printf("Service=%f\n", c.Service)
	fmt.Printf("Start=%f\n", c.Start)
	fmt.Printf("Completion=%f\n", c.Completion)

}
