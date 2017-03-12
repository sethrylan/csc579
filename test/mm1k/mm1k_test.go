package mm1k_test

import (
  "testing"
  "mm1k"
)

func TestSimulate(t *testing.T) {
  λ := 0.5
  µ := 1.0
  k := 1
  c := 10
  completes, _ := mm1k.Simulate(λ, µ, k, c, 1)
  if (len(completes) != c) {
    t.Errorf("Expected %d completes, got %d", c, len(completes))
  }
}
