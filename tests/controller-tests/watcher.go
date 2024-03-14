package main

import (
  "context"  
  "log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func watchDeployments(
  clientset *kubernetes.Clientset,
  selector *metav1.ListOptions,
) error {
  deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

  watch, err := deploymentsClient.Watch(context.Background(), *selector)
  if err != nil {
    return err
  }

  for event := range watch.ResultChan() {
    item := event.Object.(*appsv1.Deployment)
    log.Printf("  [W] Event: %s, deployment: %s\n", event.Type, item.Name)
  }

  return nil
}
