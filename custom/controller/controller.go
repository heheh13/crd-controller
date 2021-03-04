package controller

import (
	"context"
	"fmt"
	"log"
	"time"

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

//Controller is the controller implementation for Custom destroyment Resource.
type Controller struct {
	//kubeclientset is the standard kubernetes clientset
	kubeclientset kubernetes.Interface
	//sampleclientset is a clientset for our own api group
	sampleclientset   clientset.Interface
	deploymentsLister appslisters.DeploymentLister

	destroymentLister lister.DestroymentLister
	destroymentSyncd  cache.InformerSynced
	//workqueue is a rate limited workqueue
	// main goal of using a workque to never work with the two or more api object  of same type at the same time
	workqueue workqueue.RateLimitingInterface
}

//NewController a constructor.
// returns a new controller for controlling the custom resource.
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
		AddFunc: controller.enqueDestroyment,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueDestroyment(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			controller.enqueDestroyment(obj)
		},
	})

	// able to use a deploymentinformer??...
	// see https://github.com/kubernetes/sample-controller/blob/master/controller.go#L129

	return controller

}
func (c *Controller) enqueDestroyment(obj interface{}) {
	log.Println("Enqueueing Foo. . . ")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// Run will set up the event handlers for the type we are interested
// will block until stopCh is closed
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {

	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	log.Println("starting Controller")

	log.Println("waiting for caches to sync")
	//fmt.Println("......................................................")
	//des, err := c.sampleclientset.HehehV1().Destroyments("default").Get(context.TODO(), "destroyment", metav1.GetOptions{})
	//if err != nil {
	//	fmt.Println("err===========" + err.Error())
	//}
	//spew.Dump(des)
	//
	//fmt.Println("..................................................................")
	if !cache.WaitForCacheSync(stopCh, c.destroymentSyncd) {
		return fmt.Errorf("failed to sync")
	}
	log.Println("starting workers!!")

	//launch <threadiness> worker to process Custom resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

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

// will implement the bussness logic here..
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
		log.Printf("\n\n deployment created %s\n\n", deployment.Name)
	}

	if destroyment.Spec.Replicas != nil && *destroyment.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("destroyment %s replicas: %d, deployment replicas: %d\n\n", name, *destroyment.Spec.Replicas, *deployment.Spec.Replicas)
		if *destroyment.Spec.Replicas > *deployment.Spec.Replicas {
			log.Println("scalling up deploymnet to ", *destroyment.Spec.Replicas, ".......")
		} else {
			log.Println("scaling down deployment to ", *destroyment.Spec.Replicas, "............")
		}

		deployment, err = c.kubeclientset.AppsV1().Deployments(destroyment.Namespace).Update(context.TODO(), newDeployment(*destroyment), metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	err = c.updateDestroymnetStatus(destroyment, deployment)
	if err != nil {
		return err
	}

	serviceName := destroyment.Name
	if serviceName == "" {
		utilruntime.HandleError(fmt.Errorf("must provide service name"))
		return nil
	}
	//if service avail able check for it
	//else create a new one
	service, err := c.kubeclientset.CoreV1().Services(destroyment.Namespace).Get(context.TODO(), destroyment.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(destroyment.Namespace).Create(context.TODO(), newService(*destroyment), metav1.CreateOptions{})
		if err != nil {
			log.Println("error while creating serviecs...")
			return err
		}
		log.Printf("\n\nservice %s created ....\n\n", service.Name)
	} else if err != nil {
		log.Println("error while getting serviecs...")
		return err
	}

	_, err = c.kubeclientset.CoreV1().Services(destroyment.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	//fmt.Println(service)
	return nil
}

func (c *Controller) updateDestroymnetStatus(destroyment *v1.Destroyment, deployment *appsv1.Deployment) error {
	fmt.Println("updating status!!!!!!!!!!!")
	destroymentcopy := destroyment.DeepCopy()
	destroymentcopy.Status.Phase = "running..."
	destroymentcopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	destroymentcopy.Status.Replicas = deployment.Status.Replicas
	//_, err := c.sampleclientset.CrdcntrlrV1alpha1().Foos(foo.Namespace).Update(fooCopy)
	_, err := c.sampleclientset.HehehV1().Destroyments(destroymentcopy.Namespace).UpdateStatus(context.TODO(), destroymentcopy, metav1.UpdateOptions{})
	return err

}

func newDeployment(destroyment v1.Destroyment) *appsv1.Deployment {
	labels := map[string]string{
		"k8s.io/app": destroyment.Name,
		"controller": destroyment.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      destroyment.Name,
			Namespace: destroyment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&destroyment, v1.SchemeGroupVersion.WithKind("Destroyment")),
			},
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

func newService(destroyment v1.Destroyment) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: destroyment.Name,
			Labels: map[string]string{
				"k8s.io/app":  destroyment.Name,
				"k8s.io/name": destroyment.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&destroyment, v1.SchemeGroupVersion.WithKind("Destroyment")),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:        "",
					Protocol:    corev1.ProtocolTCP,
					AppProtocol: nil,
					Port:        destroyment.Spec.Container.Port,
				},
			},
			Type: corev1.ServiceType(destroyment.Spec.ServiceSpec.ServiceType),
			Selector: map[string]string{
				"k8s.io/app": destroyment.Name,
			},
		},
	}
}
