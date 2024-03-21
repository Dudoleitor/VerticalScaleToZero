package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)

var listenPort string
var forwardingToPort string

func main() {
  log.SetOutput(os.Stdout)

  listenPort, forwardingToPort = readEnv()

  remote, err := url.Parse("http://localhost:" + forwardingToPort)
  if err != nil {
    panic(err) }
  
  clientset := getClientset()
  resHandler := newResourceHandler(clientset)

  proxyReqHandler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
                  resHandler.reqIncoming()
                  defer resHandler.reqServed()

                  log.Println(r.URL)
                  r.Host = remote.Host
                  p.ServeHTTP(w, r)
          }
  }

  turnOffContainer(clientset)  // The container needs to be turned off at the beginning
  proxy := httputil.NewSingleHostReverseProxy(remote)
  http.HandleFunc("/", proxyReqHandler(proxy))

  log.Println("Listening on port: " + listenPort)

  err = http.ListenAndServe(listenPort, nil)
  if err != nil {
    panic(err) }
}

func getClientset() *kubernetes.Clientset {
  var config *rest.Config
  var configErr error
  config, configErr = rest.InClusterConfig()
  if configErr != nil {
    panic(configErr)
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err) }
  return clientset
}

// This function sets the listenPort and the forwardingToPort
func readEnv() (string, string) {
  listenPort := os.Getenv("LISTEN_PORT")
  if listenPort == "" {
    log.Println("LISTEN_PORT not set, using default 8080")
    listenPort = ":8080"
  } else {
    listenPort = ":" + listenPort
  }

  forwardingToPort := os.Getenv("FORWARDING_TO_PORT")
  if forwardingToPort == "" {
    log.Println("FORWARDING_TO_PORT not set, using default 80")
    forwardingToPort = "80"
  }
  return listenPort, forwardingToPort
}
