package main

import(
  "fmt"
  "net/http"
  "time"
)

var target1 = "http://exampleworkload-not-proxied:80/?n="
var target2 = "http://exampleworkload-proxied:8080/?n="

// Roughly corresponding to 50ms, 100ms, 200ms, 500ms, 1s
var n_array = []int{1, 19, 200, 300, 790, 1490}
var requests = 100

func main() {
  policy := waitpolicy_constant
  results2 := load_test(n_array, requests, policy, target2)
  fmt.Println("Results for proxied target: ", results2)

}

// This function takes care of making requests to the target server,
// measuring the time it takes to get a response and storing it in a 2D array.
// The time to wait between requests is determined by the waitpolicy function.
func load_test(
  n_array []int,
  requests int,
  waitpolicy func()time.Duration,
  target string,
  ) [][]float64 {
  results := make([][]float64, len(n_array))

  fmt.Println("Testing different values of n for ", requests, " requests")
  fmt.Println("n time")

  for i, n := range n_array {
    url := target + fmt.Sprintf("%d", n)
    results[i] = make([]float64, requests)  // Each row is a different n
    fmt.Printf("Starting %d\n", n)

    for ii := 0; ii < requests; ii++ {
      start := time.Now()
      _, err := http.Get(url)
      end := time.Now()

      results[i][ii] = end.Sub(start).Seconds()  // Each element is a different request

      if err != nil {
        fmt.Println("Error: ", err)
        panic(err)
      }
      time.Sleep(waitpolicy())
    }
  }
  return results
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
