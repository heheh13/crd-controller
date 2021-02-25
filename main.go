package main

import (
	"context"
	"fmt"
	"os/signal"
	"time"

	v1 "github.com/heheh13/crd/custom/apis/heheh.com/v1"
	"github.com/kr/pretty"

	destroymentclientset "github.com/heheh13/crd/custom/client/clientset/versioned"
	crdapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube/config")
	fmt.Print(kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal("couldn't get kubernetes config %s", err)
	}
	log.Println("custom resource Destroymet Craeting.....")

	crdClient, err := crdclientset.NewForConfig(config)
	if err != nil {
		log.Fatal("error while configuring crd client")
	}

	mycrd := &crdapi.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "destroyments.heheh.com",
		},
		Spec: crdapi.CustomResourceDefinitionSpec{
			Group: "heheh.com",
			Names: crdapi.CustomResourceDefinitionNames{
				Plural:   "destroyments",
				Singular: "destroyment",
				ShortNames: []string{
					"ds",
				},
				Kind: "Destroyment",
			},
			Scope: "Namespaced",
			Versions: []crdapi.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &crdapi.CustomResourceValidation{
						OpenAPIV3Schema: &crdapi.JSONSchemaProps{

							Type: "object",

							Properties: map[string]crdapi.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]crdapi.JSONSchemaProps{
										"replicas": {
											Type: "integer",
										},
										"container": {
											Type: "object",
											Properties: map[string]crdapi.JSONSchemaProps{
												"image": {
													Type: "string",
												},
												"port": {
													Type: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_ = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), mycrd.Name, metav1.DeleteOptions{})
	_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), mycrd, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	log.Println("custom resoucse 'destroyment' Created!")

	//deleting destroyment?
	defer func() {
		log.Println("deleteing destroyment!")
		if err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), mycrd.Name, metav1.DeleteOptions{}); err != nil {
			panic(err)
		}
		log.Println("destroyment deleted")
	}()

	log.Println("Press Ctrl+C to Create an instance of destroyment. . . .")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch ///taking input?

	log.Println("creating destroyment ... ")

	destroymentClient, err := destroymentclientset.NewForConfig(config)
	destroyment := &v1.Destroyment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "destroyment",
		},
		Spec: v1.DestroymentSpec{
			Replicas: 2,
			Container: v1.ContainerSpec{
				Image: "nginx",
				Port:  8080,
			},
		},
		Status: v1.DestroymentStatus{
			Phase: "....ing",
		},
	}
	ds, err := destroymentClient.HehehV1().Destroyments("default").Create(context.TODO(), destroyment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Println("destroyment created !!!!! ")
	fmt.Println(ds)
	pretty.Print(ds)

	defer func() {
		log.Println("deleting destroyment .......")
		if err := destroymentClient.HehehV1().Destroyments("default").Delete(
			context.TODO(), destroyment.Name, metav1.DeleteOptions{}); err != nil {
			panic(err)
		}
		log.Println("destroyment deleted")
	}()
	log.Println("Press Ctrl+C to Delete the Custom Resource 'destroyment' and the instance of d3stroyment'. . . .")

	ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}
