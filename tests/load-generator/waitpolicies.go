package main

import "time"

// This function implements a constant wait policy.
func waitpolicy_constant() time.Duration {
  rate := 300
  return time.Duration(rate) * time.Millisecond
}
