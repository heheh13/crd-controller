package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	clientset "github.com/heheh13/crd-controller/custom/client/clientset/versioned"
	customInformers "github.com/heheh13/crd-controller/custom/client/informers/externalversions"
	cntrl "github.com/heheh13/crd-controller/custom/controller"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)
	customClient := clientset.NewForConfigOrDie(config)

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	customInformerFactory := customInformers.NewSharedInformerFactory(customClient, time.Second*30)

	controller := cntrl.NewController(
		kubeClient,
		customClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		customInformerFactory.Heheh().V1().Destroyments(),
	)

	stopch := make(chan struct{})

	kubeInformerFactory.Start(stopch)
	customInformerFactory.Start(stopch)

	if err = controller.Run(2, stopch); err != nil {
		fmt.Println("error running controller ", err.Error())
	}

}
