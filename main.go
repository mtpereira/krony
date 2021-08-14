package main

import (
	"context"
	"fmt"
	"log"
	
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func connectToK8s() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("InClusterConfig: %v", err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("NewForConfig: %v", err.Error())
	}
	return clientset
}

func main() {
	clientset := connectToK8s()
	
	jobs, err := clientset.BatchV1().Jobs("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicf("main: %v", err.Error())
	}
	fmt.Printf("There are %d jobs in the cluster\n", len(jobs.Items))
}
