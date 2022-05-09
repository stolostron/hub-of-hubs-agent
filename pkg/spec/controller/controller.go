package controller

import (
	"fmt"

	clustersv1 "github.com/open-cluster-management/api/cluster/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/stolostron/hub-of-hubs-agent/pkg/helper"
	"github.com/stolostron/hub-of-hubs-agent/pkg/spec/controller/syncers"
	"github.com/stolostron/hub-of-hubs-agent/pkg/spec/controller/workers"
	consumer "github.com/stolostron/hub-of-hubs-agent/pkg/transport/consumer"
)

// AddToScheme adds only resources that have to be fetched.
// (no need to add scheme of resources that are applied as a part of generic bundles).
func AddToScheme(scheme *runtime.Scheme) error {
	if err := clustersv1.Install(scheme); err != nil {
		return fmt.Errorf("failed to install scheme: %w", err)
	}
	return nil
}

func AddControllerToManager(manager ctrl.Manager, consumer consumer.Consumer, env helper.EnvironmentManager) error {
	
	workerPool, err := workers.AddWorkerPool(ctrl.Log.WithName("workers-pool"), env.SpecWorkPoolSize, manager)
	if err != nil {
		return fmt.Errorf("failed to add worker pool to runtime manager: %w", err)
	}

	if err = syncers.AddGenericBundleSyncer(ctrl.Log.WithName("generic-bundle-syncer"), manager, 
	env.SpecEnforceHohRbac, consumer, workerPool); err != nil {
		return fmt.Errorf("failed to add bundles spec syncer to runtime manager: %w", err)
	}

	// add managed cluster labels syncer to mgr
	if err = syncers.AddManagedClusterLabelsBundleSyncer(ctrl.Log.WithName("managed-clusters-labels-syncer"), manager,
		consumer, workerPool); err != nil {
		return fmt.Errorf("failed to add managed cluster labels syncer to runtime manager: %w", err)
	}
	
	return nil
}

