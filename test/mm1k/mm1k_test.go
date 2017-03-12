package mm1k_test

import (
  "testing"
  "mm1k"
)

var simTests = []struct {
  λ        float64
  µ        float64
  k        int
  c         int
  seed    int64
  expectedRejects int
}{
  {0.1, 1.0, 10, 1000, 1, 0},
  {0.3, 1.0, 10, 1000, 1, 0},
  {0.5, 1.0, 10, 1000, 1, 0},
  {0.7, 1.0, 10, 1000, 1, 8},
  {0.9, 1.0, 10, 1000, 1, 45},
}

func TestSimulate(t *testing.T) {
  for _, tt := range simTests {
    completes, rejects := mm1k.Simulate(tt.λ, tt.µ, mm1k.NewFIFOQueue(tt.k), tt.c, tt.seed)
    if (len(completes) != tt.c) {
      t.Errorf("Expected %d completes, got %d", tt.c, len(completes))
    }
    if (len(rejects) != tt.expectedRejects) {
      t.Errorf("Expected %d completes, got %d", tt.expectedRejects, len(rejects))
    }
  }
}
