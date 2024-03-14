package main

import (
  "sync"

	"k8s.io/client-go/kubernetes"
)

type resourcesHandler struct {
  clientset *kubernetes.Clientset
  mu sync.Mutex
  inFlightRequests int
  isOn bool
}

// Create a new request handler
func newResourceHandler(clientset *kubernetes.Clientset) *resourcesHandler {
  return &resourcesHandler{
    clientset: clientset,
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
  if r.isOn == false {
    r.isOn = true
    turnOnContainer(r.clientset)
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
    r.isOn = false
    // Turn off the container
    turnOffContainer(r.clientset)
  }
}
