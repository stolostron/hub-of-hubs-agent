package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/pflag"
	compressor "github.com/stolostron/hub-of-hubs-message-compression"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"

	helper "github.com/stolostron/hub-of-hubs-agent/pkg/helper"
	"github.com/stolostron/hub-of-hubs-agent/pkg/spec/bundle"
	"github.com/stolostron/hub-of-hubs-agent/pkg/spec/controller"
	consumer "github.com/stolostron/hub-of-hubs-agent/pkg/transport/consumer"
	producer "github.com/stolostron/hub-of-hubs-agent/pkg/transport/producer"
)

const (
	METRICS_HOST               = "0.0.0.0"
	METRICS_PORT               = 9435
	TRANSPORT_TYPE_KAFKA       = "kafka"
	TRANSPORT_TYPE_SYNC_SVC    = "sync-service"
	LEADER_ELECTION_ID         = "hub-of-hubs-agent-lock"
	HOH_LOCAL_NAMESPACE        = "hoh-local"
	INCARNATION_CONFIG_MAP_KEY = "incarnation"
	BASE10                     = 10
	UINT64_SIZE                = 64
)

func main() {
	os.Exit(doMain())
}

// function to handle defers with exit, see https://stackoverflow.com/a/27629493/553720.
func doMain() int {
	// TODO test and refactor
	pflag.CommandLine.AddFlagSet(zap.FlagSet())
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	ctrl.SetLogger(zap.Logger())

	log := ctrl.Log.WithName("cmd")
	printVersion(log)

	environmentManager, err := helper.NewEnvironmentManager()
	if err != nil {
		log.Error(err, "failed to load environment variable")
		return 1
	}

	// transport layer initialization
	genericBundleChan := make(chan *bundle.GenericBundle)
	defer close(genericBundleChan)

	consumer, err := getConsumer(environmentManager, genericBundleChan)
	if err != nil {
		log.Error(err, "transport initialization error")
		return 1
	}

	mgr, err := createManager(consumer, environmentManager)
	if err != nil {
		log.Error(err, "Failed to create manager")
		return 1
	}

	consumer.Start()
	defer consumer.Stop()

	log.Info("Starting the Cmd.")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		return 1
	}

	return 0
}

func printVersion(log logr.Logger) {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

// function to choose transport type based on env var.
func getConsumer(environmentManager *helper.EnvironmentManager, genericBundleChan chan *bundle.GenericBundle) (consumer.Consumer, error) {
	switch environmentManager.TransportType {
	case TRANSPORT_TYPE_KAFKA:
		kafkaConsumer, err := consumer.NewKafkaConsumer(ctrl.Log.WithName("kafka-consumer"), environmentManager, genericBundleChan)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka-consumer: %w", err)
		}
		return kafkaConsumer, nil
	case TRANSPORT_TYPE_SYNC_SVC:
		syncService, err := consumer.NewSyncService(ctrl.Log.WithName("sync-service"), environmentManager, genericBundleChan)
		if err != nil {
			return nil, fmt.Errorf("failed to create sync-service: %w", err)
		}
		return syncService, nil
	default:
		return nil, fmt.Errorf("environment variable %q - %q is not a valid option", "TRANSPORT_TYPE", environmentManager.TransportType)
	}
}

func getProducer(environmentManager *helper.EnvironmentManager) (producer.Producer, error) {
	messageCompressor, err := compressor.NewCompressor(compressor.CompressionType(environmentManager.TransportCompressionType))
	if err != nil {
		return nil, fmt.Errorf("failed to create transport producer message-compressor: %w", err)
	}

	switch environmentManager.TransportType {
	case TRANSPORT_TYPE_KAFKA:
		kafkaProducer, err := producer.NewKafkaProducer(messageCompressor, ctrl.Log.WithName("kafka-producer"), environmentManager)
		if err != nil {
			return nil, fmt.Errorf("failed to create kafka-producer: %w", err)
		}
		return kafkaProducer, nil
	case TRANSPORT_TYPE_SYNC_SVC:
		syncServiceProducer, err := producer.NewSyncServiceProducer(messageCompressor, ctrl.Log.WithName("syncservice-producer"), environmentManager)
		if err != nil {
			return nil, fmt.Errorf("failed to create sync-service producer: %w", err)
		}
		return syncServiceProducer, nil
	default:
		return nil, fmt.Errorf("environment variable %q - %q is not a valid option", "TRANSPORT_TYPE", environmentManager.TransportType)
	}
}

func createManager(consumer consumer.Consumer, environmentManager *helper.EnvironmentManager) (ctrl.Manager, error) {
	options := ctrl.Options{
		MetricsBindAddress:      fmt.Sprintf("%s:%d", METRICS_HOST, METRICS_PORT),
		LeaderElection:          true,
		LeaderElectionID:        LEADER_ELECTION_ID,
		LeaderElectionNamespace: environmentManager.PodNameSpace,
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new manager: %w", err)
	}

	if err := controller.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, fmt.Errorf("failed to add schemes: %w", err)
	}

	if err := controller.AddControllerToManager(mgr, consumer, *environmentManager); err != nil {
		return nil, fmt.Errorf("failed to add spec syncer: %w", err)
	}

	return mgr, nil
}
