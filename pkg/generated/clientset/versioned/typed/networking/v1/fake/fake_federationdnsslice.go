/*
Copyright The Kubernetes Authors.

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

	networkingv1 "github.com/ZDragon/coredns-hot-update/pkg/apis/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFederationDNSSlices implements FederationDNSSliceInterface
type FakeFederationDNSSlices struct {
	Fake *FakeNetworkingV1
	ns   string
}

var federationdnsslicesResource = schema.GroupVersionResource{Group: "networking.synapse.sber", Version: "v1", Resource: "federationdnsslices"}

var federationdnsslicesKind = schema.GroupVersionKind{Group: "networking.synapse.sber", Version: "v1", Kind: "FederationDNSSlice"}

// Get takes name of the federationDNSSlice, and returns the corresponding federationDNSSlice object, and an error if there is any.
func (c *FakeFederationDNSSlices) Get(ctx context.Context, name string, options v1.GetOptions) (result *networkingv1.FederationDNSSlice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(federationdnsslicesResource, c.ns, name), &networkingv1.FederationDNSSlice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*networkingv1.FederationDNSSlice), err
}

// List takes label and field selectors, and returns the list of FederationDNSSlices that match those selectors.
func (c *FakeFederationDNSSlices) List(ctx context.Context, opts v1.ListOptions) (result *networkingv1.FederationDNSSliceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(federationdnsslicesResource, federationdnsslicesKind, c.ns, opts), &networkingv1.FederationDNSSliceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &networkingv1.FederationDNSSliceList{ListMeta: obj.(*networkingv1.FederationDNSSliceList).ListMeta}
	for _, item := range obj.(*networkingv1.FederationDNSSliceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested federationDNSSlices.
func (c *FakeFederationDNSSlices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(federationdnsslicesResource, c.ns, opts))

}

// Create takes the representation of a federationDNSSlice and creates it.  Returns the server's representation of the federationDNSSlice, and an error, if there is any.
func (c *FakeFederationDNSSlices) Create(ctx context.Context, federationDNSSlice *networkingv1.FederationDNSSlice, opts v1.CreateOptions) (result *networkingv1.FederationDNSSlice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(federationdnsslicesResource, c.ns, federationDNSSlice), &networkingv1.FederationDNSSlice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*networkingv1.FederationDNSSlice), err
}

// Update takes the representation of a federationDNSSlice and updates it. Returns the server's representation of the federationDNSSlice, and an error, if there is any.
func (c *FakeFederationDNSSlices) Update(ctx context.Context, federationDNSSlice *networkingv1.FederationDNSSlice, opts v1.UpdateOptions) (result *networkingv1.FederationDNSSlice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(federationdnsslicesResource, c.ns, federationDNSSlice), &networkingv1.FederationDNSSlice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*networkingv1.FederationDNSSlice), err
}

// Delete takes name of the federationDNSSlice and deletes it. Returns an error if one occurs.
func (c *FakeFederationDNSSlices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(federationdnsslicesResource, c.ns, name), &networkingv1.FederationDNSSlice{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFederationDNSSlices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(federationdnsslicesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &networkingv1.FederationDNSSliceList{})
	return err
}

// Patch applies the patch and returns the patched federationDNSSlice.
func (c *FakeFederationDNSSlices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *networkingv1.FederationDNSSlice, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(federationdnsslicesResource, c.ns, name, pt, data, subresources...), &networkingv1.FederationDNSSlice{})

	if obj == nil {
		return nil, err
	}
	return obj.(*networkingv1.FederationDNSSlice), err
}
