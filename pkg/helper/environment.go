package helper

import (
	"fmt"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/open-horizon/edge-utilities/logger/log"
	kafkaclient "github.com/stolostron/hub-of-hubs-kafka-transport/kafka-client"
)

const (
	LEAF_HUB_ID   = "LH_ID"
	POD_NAMESPACE = "POD_NAMESPACE"

	TRANSPORT_TYPE                     = "TRANSPORT_TYPE"
	TRANSPORT_MESSAGE_COMPRESSION_TYPE = "TRANSPORT_MESSAGE_COMPRESSION_TYPE"

	K8S_WORK_POOL_SIZE      = "K8S_CLIENTS_POOL_SIZE"
	KAFKA_BOOTSTRAP_SERVERS = "KAFKA_BOOTSTRAP_SERVERS"
	KAFKA_SSL_CA            = "KAFKA_SSL_CA"
	KAFKA_CONSUMER_TOPIC    = "KAFKA_CONSUMER_TOPIC"

	SYNC_SERVICE_PROTOCOL                  = "SYNC_SERVICE_PROTOCOL"
	SYNC_SERVICE_CONSUMER_HOST             = "SYNC_SERVICE_HOST"
	SYNC_SERVICE_CONSUMER_PORT             = "SYNC_SERVICE_PORT"
	SYNC_SERVICE_CONSUMER_POLLING_INTERVAL = "SYNC_SERVICE_POLLING_INTERVAL"
	SPEC_ENFORCE_HOH_RBAC                  = "ENFORCE_HOH_RBAC"

	KAFKA_PRODUCER_ID           = "KAFKA_PRODUCER_ID"
	KAFKA_PRODUCER_TOPIC        = "KAFKA_PRODUCER_TOPIC"
	KAFKA_MESSAGE_SIZE_LIMIT_KB = "KAFKA_MESSAGE_SIZE_LIMIT_KB"

	SYNC_SERVICE_PRODUCER_HOST = "SYNC_SERVICE_HOST"
	SYNC_SERVICE_PRODUCER_PORT = "SYNC_SERVICE_PORT"

	COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR = "COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR"

	// envVarKafkaProducerID       = "KAFKA_PRODUCER_ID"
	// envVarKafkaBootstrapServers = "KAFKA_BOOTSTRAP_SERVERS"
	// envVarKafkaTopic            = "KAFKA_TOPIC"
	// envVarKafkaSSLCA            = "KAFKA_SSL_CA"
	// envVarMessageSizeLimit      = "KAFKA_MESSAGE_SIZE_LIMIT_KB"
)

const (
	defaultK8sClientsPoolSize = 10
	maxMessageSizeLimit       = 987 // to make sure that the message size is below 1 MB.
)

type KafkaEnvironment struct {
	BootstrapServers     string
	SslCa                string
	ComsumerTopic        string
	ProducerId           string
	ProducerTopic        string
	ProducerMessageLimit int
}

type SyncServiceEnvironment struct {
	Protocol                string
	ConsumerHost            string
	ConsumerPort            int
	ConsumerPollingInterval int
	ProducerHost            string
	ProducerPort            int
}

type EnvironmentManager struct {
	LeafHubID                    string
	PodNameSpace                 string
	TransportType                string
	TransportCompressionType     string
	SpecWorkPoolSize             int
	SpecEnforceHohRbac           bool
	StatusDeltaCountSwitchFactor int
	Kafka                        KafkaEnvironment
	SyncService                  SyncServiceEnvironment
}

