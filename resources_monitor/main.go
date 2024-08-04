package main

import (
	"log"
  "flag"
	"os"
	"path/filepath"
  "context"
  "time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/apimachinery/pkg/api/resource"
)

var toMonitor string
var sum int64
var count int64
var delay time.Duration = 200 * time.Millisecond

func main() {
  log.SetOutput(os.Stdout)

  toMonitor = readEnv()

  clientset := getClientset()
  for true {
    cpuRequest := getPodCpuRequest(clientset)
    sum += cpuRequest.MilliValue()
    count++
    log.Println("Average CPU request: ", sum/count)
    time.Sleep(delay)
  }

}

func getPodCpuRequest(clientset *kubernetes.Clientset) resource.Quantity {
  var selector = metav1.ListOptions{
    FieldSelector: "metadata.name=" + toMonitor,
  }
  podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)

  list, err := podsClient.List(context.Background(), selector)
  if err != nil {
    panic(err)}
  if len(list.Items) != 1 {
    log.Println("No unique pod found. Exiting...")
    panic("No unique pod found")}
  pod := list.Items[0]

  if pod.Spec.Containers[0].Resources.Requests == nil {
    return resource.Quantity{}
  }
  return pod.Spec.Containers[0].Resources.Requests["cpu"]
}

func getClientset() *kubernetes.Clientset {
  var kubeconfig *string
  // Building homedir, needed for the kubeconfig file
  if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

  var config *rest.Config
  var configErr error

  // Checking if the kubeconfig file exists, otherwise supposing the program is running inside a pod
  if _, err := os.Stat(*kubeconfig); err == nil {
    // File exists
    config, configErr = clientcmd.BuildConfigFromFlags("", *kubeconfig)
  } else {
    config, configErr = rest.InClusterConfig()
  }
  if configErr != nil {
    panic(configErr)
  }
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err) }
  return clientset
}

// This function gets the toMonitor variable
func readEnv() string {
  toMonitor := os.Getenv("POD_TO_MONITOR")
  if toMonitor == "" {
    log.Println("POD_TO_MONITOR not set. Exiting...")
    panic("POD_TO_MONITOR not set")
  }
  return toMonitor
}
