package checkpoint

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func getContainerPid() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	context := context.Background()
	pod, err := clientset.CoreV1().Pods("").Get(context, "")
	if err != nil {
		log.Fatal(err)
	}
	print(pod)
}
