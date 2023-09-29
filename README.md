# otel-observability

## Set up Kepler on kind

```
git clone https://github.com/sustainable-computing-io/kepler.git

export CLUSTER_PROVIDER='kind'
export PROMETHEUS_ENABLE="true"
export GRAFANA_ENABLE="true"
make cluster-up

git clone --depth 1 https://github.com/prometheus-operator/kube-prometheus
cd kube-prometheus

kubectl apply --server-side -f manifests/setup
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
kubectl apply -f manifests/

```
> This will create the monitoring namespace and CRDs, and then wait for them to be available before creating the remaining resources. During the until loop, a response of No resources found is to be expected. This statement checks whether the resource API is created but doesn't expect the resources to be there.

## Create Kepler manifests

```
make build-manifest OPTS="CI_DEPLOY PROMETHEUS_DEPLOY DEBUG_DEPLOY"

kubectl apply -f _output/generated-manifest/deployment.yaml 

```

### Deploy Otel Collector

```
kubectl apply -f https://raw.githubusercontent.com/husky-parul/otel-observability/main/otel/otel-collector.yaml
```