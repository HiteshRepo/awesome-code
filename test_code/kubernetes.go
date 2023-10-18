package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const DefaultKubeConfigPath = "/home/hitesh/.kube/config"

func GetAllServices() {
	kubeconfig := flag.String("kubeconfig", "/home/hitesh/.kube/config", "location of your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("error while building config from kubeconfig file location : %s\n", err.Error())
		log.Println("fetching config within cluster")
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error while getting inclusterconfig : %s\n", err.Error())
			return
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	namespace := "default"

	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services.Items {
		if service.Name == "greetserver" {
			fmt.Println(service)
		}
	}

}
