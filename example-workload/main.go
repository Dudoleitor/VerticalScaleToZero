package main

import (
  "fmt"
  "time"
  "net/http"
  "strconv"
)

var testLength = 2
var numToTest = 100000000

func runTest(n int) string {
  var l int

  start := time.Now()
  for i:=0; i < 100*n; i++{
    for j:=0; j < 100; j++{
      for k:=0; k < 100; k++{
        l = numToTest / (i*j*k+1)
      }
    }
  }
  end := time.Now()
  fmt.Printf("Time taken for %d iterations: %s\n", n, end.Sub(start)) 
  fmt.Printf("Useless result: %d\n", l)
  return fmt.Sprintf("Time taken for %d iterations: %s\n", n, end.Sub(start))
}

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    testLength, err := strconv.Atoi(r.URL.Query().Get("n"))
    if (err == nil && testLength > 0) {
      w.WriteHeader(http.StatusOK)
      w.Write([]byte(runTest(testLength)))
    } else {
      w.WriteHeader(http.StatusBadRequest)
    }
  })
  http.ListenAndServe(":8080", nil)
}
