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
	time "time"

	gloosoloiov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	versioned "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned"
	internalinterfaces "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/informers/externalversions/internalinterfaces"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/listers/gloo.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SettingsInformer provides access to a shared informer and lister for
// Settingses.
type SettingsInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.SettingsLister
}

type settingsInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSettingsInformer constructs a new informer for Settings type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSettingsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSettingsInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSettingsInformer constructs a new informer for Settings type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSettingsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GlooV1().Settingses(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GlooV1().Settingses(namespace).Watch(options)
			},
		},
		&gloosoloiov1.Settings{},
		resyncPeriod,
		indexers,
	)
}

func (f *settingsInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSettingsInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *settingsInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&gloosoloiov1.Settings{}, f.defaultInformer)
}

func (f *settingsInformer) Lister() v1.SettingsLister {
	return v1.NewSettingsLister(f.Informer().GetIndexer())
}
