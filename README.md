[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Hub-of-Hubs-Agent

[![Go Reference](https://pkg.go.dev/badge/github.com/stolostron/hub-of-hubs-agent.svg)](https://pkg.go.dev/github.com/stolostron/hub-of-hubs-agent)
[![License](https://img.shields.io/github/license/stolostron/hub-of-hubs-agent)](/LICENSE)

The manager component of [Hub-of-Hubs](https://github.com/stolostron/hub-of-hubs).

Go to the [Contributing guide](CONTRIBUTING.md) to learn how to get involved.

## Prerequisites
- Hub of hubs on the ACM 2.5 environment
- Leaf hub to running the agent

## Getting Started

## Environment variables

The following environment variables are required for the most tasks below:

* `REGISTRY`, for example `quay.io/open-cluster-management-hub-of-hubs`.
* `IMAGE_TAG`, for example `latest` or `v0.1.0`.
* `LEAF_HUB_NAME`, the leaf hub name in hub-of-hubs.

## Build Image

```bash
make build-images
docker push ${REGISTRY}/hub-of-hubs-agent/${IMAGE_TAG}
```

## Disable the leaf hub addon controller from hub-of-hubs cluster

```
oc scale deploy hub-of-hubs-addon-controller -n open-cluster-management --replicas=0
oc delete manifestwork ${LEAF_HUB_NAME}-hoh-agent -n ${LEAF_HUB_NAME}
```

## Deploy the hub-of-hubs-agent on leaf hub cluster

Set the following environment variables in ./deploy/hub-of-hubs-agent.yaml

* IMAGE - the hub-of-hubs-agent controller, ${REGISTRY}/hub-of-hubs-agent/${IMAGE_TAG}
* POD_NAMESPACE - the leader election namespace
* WATCH_NAMESPACE - the watched namespaces, multiple namespace splited by comma.
* TRANSPORT_TYPE - the transport type, "kafka" or "sync-service"
* TRANSPORT_MESSAGE_COMPRESSION_TYPE - the compression type for transport message, "gzip" or "no-op"
* KAFKA_PRODUCER_ID - the ID for kafka producer, leaf hub cluster name by default
* KAFKA_PRODUCER_TOPIC - the topic for kafka producer at hub-of-hubs side, should be "spec"
* KAFKA_CONSUMER_TOPIC - the topic for kafka consumer at hub-of-hubs side, should be "status"
* KAFKA_BOOTSTRAP_SERVERS - the bootstrap server of kafka, "kafka-brokers-cluster-kafka-bootstrap.kafka.svc:9092" by default
* KAFKA_MESSAGE_SIZE_LIMIT_KB - the limit size of kafka message, should be < 1000
* SYNC_SERVICE_PROTOCOL - the protocol of sync-service, "http" by default
* SYNC_SERVICE_HOST - the host of of sync-service, "sync-service-css.sync-service.svc.cluster.local" by default
* SYNC_SERVICE_PORT - the port of of sync-service
* SYNC_SERVICE_POLLING_INTERVAL - the polling interval of sync-service
* K8S_CLIENTS_POOL_SIZE - the goroutine number to propagate the bundles on managed cluster
* SYNC_SERVICE_POLLING_INTERVAL - the interval of spec-sync
* KAFKA_SSL_CA - the authentication to connect to the kafka.
* COMPLIANCE_STATUS_DELTA_COUNT_SWITCH_FACTOR - "10" 

<!-- `POD_NAMESPACE` should usually be `open-cluster-management`.

`WATCH_NAMESPACE` can be defined empty so the controller will watch all the namespaces.

-->

```bash
oc deploy ./deploy/hub-of-hubs-rbac.yaml
oc deploy ./deploy/hub-of-hubs-agent.yaml
```


## TEST
[Hub of hubs Scenarios](https://docs.google.com/document/d/1a0cX_FzMXQaMM_HyBolXir2ce7XAab-re_KrQE9gp-c/edit)
