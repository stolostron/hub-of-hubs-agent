package placement

import (
	"fmt"

	"github.com/stolostron/hub-of-hubs-agent/pkg/helper"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/bundle"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/generic"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/syncintervals"
	"github.com/stolostron/hub-of-hubs-agent/pkg/transport/producer"
	configV1 "github.com/stolostron/hub-of-hubs-data-types/apis/config/v1"
	placementrulesV1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	placementRuleSyncLog           = "placement-rules-sync"
	PlacementRuleMsgKey            = "PlacementRule"
	OriginOwnerReferenceAnnotation = "hub-of-hubs.open-cluster-management.io/originOwnerReferenceUid"
)

// AddPlacementRulesController adds placement-rule controller to the manager.
func AddPlacementRulesController(mgr ctrl.Manager, transport producer.Producer, leafHubName string,
	incarnation uint64, _ *configV1.Config, syncIntervalsData *syncintervals.SyncIntervals,
) error {
	createObjFunction := func() bundle.Object { return &placementrulesV1.PlacementRule{} }

	// TODO datatypes.PlacementRuleMsgKey
	bundleCollection := []*generic.BundleCollectionEntry{
		generic.NewBundleCollectionEntry(fmt.Sprintf("%s.%s", leafHubName, PlacementRuleMsgKey),
			bundle.NewGenericStatusBundle(leafHubName, incarnation, cleanPlacementRule),
			func() bool { return true }),
	} // bundle predicate - always send placement rules.

	// TODO datatypes.OriginOwnerReferenceAnnotation
	ownerRefAnnotationPredicate := predicate.NewPredicateFuncs(func(object client.Object) bool {
		return helper.HasAnnotation(object, OriginOwnerReferenceAnnotation)
	})

	if err := generic.NewGenericStatusSyncController(mgr, placementRuleSyncLog, transport, bundleCollection,
		createObjFunction, ownerRefAnnotationPredicate, syncIntervalsData.GetPolicies); err != nil {
		return fmt.Errorf("failed to add placement rules controller to the manager - %w", err)
	}

	return nil
}

func cleanPlacementRule(object bundle.Object) {
	placementrule, ok := object.(*placementrulesV1.PlacementRule)
	if !ok {
		panic("Wrong instance passed to clean placement-rule function, not a placement-rule")
	}
	// clean spec. no need for it.
	placementrule.Spec = placementrulesV1.PlacementRuleSpec{}
}
