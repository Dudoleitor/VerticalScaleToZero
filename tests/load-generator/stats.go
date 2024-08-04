package main

import "slices"

// This function computes the mean of a slice of float64s.
func mean(data []float64) float64 {
  sum := 0.0
  for _, value := range data {
    sum += value
  }
  return sum / float64(len(data))
}

// This function computes the variance of a slice of float64s, given the mean.
func variance(data []float64, mean float64) float64 {
  sum := 0.0
  for _, value := range data {
    sum += (value - mean) * (value - mean)
  }
  return sum / float64(len(data))
}

// This function computes the median of a slice of float64s.
func median(data []float64) float64 {
  // Sort the data
  sorted := make([]float64, len(data))
  copy(sorted, data)
  slices.Sort(sorted)

  // Finding the median
  if len(sorted) % 2 == 0 {
    return (sorted[len(sorted)/2] + sorted[len(sorted)/2-1]) / 2
  } else {
    return sorted[len(sorted)/2]
  }
}

// This function computes the mode of a slice of float64s.
func fashion(data []float64) float64 {
  // Count the frequency of each value
  freq := make(map[float64]int)
  for _, value := range data {
    freq[value]++
  }

  // Find the most frequent value
  max := 0
  var mode float64
  for key, value := range freq {
    if value > max {
      max = value
      mode = key
    }
  }
  return mode
}

// This function computes the 95th percentile of a slice of float64s.
func percentile(data []float64) float64 {
  // Sort the data
  sorted := make([]float64, len(data))
  copy(sorted, data)
  slices.Sort(sorted)

  // Finding the 95th percentile
  index := int(0.95 * float64(len(sorted)))
  return sorted[index]
}
