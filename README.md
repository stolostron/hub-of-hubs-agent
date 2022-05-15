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

## Disable the leaf hub addon controller from hub-of-hubs cluster

```
oc scale deploy hub-of-hubs-addon-controller -n open-cluster-management --replicas=0
oc delete manifestwork ${LEAF_HUB_NAME}-hoh-agent -n ${LEAF_HUB_NAME}
```

## Deploy the hub-of-hubs-agent on leaf hub cluster

The following environment variables are required for the most tasks below:

* `REGISTRY`, for example `quay.io/open-cluster-management-hub-of-hubs`.
* `IMAGE_TAG`, for example `latest` or `v0.1.0`.
* `LEAF_HUB_NAME`, the leaf hub name
* `KAFKA_BOOTSTRAP_SERVER`, the bootstrap server of kafka
* `KAFKA_SSL_CA`, the authentication to connect to the kafka.

### Deploy using the image
- Build Image
```bash
make build-images
docker push ${REGISTRY}/hub-of-hubs-agent:${IMAGE_TAG}
```
- Run 
```bash
oc apply -n open-cluster-management -f ./deploy/hub-of-hubs-rbac.yaml
envsub < ./deploy/hub-of-hubs-agent.yaml | oc apply -n open-cluster-management -f -
```

### Or Deploy using the binary file
```
./bin/hub-of-hubs-agent --kubeconfig=$LEAF_HUB_CONFIG --leaf-hub-name=${LEAF_HUB_NAME} --kafka-bootstrap-server=${KAFKA_BOOTSTRAP_SERVER} --kafka-ssl-ca=${KAFKA_SSL_CA}
```

## TEST
[Hub of hubs Scenarios](https://docs.google.com/document/d/1a0cX_FzMXQaMM_HyBolXir2ce7XAab-re_KrQE9gp-c/edit)
