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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	hehehcomv1 "github.com/heheh13/crd-controller/custom/apis/heheh.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDestroyments implements DestroymentInterface
type FakeDestroyments struct {
	Fake *FakeHehehV1
	ns   string
}

var destroymentsResource = schema.GroupVersionResource{Group: "heheh.com", Version: "v1", Resource: "destroyments"}

var destroymentsKind = schema.GroupVersionKind{Group: "heheh.com", Version: "v1", Kind: "Destroyment"}

// Get takes name of the destroyment, and returns the corresponding destroyment object, and an error if there is any.
func (c *FakeDestroyments) Get(ctx context.Context, name string, options v1.GetOptions) (result *hehehcomv1.Destroyment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(destroymentsResource, c.ns, name), &hehehcomv1.Destroyment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*hehehcomv1.Destroyment), err
}

// List takes label and field selectors, and returns the list of Destroyments that match those selectors.
func (c *FakeDestroyments) List(ctx context.Context, opts v1.ListOptions) (result *hehehcomv1.DestroymentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(destroymentsResource, destroymentsKind, c.ns, opts), &hehehcomv1.DestroymentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &hehehcomv1.DestroymentList{ListMeta: obj.(*hehehcomv1.DestroymentList).ListMeta}
	for _, item := range obj.(*hehehcomv1.DestroymentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested destroyments.
func (c *FakeDestroyments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(destroymentsResource, c.ns, opts))

}

// Create takes the representation of a destroyment and creates it.  Returns the server's representation of the destroyment, and an error, if there is any.
func (c *FakeDestroyments) Create(ctx context.Context, destroyment *hehehcomv1.Destroyment, opts v1.CreateOptions) (result *hehehcomv1.Destroyment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(destroymentsResource, c.ns, destroyment), &hehehcomv1.Destroyment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*hehehcomv1.Destroyment), err
}

// Update takes the representation of a destroyment and updates it. Returns the server's representation of the destroyment, and an error, if there is any.
func (c *FakeDestroyments) Update(ctx context.Context, destroyment *hehehcomv1.Destroyment, opts v1.UpdateOptions) (result *hehehcomv1.Destroyment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(destroymentsResource, c.ns, destroyment), &hehehcomv1.Destroyment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*hehehcomv1.Destroyment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDestroyments) UpdateStatus(ctx context.Context, destroyment *hehehcomv1.Destroyment, opts v1.UpdateOptions) (*hehehcomv1.Destroyment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(destroymentsResource, "status", c.ns, destroyment), &hehehcomv1.Destroyment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*hehehcomv1.Destroyment), err
}

// Delete takes name of the destroyment and deletes it. Returns an error if one occurs.
func (c *FakeDestroyments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(destroymentsResource, c.ns, name), &hehehcomv1.Destroyment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDestroyments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(destroymentsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &hehehcomv1.DestroymentList{})
	return err
}

// Patch applies the patch and returns the patched destroyment.
func (c *FakeDestroyments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *hehehcomv1.Destroyment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(destroymentsResource, c.ns, name, pt, data, subresources...), &hehehcomv1.Destroyment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*hehehcomv1.Destroyment), err
}
