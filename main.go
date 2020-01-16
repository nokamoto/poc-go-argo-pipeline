package main

import (
	"fmt"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	println(cfg.String())

	namespace := "default"

	helloWorldWorkflow := wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "hello-world-",
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: "whalesay",
			Templates: []wfv1.Template{
				{
					Name: "whalesay",
					Container: &corev1.Container{
						Image:   "docker/whalesay:latest",
						Command: []string{"cowsay", "hello world"},
					},
				},
			},
		},
	}

	wfClient := wfclientset.NewForConfigOrDie(cfg).ArgoprojV1alpha1().Workflows(namespace)
	createdWf, err := wfClient.Create(&helloWorldWorkflow)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Workflow %s submitted\n", createdWf.Name)
}