func NewEnvironmentManager() (*EnvironmentManager, error) {
	podNameSpace, exist := os.LookupEnv(POD_NAMESPACE)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", POD_NAMESPACE)
	}

	transportType, exist := os.LookupEnv(TRANSPORT_TYPE)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", TRANSPORT_TYPE)
	}

	leafHubId, exist := os.LookupEnv(LEAF_HUB_ID)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", LEAF_HUB_ID)
	}

	kafkaBootstrapServers, exist := os.LookupEnv(KAFKA_BOOTSTRAP_SERVERS)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_BOOTSTRAP_SERVERS)
	}

	kafkaSslBase64EncodedCertificate, exist := os.LookupEnv(KAFKA_SSL_CA)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_SSL_CA)
	}

	kafkaConsumerTopic, exist := os.LookupEnv(KAFKA_CONSUMER_TOPIC)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_CONSUMER_TOPIC)
	}

	syncServiceProtocol, exist := os.LookupEnv(SYNC_SERVICE_PROTOCOL)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_PROTOCOL)
	}

	syncServiceConsumerHost, exist := os.LookupEnv(SYNC_SERVICE_CONSUMER_HOST)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_CONSUMER_HOST)
	}

	syncServiceConsumerPort, exit := os.LookupEnv(SYNC_SERVICE_CONSUMER_PORT)
	if !exit {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_CONSUMER_PORT)
	}
	syncServiceValidConsumerPort, err := strconv.Atoi(syncServiceConsumerPort)
	if err != nil || syncServiceValidConsumerPort < 0 {
		return nil, fmt.Errorf("environment variable %q is not valid", SYNC_SERVICE_CONSUMER_PORT)
	}

	syncServiceConsumerPollingInterval, exist := os.LookupEnv(SYNC_SERVICE_CONSUMER_POLLING_INTERVAL)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_CONSUMER_POLLING_INTERVAL)
	}
	syncServiceValidConsumerPollingInterval, err := strconv.Atoi(syncServiceConsumerPollingInterval)
	if err != nil || syncServiceValidConsumerPollingInterval < 0 {
		return nil, fmt.Errorf("environment variable %q is not valid", SYNC_SERVICE_CONSUMER_POLLING_INTERVAL)
	}

	syncServiceProducerHost, exist := os.LookupEnv(SYNC_SERVICE_PRODUCER_HOST)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_PRODUCER_HOST)
	}

	syncServiceProducerPort, exit := os.LookupEnv(SYNC_SERVICE_PRODUCER_PORT)
	if !exit {
		return nil, fmt.Errorf("not found environment variable: %q", SYNC_SERVICE_PRODUCER_PORT)
	}
	syncServiceValidProducerPort, err := strconv.Atoi(syncServiceProducerPort)
	if err != nil || syncServiceValidConsumerPort < 0 {
		return nil, fmt.Errorf("environment variable %q is not valid", SYNC_SERVICE_PRODUCER_PORT)
	}

	k8sWorkPoolSize, exist := os.LookupEnv(K8S_WORK_POOL_SIZE)
	if !exist {
		log.Warning("not found environment variable %q", K8S_WORK_POOL_SIZE)
	}
	k8sValidWorkPoolSize, err := strconv.Atoi(k8sWorkPoolSize)
	if err != nil || k8sValidWorkPoolSize < 1 {
		log.Warning("environment variable %s is not valid, default work pool size %d", K8S_WORK_POOL_SIZE, defaultK8sClientsPoolSize)
		k8sValidWorkPoolSize = defaultK8sClientsPoolSize
	}

	enforceHohRbacString, exist := os.LookupEnv(SPEC_ENFORCE_HOH_RBAC)
	if !exist {
		return nil, fmt.Errorf("not found environment variable %q", SYNC_SERVICE_CONSUMER_POLLING_INTERVAL)
	}

	enforceHohRbac, err := strconv.ParseBool(enforceHohRbacString)
	if err != nil {
		return nil, fmt.Errorf("environment variable %q is not valid", SYNC_SERVICE_CONSUMER_POLLING_INTERVAL)
	}

	transportMessageCompressionType, exist := os.LookupEnv(TRANSPORT_MESSAGE_COMPRESSION_TYPE)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", TRANSPORT_MESSAGE_COMPRESSION_TYPE)
	}

	kafkaProducerId, exist := os.LookupEnv(KAFKA_PRODUCER_ID)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_PRODUCER_ID)
	}

	kafkaProducerTopic, exist := os.LookupEnv(KAFKA_PRODUCER_TOPIC)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_PRODUCER_TOPIC)
	}

	kafkaMessageSizeLimit, exist := os.LookupEnv(KAFKA_MESSAGE_SIZE_LIMIT_KB)
	if !exist {
		return nil, fmt.Errorf("not found environment variable: %q", KAFKA_MESSAGE_SIZE_LIMIT_KB)
	}
	kafkaValidMessageSizeLimit, err := strconv.Atoi(kafkaMessageSizeLimit)
	if err != nil || kafkaValidMessageSizeLimit < 0 {
		return nil, fmt.Errorf("environment variable %q is not valid", KAFKA_MESSAGE_SIZE_LIMIT_KB)
	}

	if kafkaValidMessageSizeLimit > maxMessageSizeLimit {
		return nil, fmt.Errorf("%s - size must not exceed %d", KAFKA_MESSAGE_SIZE_LIMIT_KB, maxMessageSizeLimit)
	}

	statusCountSwitchFactor, exit := os.LookupEnv(COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR)
	if !exit {
		return nil, fmt.Errorf("not found environment variable: %q", COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR)
	}
	statusValidCountSwitchFactor, err := strconv.Atoi(statusCountSwitchFactor)
	if err != nil || statusValidCountSwitchFactor < 0 {
		return nil, fmt.Errorf("environment variable %q is not valid", COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR)
	}

	environmentManager := &EnvironmentManager{
		LeafHubID:                    leafHubId,
		PodNameSpace:                 podNameSpace,
		TransportType:                transportType,
		TransportCompressionType:     transportMessageCompressionType,
		SpecWorkPoolSize:             k8sValidWorkPoolSize,
		SpecEnforceHohRbac:           enforceHohRbac,
		StatusDeltaCountSwitchFactor: statusValidCountSwitchFactor,
		Kafka: KafkaEnvironment{
			BootstrapServers:     kafkaBootstrapServers,
			SslCa:                kafkaSslBase64EncodedCertificate,
			ComsumerTopic:        kafkaConsumerTopic,
			ProducerId:           kafkaProducerId,
			ProducerTopic:        kafkaProducerTopic,
			ProducerMessageLimit: kafkaValidMessageSizeLimit,
		},
		SyncService: SyncServiceEnvironment{
			Protocol:                syncServiceProtocol,
			ConsumerHost:            syncServiceConsumerHost,
			ConsumerPort:            syncServiceValidConsumerPort,
			ConsumerPollingInterval: syncServiceValidConsumerPollingInterval,
			ProducerHost:            syncServiceProducerHost,
			ProducerPort:            syncServiceValidProducerPort,
		},
	}
	return environmentManager, nil
}

