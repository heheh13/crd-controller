package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"

	v1 "github.com/heheh13/crd-controller/custom/apis/heheh.com/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/util/wait"

	clientset "github.com/heheh13/crd-controller/custom/client/clientset/versioned"
	informer "github.com/heheh13/crd-controller/custom/client/informers/externalversions/heheh.com/v1"
	lister "github.com/heheh13/crd-controller/custom/client/listers/heheh.com/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeclientset   kubernetes.Interface
	sampleclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister

	destroymentLister lister.DestroymentLister
	destroymentSyncd  cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
}

func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	destroymentInformer informer.DestroymentInformer) *Controller {

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		destroymentLister: destroymentInformer.Lister(),
		destroymentSyncd:  destroymentInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Destroyments"),
	}

	log.Println("Setting up event handlers")

	destroymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueueFoo(newObj)
		},
		DeleteFunc: func(obj interface{}) {},
	})

	return controller

}
func (c *Controller) enqueueFoo(obj interface{}) {
	log.Println("Enqueueing Foo. . . ")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) Run(stopCh <-chan struct{}) error {

	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	log.Println("starting Controller")

	log.Println("waiting for caches to sync")
	fmt.Println("......................................................")
	des, err := c.sampleclientset.HehehV1().Destroyments("default").Get(context.TODO(), "destroyment", metav1.GetOptions{})
	if err != nil {
		fmt.Println("err===========" + err.Error())
	}
	spew.Dump(des)
	fmt.Println("..................................................................")
	if !cache.WaitForCacheSync(stopCh, c.destroymentSyncd) {
		return fmt.Errorf("failed to sync")
	}
	log.Println("starting workers!!")

	go wait.Until(c.runWorker, time.Second, stopCh)

	log.Println("started workers")

	<-stopCh
	log.Println("shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}

func (c *Controller) processNextWorkItem() bool {

	obj, quit := c.workqueue.Get()

	if quit {
		return false
	}
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expectec string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("erro syncing %s : %s  requing ", key, err)
		}
		c.workqueue.Forget(obj)
		log.Println("succesfully synced")
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true

}

func (c *Controller) syncHandler(key string) interface{} {

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key %s", key))
		return nil
	}
	destroyment, err := c.destroymentLister.Destroyments(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("foo %s in workque not exists", key))
			return nil
		}
		return err
	}
	deploymentName := destroyment.Name

	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("%s deployment name must be specified", key))
		return nil
	}

	deployment, err := c.deploymentsLister.Deployments(destroyment.Namespace).Get(deploymentName)

	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(destroyment.Namespace).Create(context.TODO(), newDeployment(*destroyment), metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	if destroyment.Spec.Replicas != nil && *destroyment.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("destroyment %s replicas: %d, deployment replicas: %d\n\n", name, *destroyment.Spec.Replicas, *deployment.Spec.Replicas)
		if *destroyment.Spec.Replicas > *destroyment.Spec.Replicas {
			log.Println("scalling up deploymnet to ", *destroyment.Spec.Replicas, ".......")
		} else {
			log.Println("scaling down deployment to ", *destroyment.Spec.Replicas, "............")
		}

		deployment, err = c.kubeclientset.AppsV1().Deployments(destroyment.Namespace).Update(context.TODO(), newDeployment(*destroyment), metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	err = c.updateFooStatus(destroyment, deployment)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) updateFooStatus(destroyment *v1.Destroyment, deployment *appsv1.Deployment) error {
	fmt.Println("update yet to implement")
	return nil
}

func newDeployment(destroyment v1.Destroyment) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": destroyment.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      destroyment.Name,
			Namespace: destroyment.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: destroyment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  destroyment.Name,
							Image: destroyment.Spec.Container.Image,
						},
					},
				},
			},
		},
	}
}
