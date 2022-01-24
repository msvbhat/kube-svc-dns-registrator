package main

import (
	"context"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Endpoints().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: onEndpointAdd,
	})
	informer.Run(stopper)
}

func onEndpointAdd(obj interface{}) {
	// casting obj as Endpoints
	ep := obj.(*v1.Endpoints)
	log.Printf("New Endpoint added: %s", ep.Name)

	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	svc, err := clientset.CoreV1().Services(ep.Namespace).Get(context.TODO(), ep.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Unable to find the service: %s", ep.Name)
	} else {
		log.Printf("The service found: %s and annotations are: %s", svc.Name, svc.Annotations)
	}

}
