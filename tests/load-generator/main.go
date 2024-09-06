package main

import(
  "fmt"
  "net/http"
  "time"
)

type Results struct {
  response_times []float64
  cpuAvg float64
  cpuVar float64
}

var target = "http://exampleworkload-proxied:8080/?n="

// Roughly corresponding to 50ms, 100ms, 200ms, 500ms, 1s
var n_array = []int{1, 305, 510, 800, 1525, 3005}

// var n_array = []int{10, 4000}
var requests = 100

func main() {
  policy := waitpolicy_exponential
  results := load_test(n_array, requests, policy, target)
  print_results(results)
}

// This function takes care of making requests to the target server,
// measuring the time it takes to get a response and storing it in a 2D array.
// The time to wait between requests is determined by the waitpolicy function.
func load_test(
  n_array []int,
  requests int,
  waitpolicy func()time.Duration,
  target string,
  ) []Results {
  results := make([]Results, len(n_array))  // Each row is a different n (different results)

  fmt.Println("Testing different values of n for ", requests, " requests")

  for i, n := range n_array {
    url := target + fmt.Sprintf("%d", n)
    results[i].response_times = make([]float64, requests)  // Each row is a different n
    fmt.Printf("Starting %d\n", n)
    startMonitor()

    for ii := 0; ii < requests; ii++ {
      start := time.Now()
      _, err := http.Get(url)
      end := time.Now()

      results[i].response_times[ii] = end.Sub(start).Seconds()  // Each element is a different request

      if err != nil {
        fmt.Println("Error: ", err)
        panic(err)
      }
      time.Sleep(waitpolicy())
    }
    results[i].cpuAvg, results[i].cpuVar = stopMonitor()
  }
  return results
}

func print_results(res []Results) {
  for i, row := range res {
    fmt.Println("Results for n = ", n_array[i])
    fmt.Println("Response time:")
    mean := mean(row.response_times)
    fmt.Println("  Mean: ", mean)
    fmt.Println("  Variance: ", variance(row.response_times, mean))
    fmt.Println("  Median: ", median(row.response_times))
    fmt.Println("  Mode: ", fashion(row.response_times))
    fmt.Println("  95th percentile: ", percentile(row.response_times))
    fmt.Println("CPU request:")
    fmt.Println("  Mean: ", row.cpuAvg)
    fmt.Println("  Variance: ", row.cpuVar)
    fmt.Println()
  }
}

func run_n_test(target string) {
  fmt.Println("Testing different values of n for target: ", target)
  fmt.Println("n time")
  for i := 1; i <= 1500; i++ {
    url := target + fmt.Sprintf("%d", i)

    start := time.Now()
    _, err := http.Get(url)
    end := time.Now()


    if err != nil {
      fmt.Println("Error: ", err)
      break
    } else {
      fmt.Printf("%d %s\n", i, end.Sub(start))
    }
  }
}
