package mm1k_test

import (
  "testing"
  "mm1k"
)

func TestSimulate(t *testing.T) {
  λ := 0.5
  µ := 1.0
  K := 1
  C := 10
  completes, _ := mm1k.Simulate(λ, µ, mm1k.NewFIFOQueue(K), C, 1)
  if (len(completes) != C) {
    t.Errorf("Expected %d completes, got %d", C, len(completes))
  }
}
