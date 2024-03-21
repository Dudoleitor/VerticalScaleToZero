package main

import (
  "os"
  "context"
  "log"
  "time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/api/resource"
  "k8s.io/client-go/util/retry"
	"k8s.io/apimachinery/pkg/util/wait"
)

var defaultCpuLimit = resource.MustParse("50m")
var defaultCpuRequest = resource.MustParse("50m")

var myName = os.Getenv("HOSTNAME")

var previousCpuLimit, previousCpuRequest resource.Quantity

// Defining the selector to retrieve the pod
var selector = metav1.ListOptions{
  LabelSelector: "to_proxy=true",
}

var retryPolicy = wait.Backoff{
  Duration: 500 * time.Millisecond,
  Steps:    10,
  Factor:   1.0,
  Jitter:   0.1,
}

// Turning off the targer container
// by setting cpu resources to 1m
func turnOffContainer(
  clientset *kubernetes.Clientset,
) error {
  // Getting the current allocated cpu
  previousCpuLimit, previousCpuRequest = getCpuResource(clientset)

  log.Println("Turning off container")
  return setCpuResource(clientset, resource.MustParse("1m"), resource.MustParse("1m"))
}

// Turning on the targer container
// by setting cpu resources to 50m
func turnOnContainer(
  clientset *kubernetes.Clientset,
) error {
  log.Println("Turning on container")
  return setCpuResource(clientset, previousCpuLimit, previousCpuRequest)
}


func setCpuResource(
  clientset *kubernetes.Clientset,
  limit resource.Quantity,
  request resource.Quantity,
) error {

  retryErr := retry.RetryOnConflict(retryPolicy, func() error {
    pod := getMyPod(clientset)

    (*pod).Spec.Containers[0].Resources.Limits["cpu"] = limit
    (*pod).Spec.Containers[0].Resources.Requests["cpu"] = request
    _, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).Update(context.TODO(), pod, metav1.UpdateOptions{})
    return err
  })
  return retryErr 
}

func getCpuResource(
  clientset *kubernetes.Clientset,
) (resource.Quantity, resource.Quantity) {

  pod := getMyPod(clientset)

  if pod == nil {
    panic("Pod not found while getting cpu resources")
  }

  var cpuLimit, cpuRequest resource.Quantity
  if (*pod).Spec.Containers[0].Resources.Limits == nil {
    cpuLimit = defaultCpuLimit
  } else {
    cpuLimit = (*pod).Spec.Containers[0].Resources.Limits["cpu"]
  }

  if (*pod).Spec.Containers[0].Resources.Requests == nil {
    cpuRequest = defaultCpuRequest
  } else {
    cpuRequest = (*pod).Spec.Containers[0].Resources.Requests["cpu"]
  }

  return cpuLimit, cpuRequest
}


func getMyPod(
  clientset *kubernetes.Clientset,
) *apiv1.Pod {

  //Retrieving the pod
  podsList, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(context.Background(), selector)
  if err != nil {
    panic(err) }

  for _, pod := range podsList.Items {
    if pod.Name == myName {
      return &pod }
  }

  return nil
}
