package main

import (
  "sync"
  "time"
	"k8s.io/client-go/kubernetes"
  "log"
)

type resourcesHandler struct {
  clientset *kubernetes.Clientset
  delayS int
  mu sync.Mutex
  ch chan bool
  inFlightRequests int
  isOn bool
}

// Create a new request handler
func newResourceHandler(clientset *kubernetes.Clientset, delaySec int) *resourcesHandler {
  return &resourcesHandler{
    clientset: clientset,
    delayS: delaySec,
    ch: make(chan bool, 1),
    inFlightRequests: 0,
    isOn: false,
  }
}

// This function will be called when a request is received
func (r *resourcesHandler) reqIncoming() {
  r.mu.Lock()
  defer r.mu.Unlock()

  r.inFlightRequests++
  // Turn on the container
  if r.inFlightRequests == 1 {  // The first request
    if r.isOn == false {  // Delay already expired
      r.isOn = true
      turnOnContainer(r.clientset)
    } else {
      if r.delayS > 0 {
        select {
          case r.ch <- false:  // Signal to stop the delay, no need to turn off the container
          default: log.Println("WARN: Signal channel is full")
        }
      }
    }
  }
}

// This function will be called when a request has been served
func (r *resourcesHandler) reqServed() {
  r.mu.Lock()
  defer r.mu.Unlock()

  if r.inFlightRequests == 0 {
    panic("Fatal error, inFlightRequests is zero") }

  r.inFlightRequests--

  if r.inFlightRequests == 0 { 
    if r.delayS > 0 {
      // Start time to turn off the turn off the container
      r.turnOffDelay(time.Duration(r.delayS) * time.Second)
    } else {
      r.isOn = false
      turnOffContainer(r.clientset)
    }
  }
}

func (r *resourcesHandler) turnOffDelay(delay time.Duration) {
  interrupted := false
  localMu := sync.Mutex{}

  // Emptying the channel
  for len(r.ch) > 0 { 
    <- r.ch
  }

  go func() {  // Goroutine to wait for the delay
    time.Sleep(delay)
  
    localMu.Lock()
    defer localMu.Unlock()
    if !interrupted {
      select{
        case r.ch <- true:
        default: log.Println("WARN: Signal channel is full")
      }
    }
  } ()

  go func() { // Goroutine to wait for the signal and turn off the container
    todo := <- r.ch // Waiting for a signal
    
    if todo == true {  // Need to turn off the container
      r.mu.Lock()
      defer r.mu.Unlock()
      r.isOn = false
      turnOffContainer(r.clientset)
    } else {
      localMu.Lock()
      defer localMu.Unlock()
      interrupted = true
    }
  }()
}