func (environmentManager *EnvironmentManager) GetKafkaConfigMap() (*kafka.ConfigMap, error) {
	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers":       environmentManager.Kafka.BootstrapServers,
		"client.id":               environmentManager.LeafHubID,
		"group.id":                environmentManager.LeafHubID,
		"auto.offset.reset":       "earliest",
		"enable.auto.commit":      "false",
		"socket.keepalive.enable": "true",
		"log.connection.close":    "false", // silence spontaneous disconnection logs, kafka recovers by itself.
	}
	err := loadSslToConfigMap(kafkaConfigMap)
	if err != nil {
		return kafkaConfigMap, fmt.Errorf("failed to configure kafka-consumer - %w", err)
	}
	return kafkaConfigMap, nil
}

func (environmentManager *EnvironmentManager) GetProducerKafkaConfigMap() (*kafka.ConfigMap, error) {
	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers":       environmentManager.Kafka.BootstrapServers,
		"client.id":               environmentManager.Kafka.ProducerId,
		"acks":                    "1",
		"retries":                 "0",
		"socket.keepalive.enable": "true",
		"log.connection.close":    "false", // silence spontaneous disconnection logs, kafka recovers by itself.
	}

	err := loadSslToConfigMap(kafkaConfigMap)
	if err != nil {
		return kafkaConfigMap, fmt.Errorf("failed to configure kafka-producer - %w", err)
	}
	return kafkaConfigMap, nil
}

func loadSslToConfigMap(kafkaConfigMap *kafka.ConfigMap) error {
	if sslBase64EncodedCertificate, found := os.LookupEnv(KAFKA_SSL_CA); found {
		certFileLocation, err := kafkaclient.SetCertificate(&sslBase64EncodedCertificate)
		if err != nil {
			return fmt.Errorf("failed to SetCertificate - %w", err)
		}

		if err = kafkaConfigMap.SetKey("security.protocol", "ssl"); err != nil {
			return fmt.Errorf("failed to SetKey security.protocol - %w", err)
		}

		if err = kafkaConfigMap.SetKey("ssl.ca.location", certFileLocation); err != nil {
			return fmt.Errorf("failed to SetKey ssl.ca.location - %w", err)
		}
	}
	return nil
}
