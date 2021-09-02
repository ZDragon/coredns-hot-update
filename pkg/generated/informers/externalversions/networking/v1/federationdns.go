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
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	networkingv1 "github.com/ZDragon/coredns-hot-update/pkg/apis/networking/v1"
	versioned "github.com/ZDragon/coredns-hot-update/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/ZDragon/coredns-hot-update/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/ZDragon/coredns-hot-update/pkg/generated/listers/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// FederationDNSInformer provides access to a shared informer and lister for
// FederationDNSs.
type FederationDNSInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.FederationDNSLister
}

type federationDNSInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFederationDNSInformer constructs a new informer for FederationDNS type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFederationDNSInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFederationDNSInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFederationDNSInformer constructs a new informer for FederationDNS type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFederationDNSInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1().FederationDNSs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1().FederationDNSs(namespace).Watch(context.TODO(), options)
			},
		},
		&networkingv1.FederationDNS{},
		resyncPeriod,
		indexers,
	)
}

func (f *federationDNSInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFederationDNSInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *federationDNSInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&networkingv1.FederationDNS{}, f.defaultInformer)
}

func (f *federationDNSInformer) Lister() v1.FederationDNSLister {
	return v1.NewFederationDNSLister(f.Informer().GetIndexer())
}