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

	cephrookiov1 "github.com/koor-tech/koor/pkg/apis/ceph.rook.io/v1"
	versioned "github.com/koor-tech/koor/pkg/client/clientset/versioned"
	internalinterfaces "github.com/koor-tech/koor/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/koor-tech/koor/pkg/client/listers/ceph.rook.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CephBucketTopicInformer provides access to a shared informer and lister for
// CephBucketTopics.
type CephBucketTopicInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CephBucketTopicLister
}

type cephBucketTopicInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCephBucketTopicInformer constructs a new informer for CephBucketTopic type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCephBucketTopicInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCephBucketTopicInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCephBucketTopicInformer constructs a new informer for CephBucketTopic type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCephBucketTopicInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CephV1().CephBucketTopics(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CephV1().CephBucketTopics(namespace).Watch(context.TODO(), options)
			},
		},
		&cephrookiov1.CephBucketTopic{},
		resyncPeriod,
		indexers,
	)
}

func (f *cephBucketTopicInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCephBucketTopicInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cephBucketTopicInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cephrookiov1.CephBucketTopic{}, f.defaultInformer)
}

func (f *cephBucketTopicInformer) Lister() v1.CephBucketTopicLister {
	return v1.NewCephBucketTopicLister(f.Informer().GetIndexer())
}
