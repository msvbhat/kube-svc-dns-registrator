package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset    *kubernetes.Clientset
	kubeconfig   string
	dnsName      string
	hostedZoneId string
)

const (
	controllerAnnotation = "kube-svc-route53-registrator"
)

func init() {
	kubeconfig = os.Getenv("KUBECONFIG")
	dnsName = os.Getenv("DNS_NAME")
	hostedZoneId = os.Getenv("HOSTED_ZONE_ID")
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Endpoints().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: onEndpointAdd,
		UpdateFunc: func(old, new interface{}) {
			onEndpointAdd(new)
		},
	})
	informer.Run(stopper)
}

// isServiceElligible returns true if we should create Route53 record for the service
func isServiceElligible(svc *v1.Service) bool {
	annotation, ok := svc.Annotations[controllerAnnotation]
	if !ok || annotation != "true" {
		return false
	}
	return true
}

func extractIpAddresses(ep *v1.Endpoints) []string {
	ips := make([]string, 0)
	for _, subset := range ep.Subsets {
		for _, address := range subset.Addresses {
			ips = append(ips, address.IP)
		}
	}
	return ips
}

func onEndpointAdd(obj interface{}) {
	// casting obj as Endpoints
	ep := obj.(*v1.Endpoints)
	log.Printf("New Endpoint added: %s", ep.Name)

	svc, err := clientset.CoreV1().Services(ep.Namespace).Get(context.TODO(), ep.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("Service %s is not found in namespace %s", ep.Name, ep.Namespace)
			return
		} else {
			log.Printf("Error finding service %s", ep.Name)
			return
		}
	}
	if !isServiceElligible(svc) {
		log.Printf("Skipping service %s since we are not responsible for it.", svc.Name)
		return
	}
	ips := extractIpAddresses(ep)
	name := fmt.Sprintf("%s.%s", svc.Name, dnsName)
	log.Printf("The IP Addresses of %s are: %v", name, ips)
	route53CreateRecord(hostedZoneId, name, ips)
}
