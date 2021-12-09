# GKE Hubble Export

This is a grpc server wrapper that re-export the cilium agent's observer service and peer service
from the local domain socket. And it can be consumed by the Hubble Relay on the GKE Dataplane V2.

## Notice

The current implementation re-exports the observer service without mTLS.
Feel free to modify it.

## Example 

```shell
# create gke with --enable-dataplane-v2
gcloud beta container clusters create "cluster-1" \
  --cluster-version "1.21.5-gke.1302" \
  --enable-dataplane-v2 

# deploy gke-hubble-export + hubble-relay + hubble-ui
kubectl apply -f example.yaml

# access the hubble-ui
kubectl port-forward svc/hubble-ui 8081:80

```