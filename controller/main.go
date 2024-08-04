package main

import (
	"flag"
	"os"
	"path/filepath"
  "log"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Defining the selector to retrieve the deploymentsClient
var selector = metav1.ListOptions{
  LabelSelector: "to_proxy=true",
}

func main() {
  log.SetOutput(os.Stdout)

  clientset := getClientset()

  // Starting the watcher
  log.Println("Starting watcher...")
  watchDeployments(clientset, &selector)
}

// Helper functions
func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
func deploymentToString(deployments []appsv1.Deployment) []string {
  var deploymentNames []string
  for _, deployment := range deployments {
    deploymentNames = append(deploymentNames, deployment.Name)
  }
  return deploymentNames
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
