package main

import (
  "context"  
  "os"
  "log"
  "fmt"

	appsv1 "k8s.io/api/apps/v1"
  typedv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

const containerName = "function-proxy"
const containerImage = "http_red:v0.0.1"

func watchDeployments(
  clientset *kubernetes.Clientset,
  selector *metav1.ListOptions,
) error {
  log.SetOutput(os.Stdout)
  deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

  watch, err := deploymentsClient.Watch(context.Background(), *selector)
  if err != nil {
    return err }

  // Iterating over the events
  for event := range watch.ResultChan() {
    item := event.Object.(*appsv1.Deployment)
    log.Printf("  [W] Event: %s, deployment: %s\n", event.Type, item.Name)

    switch event.Type {
    case "ADDED", "MODIFIED":
      handlePodInjection(deploymentsClient, *item, 8080, 80)
    }
  }
  return nil
}

func handlePodInjection(
  deploymentsClient typedv1.DeploymentInterface,
  deployment appsv1.Deployment,
  containerPort int32,
  portToForward int32,
) error {
  // Defining the proxy container that will be injected
  proxyContainer := apiv1.Container{
    Name: containerName,
    Image: containerImage,
    Ports: []apiv1.ContainerPort{
      {
        Name:          "http",
        Protocol:      apiv1.ProtocolTCP,
        ContainerPort: containerPort,
      },
    },
    Env: []apiv1.EnvVar{
      {
        Name:  "LISTEN_PORT",
        Value: fmt.Sprintf("%d", containerPort),
      },
      {
        Name:  "FORWARDING_TO_PORT",
        Value: fmt.Sprintf("%d", portToForward),
      },
    },
  }

  retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    containers := &deployment.Spec.Template.Spec.Containers

    // Looking for the proxy container
    existingContainer := checkIfProxyContainerExists(*containers)

    if existingContainer != nil {
      log.Println("Proxy container already exists in", deployment.Name)
      if checkIfProxyContainerIsAsDesired(
        *existingContainer,
        containerPort,
        portToForward) {
        log.Printf("Proxy container is as desired in %s, skipping\n", deployment.Name)
        return nil
      }

      log.Printf("Proxy container is not as desired in %s, deleting it\n", deployment.Name)
      delErr := deleteContainer(deploymentsClient, deployment, existingContainer)
      if delErr != nil {
        log.Println("Failed to delete proxy container:", delErr)
        return delErr
      }
    }

    // Injecting the proxy container
    *containers = append(
      *containers,
      proxyContainer)
    log.Println("Injecting into", deployment.Name)

    _, updateErr := deploymentsClient.Update(
      context.TODO(),
      &deployment,
      metav1.UpdateOptions{},
    )
    if updateErr != nil {
      log.Println("Failed to update deployment:", updateErr)
      return updateErr
    }
    return nil
	})
  return retryErr
}

// Check if the proxy container already exists,
// if so, return the container else return nil
func checkIfProxyContainerExists(
  containers []apiv1.Container,
) *apiv1.Container {
  for _, container := range containers {
    if container.Name == containerName {
      return &container
    }
  }
  return nil
}

// Check if the proxy container has already been 
// injected matches the desired configuration
func checkIfProxyContainerIsAsDesired(
  container apiv1.Container,
  containerPort int32,
  portToForward int32,
) bool {
  if container.Name != containerName {
    return false }

  if container.Image != containerImage { 
    return false }

  if len(container.Ports) != 1 {
    return false }

  if container.Ports[0].Name != "http" && 
    container.Ports[0].Protocol != apiv1.ProtocolTCP &&
    container.Ports[0].ContainerPort != containerPort {
    return false }

  if len(container.Env) != 2 {
    return false }

  if container.Env[0].Name != "LISTEN_PORT" &&
    container.Env[0].Value != fmt.Sprintf("%d", containerPort) &&
    container.Env[1].Name != "FORWARDING_TO_PORT" &&
    container.Env[1].Value != fmt.Sprintf("%d", portToForward) {
    return false }

  return true
}

// Delete the proxy container from the deployment
func deleteContainer(
  deploymentsClient typedv1.DeploymentInterface,
  deployment appsv1.Deployment,
  container *apiv1.Container,
) error {
  retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    containers := &deployment.Spec.Template.Spec.Containers
    newContainers := []apiv1.Container{}

    // Iterating over the containers and skipping the given one
    for _, container_checking := range *containers {
      if &container_checking != container {
        newContainers = append(newContainers, container_checking)
      }
    }

    // Overwriting the containers with the new list
    deployment.Spec.Template.Spec.Containers = newContainers

    _, updateErr := deploymentsClient.Update(
      context.TODO(),
      &deployment,
      metav1.UpdateOptions{},
    )
    if updateErr != nil {
      log.Println("Failed to update deployment:", updateErr)
      return updateErr
    }
    return nil
  })
  return retryErr
}
