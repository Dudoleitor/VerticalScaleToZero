package main

import "slices"

type RollingMean struct {
  sum int64
  count int64
}

type RollingVariance struct {
  sum int64
  sumOfSquares int64
  count int64
}

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

func getRollingMean (d RollingMean) float64 {
  return float64(d.sum) / float64(d.count)
}
func updateRollingMean (d RollingMean, newValue int64) RollingMean {
  d.sum += newValue
  d.count++
  return d
}

func getRollingVariance (d RollingVariance) float64 {
  return (float64(d.sumOfSquares) - ((float64(d.sum)*float64(d.sum))/float64(d.count))) / (float64(d.count) - 1)
}
func updateRollingVariance (d RollingVariance, newValue int64) RollingVariance {
  d.sum += newValue
  d.sumOfSquares += newValue * newValue
  d.count++
  return d
}
