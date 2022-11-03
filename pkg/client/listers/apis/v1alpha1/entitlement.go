/*
Copyright The KCP Authors.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
)

// EntitlementLister helps list Entitlements.
// All objects returned here must be treated as read-only.
type EntitlementLister interface {
	// List lists all Entitlements in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Entitlement, err error)
	// Get retrieves the Entitlement from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Entitlement, error)
	EntitlementListerExpansion
}

// entitlementLister implements the EntitlementLister interface.
type entitlementLister struct {
	indexer cache.Indexer
}

// NewEntitlementLister returns a new EntitlementLister.
func NewEntitlementLister(indexer cache.Indexer) EntitlementLister {
	return &entitlementLister{indexer: indexer}
}

// List lists all Entitlements in the indexer.
func (s *entitlementLister) List(selector labels.Selector) (ret []*v1alpha1.Entitlement, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Entitlement))
	})
	return ret, err
}

// Get retrieves the Entitlement from the index for a given name.
func (s *entitlementLister) Get(name string) (*v1alpha1.Entitlement, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("entitlement"), name)
	}
	return obj.(*v1alpha1.Entitlement), nil
}