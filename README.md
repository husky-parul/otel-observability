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

### Get Kepler Prom metrics into otel collector via logging

```

kubectl get pods

NAME                              READY   STATUS    RESTARTS   AGE
kepler-exporter-vg4md             1/1     Running   0          65m
otel-collector-5fbbf4d69b-lwxvk   1/1     Running   0          9m53s


oc logs otel-collector-5fbbf4d69b-lwxvk

StartTimestamp: 2023-09-29 13:21:49.808 +0000 UTC
Timestamp: 2023-09-29 13:31:54.808 +0000 UTC
Value: 0.000000
NumberDataPoints #97
Data point attributes:
     -> container_id: Str(c75251098943186b0671bb5b16d1e85423e6cc36582ca8739320d05b828945e7)
     -> container_name: Str(init-config-reloader)
     -> container_namespace: Str(monitoring)
     -> mode: Str(dynamic)
     -> pod_name: Str(prometheus-k8s-1)
StartTimestamp: 2023-09-29 13:21:49.808 +0000 UTC
Timestamp: 2023-09-29 13:31:54.808 +0000 UTC
Value: 0.000000
NumberDataPoints #98
Data point attributes:
     -> container_id: Str(c82bfdaeb7f86417f0f1d76d4f70751a1e28016a405a56832c1a56cbaf4e2960)
     -> container_name: Str(init-config-reloader)
     -> container_namespace: Str(monitoring)
     -> mode: Str(dynamic)
     -> pod_name: Str(alertmanager-main-2)
StartTimestamp: 2023-09-29 13:21:49.808 +0000 UTC
Timestamp: 2023-09-29 13:31:54.808 +0000 UTC
Value: 0.000000
NumberDataPoints #99
Data point attributes:
     -> container_id: Str(cb47a077fb14efea4f9744de5a4e710cf3e1a1ac498a50d667155d3255ce2c25)
     -> container_name: Str(init-config-reloader)
     -> container_namespace: Str(monitoring)
     -> mode: Str(dynamic)
     -> pod_name: Str(alertmanager-main-1)
StartTimestamp: 2023-09-29 13:21:49.808 +0000 UTC
Timestamp: 2023-09-29 13:31:54.808 +0000 UTC
Value: 0.000000
NumberDataPoints #100
Data point attributes:
     -> container_id: Str(cbf8e648e2ad453d1e20aec528fb695d6fcf2c009e86326add8d6c329399e1b1)
     -> container_name: Str(kube-rbac-proxy-self)
     -> container_namespace: Str(monitoring)
     -> mode: Str(dynamic)
     -> pod_name: Str(kube-state-metrics-59cfdf494-wf49t)



```