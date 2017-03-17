package mm1k_test

import (
	"io/ioutil"
	"log"
	"math"
	"mm1k"
	"testing"
)

var simTests = []struct {
	λ                   float64
	µ                   float64
	k                   int
	c                   int
	seed                int64
	expectedRejects     int
	expectedTime        float64
	expectedAverageWait float64
}{
	{0.10, 1.0, 10, 1000, 1., 0, 9763.68383, 1.04583},
	{0.30, 1.0, 10, 1000, 1., 0, 3255.77326, 1.31954},
	{0.50, 1.0, 10, 1000, 1., 0, 1955.29853, 1.84465},
	{0.70, 1.0, 10, 1000, 1., 8, 1408.33793, 2.63453},
	{0.90, 1.0, 10, 1000, 1., 45, 1136.25512, 4.13247},
	{0.75, 1.0, 20, 1000, 42, 0, 1258.953714, 4.91566},        // p1.task1
	{0.85, 1.0, 20, 1000, 42, 4, 1115.46965, 7.38695},         // p1.task1
	{0.95, 1.0, 20, 1000, 42, 47, 1057.682103, 10.73246},      // p1.task1
	{0.85, 1.0, 10, 1000, 42, 53, 1190.409935, 4.72365},       // p1.task2
	{0.85, 1.0, 20, 1000, 42, 4, 1115.469653, 7.38695},        // p1.task2
	{0.95, 1.0, 100, 100000, 42, 26, 105349.449478, 18.90901}, // p1.task4
}

func TestSimulate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	for _, tt := range simTests {
		completes, rejects := mm1k.Simulate(tt.λ, tt.µ, mm1k.NewFIFO(tt.k), tt.c, tt.seed)

		if len(completes) != tt.c {
			t.Errorf("Expected %d completes, got %d", tt.c, len(completes))
		}

		if len(rejects) != tt.expectedRejects {
			t.Errorf("Expected %d rejects, got %d", tt.expectedRejects, len(rejects))
		}

		finalTime := completes[len(completes)-1].Departure
		if diff := math.Abs(finalTime - tt.expectedTime); diff > .00001 {
			t.Errorf("Expected %f clock time, got %f", tt.expectedTime, finalTime)
		}

		w := mm1k.Mean(completes, mm1k.Wait)
		if diff := math.Abs(w - tt.expectedAverageWait); diff > .00001 {
			t.Errorf("Expected %f W̄, got %f", tt.expectedAverageWait, w)
		}

	}
}
