/*
Copyright Mehedi Hasan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/heheh13/crd-controller/custom/apis/heheh.com/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DestroymentLister helps list Destroyments.
// All objects returned here must be treated as read-only.
type DestroymentLister interface {
	// List lists all Destroyments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Destroyment, err error)
	// Destroyments returns an object that can list and get Destroyments.
	Destroyments(namespace string) DestroymentNamespaceLister
	DestroymentListerExpansion
}

// destroymentLister implements the DestroymentLister interface.
type destroymentLister struct {
	indexer cache.Indexer
}

// NewDestroymentLister returns a new DestroymentLister.
func NewDestroymentLister(indexer cache.Indexer) DestroymentLister {
	return &destroymentLister{indexer: indexer}
}

// List lists all Destroyments in the indexer.
func (s *destroymentLister) List(selector labels.Selector) (ret []*v1.Destroyment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Destroyment))
	})
	return ret, err
}

// Destroyments returns an object that can list and get Destroyments.
func (s *destroymentLister) Destroyments(namespace string) DestroymentNamespaceLister {
	return destroymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DestroymentNamespaceLister helps list and get Destroyments.
// All objects returned here must be treated as read-only.
type DestroymentNamespaceLister interface {
	// List lists all Destroyments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Destroyment, err error)
	// Get retrieves the Destroyment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Destroyment, error)
	DestroymentNamespaceListerExpansion
}

// destroymentNamespaceLister implements the DestroymentNamespaceLister
// interface.
type destroymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Destroyments in the indexer for a given namespace.
func (s destroymentNamespaceLister) List(selector labels.Selector) (ret []*v1.Destroyment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Destroyment))
	})
	return ret, err
}

// Get retrieves the Destroyment from the indexer for a given namespace and name.
func (s destroymentNamespaceLister) Get(name string) (*v1.Destroyment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("destroyment"), name)
	}
	return obj.(*v1.Destroyment), nil
}
