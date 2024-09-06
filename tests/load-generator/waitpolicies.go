package main

import "time"
import "math"
import "math/rand"

// This function implements a constant wait policy.
func waitpolicy_constant() time.Duration {
  rate := 100
  return time.Duration(rate) * time.Millisecond
}

func waitpolicy_exponential() time.Duration {
  mean := 1200
  rnd := rand.Float64()
  waitTime := -1 * float64(mean) * math.Log(1-rnd)

  return time.Duration(waitTime) * time.Millisecond
}
