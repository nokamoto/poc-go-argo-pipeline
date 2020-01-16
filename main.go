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

	value := "{{steps.generate-parameter.outputs.parameters.hello-param}}"

	helloWorldWorkflow := wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "output-parameter-",
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: "output-parameter",
			Templates: []wfv1.Template{
				{
					Name: "output-parameter",
					Steps: [][]wfv1.WorkflowStep{
						{
							{
								Name:     "generate-parameter",
								Template: "whalesay",
							},
						},
						{
							{
								Name:     "consume-parameter",
								Template: "print-message",
								Arguments: wfv1.Arguments{
									Parameters: []wfv1.Parameter{
										{
											Name:  "message",
											Value: &value,
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "whalesay",
					Container: &corev1.Container{
						Image:   "docker/whalesay:latest",
						Command: []string{"sh", "-c"},
						Args:    []string{"echo -n hello world > /tmp/output.txt"},
					},
					Outputs: wfv1.Outputs{
						Parameters: []wfv1.Parameter{
							{
								Name: "hello-param",
								ValueFrom: &wfv1.ValueFrom{
									Path: "/tmp/output.txt",
								},
							},
						},
					},
				},
				{
					Name: "print-message",
					Inputs: wfv1.Inputs{
						Parameters: []wfv1.Parameter{
							{
								Name: "message",
							},
						},
					},
					Container: &corev1.Container{
						Image:   "docker/whalesay:latest",
						Command: []string{"cowsay"},
						Args:    []string{"{{inputs.parameters.message}}"},
					},
				},
			},
			//PodGC: &wfv1.PodGC{
			//	Strategy: wfv1.PodGCOnWorkflowSuccess,
			//},
		},
	}

	wfClient := wfclientset.NewForConfigOrDie(cfg).ArgoprojV1alpha1().Workflows(namespace)
	createdWf, err := wfClient.Create(&helloWorldWorkflow)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Workflow %s submitted\n", createdWf.Name)

	select {} // sleep forever
}
