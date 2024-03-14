package main

import (
  "fmt"
  "net/http"
  "sync"
)

var target = "http://localhost:8080/"
var waitGroup sync.WaitGroup
var subroutinesSemaphore sync.WaitGroup

func main() {
  subroutineWork := func() {
        defer waitGroup.Done()
        subroutinesSemaphore.Wait()
        http.Get(target)
      }
  subroutinesSemaphore.Add(1)

  fmt.Println("Spawning subroutines")
  for i := 0; i < 5000; i++ {
    waitGroup.Add(1)
    go subroutineWork()
  }

  fmt.Println("Starting subroutines")
  subroutinesSemaphore.Done()
  waitGroup.Wait()
}

