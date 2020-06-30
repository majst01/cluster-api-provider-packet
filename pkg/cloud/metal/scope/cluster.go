/*
Copyright 2020 The Kubernetes Authors.

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

package scope

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	infrav1 "github.com/metal-stack/cluster-api-provider-metal/api/v1alpha3"

	"k8s.io/klog/klogr"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	Client       client.Client
	Logger       logr.Logger
	Cluster      *clusterv1.Cluster
	MetalCluster *infrav1.MetalCluster
}

// NewClusterScope creates a new ClusterScope from the supplied parameters.
// This is meant to be called for each reconcile iteration only on MetalClusterReconciler.
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("Cluster is required when creating a ClusterScope")
	}
	if params.MetalCluster == nil {
		return nil, errors.New("MetalCluster is required when creating a ClusterScope")
	}
	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	helper, err := patch.NewHelper(params.MetalCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ClusterScope{
		Logger:       params.Logger,
		client:       params.Client,
		Cluster:      params.Cluster,
		MetalCluster: params.MetalCluster,
		patchHelper:  helper,
	}, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	MetalCluster *infrav1.MetalCluster
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.patchHelper.Patch(context.TODO(), s.MetalCluster)
}

// Name returns the cluster name.
func (s *ClusterScope) Name() string {
	return s.Cluster.GetName()
}

// Namespace returns the cluster namespace.
func (s *ClusterScope) Namespace() string {
	return s.Cluster.GetNamespace()
}

// SetReady sets the MetalCluster Ready Status
func (s *ClusterScope) SetReady() {
	s.MetalCluster.Status.Ready = true
}