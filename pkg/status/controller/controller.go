// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"

	"github.com/stolostron/hub-of-hubs-agent/pkg/helper"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/apps"
	configCtrl "github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/config"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/controlinfo"
	localpolicies "github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/local_policies"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/localplacement"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/managedclusters"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/placement"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/policies"
	"github.com/stolostron/hub-of-hubs-agent/pkg/status/controller/syncintervals"
	"github.com/stolostron/hub-of-hubs-agent/pkg/transport/producer"
	configV1 "github.com/stolostron/hub-of-hubs-data-types/apis/config/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// AddControllers adds all the controllers to the Manager.
func AddControllers(mgr ctrl.Manager, pro producer.Producer, env helper.EnvironmentManager, incarnation uint64) error {
	config := &configV1.Config{}
	if err := configCtrl.AddConfigController(mgr, config); err != nil {
		return fmt.Errorf("failed to add ConfigMap controller: %w", err)
	}

	syncIntervals := syncintervals.NewSyncIntervals()
	if err := syncintervals.AddSyncIntervalsController(mgr, syncIntervals); err != nil {
		return fmt.Errorf("failed to add SyncIntervals controller: %w", err)
	}

	if err := policies.AddPoliciesStatusController(mgr, pro, env, incarnation, config, syncIntervals); err != nil {
		return fmt.Errorf("failed to add PoliciesStatusController controller: %w", err)
	}

	addControllerFunctions := []func(ctrl.Manager, producer.Producer, string, uint64, *configV1.Config, *syncintervals.SyncIntervals) error{
		managedclusters.AddClustersStatusController,
		placement.AddPlacementRulesController,
		placement.AddPlacementsController,
		placement.AddPlacementDecisionsController,
		apps.AddSubscriptionStatusesController,
		apps.AddSubscriptionReportsController,
		localpolicies.AddLocalPoliciesController,
		localplacement.AddLocalPlacementRulesController,
		controlinfo.AddControlInfoController,
	}

	for _, addControllerFunction := range addControllerFunctions {
		if err := addControllerFunction(mgr, pro, env.LeafHubID, incarnation, config, syncIntervals); err != nil {
			return fmt.Errorf("failed to add controller: %w", err)
		}
	}

	return nil
}
