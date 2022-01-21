# GKE Hubble Export

Currently, Hubble UI can't be install on the GKE Dataplane V2 because GKE doesn't export cilium agent's observer service and peer service through tcp socket which is required by the Hubble Relay.

So, this repo is just a simple grpc proxy that re-export the above services from the local domain socket to a tcp socket. And it can be consumed by the Hubble Relay on the GKE Dataplane V2, thus Hubble UI can be install.

## Warning

The current implementation re-exports the above services without mTLS.
Feel free to modify it.

## Example

```shell
# create gke with --enable-dataplane-v2
gcloud beta container clusters create "cluster-1" \
  --cluster-version "1.21.6-gke.1500" \
  --enable-dataplane-v2

# deploy gke-hubble-export + hubble-relay + hubble-ui in the default namespace.
kubectl apply -f example.yaml

# access the hubble-ui
kubectl port-forward svc/hubble-ui 8081:80

```