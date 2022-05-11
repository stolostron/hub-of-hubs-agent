package managedclusters

import (
	"fmt"

	"github.com/stolostron/hub-of-hubs-agent/pkg/helper"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/bundle"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/generic"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/syncintervals"
	producer "github.com/stolostron/hub-of-hubs-agent/pkg/transport/producer"
	datatypes "github.com/stolostron/hub-of-hubs-data-types"
	configV1 "github.com/stolostron/hub-of-hubs-data-types/apis/config/v1"
	clusterV1 "open-cluster-management.io/api/cluster/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	clusterStatusSyncLogName          = "clusters-status-sync"
	managedClusterManagedByAnnotation = "hub-of-hubs.open-cluster-management.io/managed-by"
)

// mgr, pro, env.LeafHubID, incarnation, config, syncIntervals
// AddClustersStatusController adds managed clusters status controller to the manager.
func AddClustersStatusController(mgr ctrl.Manager, producer producer.Producer, leafHubName string,
	incarnation uint64, hubOfHubsConfig *configV1.Config, syncIntervals *syncintervals.SyncIntervals,
) error {
	createObjFunction := func() bundle.Object { return &clusterV1.ManagedCluster{} }
	transportBundleKey := fmt.Sprintf("%s.%s", leafHubName, datatypes.ManagedClustersMsgKey)

	// update bundle object
	manipulateObjFunc := func(object bundle.Object) {
		helper.AddAnnotations(object, map[string]string{
			managedClusterManagedByAnnotation: leafHubName,
		})
	}

	predicateFunc := func() bool { // bundle predicate
		return hubOfHubsConfig.Spec.AggregationLevel == configV1.Full ||
			hubOfHubsConfig.Spec.AggregationLevel == configV1.Minimal
		// at this point send all managed clusters even if aggregation level is minimal
	}

	bundleCollection := []*generic.BundleCollectionEntry{ // single bundle for managed clusters
		generic.NewBundleCollectionEntry(transportBundleKey,
			bundle.NewGenericStatusBundle(leafHubName, incarnation, manipulateObjFunc),
			predicateFunc),
	}

	if err := generic.NewGenericStatusSyncController(mgr, clusterStatusSyncLogName, producer, bundleCollection,
		createObjFunction, nil, syncIntervals.GetManagerClusters); err != nil {
		return fmt.Errorf("failed to add managed clusters controller to the manager - %w", err)
	}

	return nil
}
