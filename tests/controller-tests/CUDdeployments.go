package main

import (
	"context"
	"log"
  "fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createDeployments(
  clientset *kubernetes.Clientset,
  deployments []appsv1.Deployment,
) error {
  deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

  for _, deployment := range deployments {
    result, err := deploymentsClient.Create(
      context.TODO(),
      &deployment,
      metav1.CreateOptions{},
      )

    if err != nil {
      return err
    }
    log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
  }
  return nil
}

func deleteDeployments(
  clientset *kubernetes.Clientset,
  deploymentNames []string,
) error {
  deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

  for _, deploymentName := range deploymentNames {
    deletePolicy := metav1.DeletePropagationForeground
    if err := deploymentsClient.Delete(
      context.TODO(),
      deploymentName,
      metav1.DeleteOptions{
        PropagationPolicy: &deletePolicy,
      },
      ); err != nil {
      return err
    }
    log.Printf("Deleted deployment %q.\n", deploymentName)
  }
  return nil
}

func updateDeployments(
  clientset *kubernetes.Clientset,
  selector metav1.ListOptions,
) error {
  deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
  results, getErr := deploymentsClient.List(context.TODO(), selector)
  if getErr != nil {
    panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
  }
  log.Printf("Found %d deployments with a matching label\n", len(results.Items))

  for _, result := range results.Items {
    result.Spec.Replicas = int32Ptr(10)
    _, updateErr := deploymentsClient.Update(
      context.TODO(),
      &result,
      metav1.UpdateOptions{},
    )
    if updateErr != nil {
      log.Printf("Failed to update deployment: %v\n", updateErr)
      return updateErr
    }
  }
  return nil
}

func listPods(
  clientset *kubernetes.Clientset,
  selector metav1.ListOptions,
) error {
  podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
  podList, err := podsClient.List(context.TODO(), selector)
  if err != nil {
    panic(fmt.Errorf("Failed to list pods: %v", err))
  }
  log.Printf("There are %d pods in the cluster\n", len(podList.Items))
  for _, pod := range podList.Items {
    log.Printf("Pod name %s\n", pod.Name)
  }
  return nil
}
