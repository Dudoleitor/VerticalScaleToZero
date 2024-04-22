package main

import(
  "fmt"
  "net/http"
  "time"
)

var target1 = "http://exampleworkload-not-proxied:80/?n="
var target2 = "http://exampleworkload-proxied:8080/?n="


func main() {
  run_test(target1)
  run_test(target2)
}

func run_test(target string) {
  fmt.Println("Testing different values of n for target: ", target)
  fmt.Println("n time")
  for i := 1; i <= 650; i++ {
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
