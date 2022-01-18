package checkpoint

import (
	"context"
	"httpInterceptor/config"
	_ "httpInterceptor/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func GetContainerPid() string {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}
	context := context.Background()
	podName := config.GetPodName()
	pod, err := clientset.CoreV1().Pods("").Get(context, podName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	print(pod)
	return "Oi"
}
