package main

import (
  "bufio"
  // time "time"
	"flag"
	"log"
  "fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

  log.SetOutput(os.Stdout)

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
		panic(err)
	}

  // Defining two deployments
  deployments := [2]appsv1.Deployment{}
	deployments[0] = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
      Labels: map[string]string{
        "app": "demo",
        "to_proxy": "true",
      },
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
          "to_proxy": "true",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
            "to_proxy": "true",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:latest",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	deployments[1] = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
      Labels: map[string]string{
        "app": "test",
        "to_proxy": "true",
      },
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test",
          "to_proxy": "true",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test",
            "to_proxy": "true",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "jitesoft/lighttpd:latest",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
              Resources: apiv1.ResourceRequirements{
                Requests: apiv1.ResourceList{
                  apiv1.ResourceCPU: resource.MustParse("50m"),
                  apiv1.ResourceMemory: resource.MustParse("50Mi"),
                },
                Limits: apiv1.ResourceList{
                  apiv1.ResourceCPU: resource.MustParse("50m"),
                  apiv1.ResourceMemory: resource.MustParse("50Mi"),
                },
              },
              ResizePolicy: []apiv1.ContainerResizePolicy{
                {
                  ResourceName: apiv1.ResourceCPU,
                  RestartPolicy: apiv1.NotRequired,
                },
                {
                  ResourceName: apiv1.ResourceMemory,
                  RestartPolicy: apiv1.RestartContainer,
                },
              },
						},
					},
				},
			},
		},
	}

  // Defining the selector to retrieve the deploymentsClient
  selector := metav1.ListOptions{
    LabelSelector: "to_proxy=true",
  }

  // Listing pods
  listPods(clientset, selector)

  // Create Deployments
	log.Println("Creating deployments...")
  createErr := createDeployments(clientset, deployments[:])

	if createErr != nil {
		panic(fmt.Errorf("Update failed: %v", createErr)) }

	// Update Deployment
	prompt()
	log.Println("Updating deployment...")
  updateErr := updateDeployments(clientset, selector)

	if updateErr != nil {
		panic(fmt.Errorf("Update failed: %v", updateErr)) }
	log.Println("Updated deployments...")

	// Delete Deployment
	prompt()
	log.Println("Deleting deployments...")
  deleteErr := deleteDeployments(clientset, deploymentToString(deployments[:]))

	if deleteErr != nil {
		panic(fmt.Errorf("Update failed: %v", deleteErr)) }
  
  log.Println("Done!")
}

func prompt() {
	log.Printf("-> Press Return key to continue.\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	log.Println()
}
// func prompt() {
//   sleepTime := 5
//   log.Printf("-> Waiting %d seconds...\n", sleepTime)
//   time.Sleep(time.Duration(sleepTime) * time.Second)
// }

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
func deploymentToString(deployments []appsv1.Deployment) []string {
  var deploymentNames []string
  for _, deployment := range deployments {
    deploymentNames = append(deploymentNames, deployment.Name)
  }
  return deploymentNames
}
