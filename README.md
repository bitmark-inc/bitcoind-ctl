# Bitcoind Controller

Bitcoind controller manages bitcoind process' lifecycle either by UNIX subprocess in local environment or by kubernetes API in cluster environment.

The functionality of this service is mainly to control the start/stop of bitcoind. It assists the pod controller to manage the bitcoind lifecycle through API endpoints. The reason to create this separate serivce instead of

## Local Mode

In local mode, it controls the start/stop of bitcoind by managing the local bitcoind process.

## Kubernetes Mode

In cluster mode, it controls the statefulset of its corresponded bitcoind through with in a Kubernetes cluster. It requires the `PATCH` permission of its related statefulset since it controls the bitcoind's existence by updating the replica number in the statefulset.


### Environment Variables

#### Global

| Variable | Type | Description |
|-|-|-|
| AUTONOMY_WALLET_OWNER_DID | string | Autonomy client's DID |
| AUTONOMY_WALLET_BITCOIND_NETOWORK | string | bitcoind network |
| AUTONOMY_WALLET_BITCOIND_ENDPOINT | string | bitcoind endpoint |

#### Kubernetes

| Variable | Type | Description |
|-|-|-|
| AUTONOMY_WALLET_K8S_USE_LOCAL_CONTEXT | boolean | Use local Kubernetes context   |
| AUTONOMY_WALLET_K8S_NAMESPACE | string | Kubernetes namespace |

#### Local

| Variable | Type | Description |
|-|-|-|
| AUTONOMY_WALLET_LOCAL_BITCOIND_PATH | string | the path of bitcoind binary |
| AUTONOMY_WALLET_LOCAL_BITCOIND_CONF_PATH | string | the path of bitcoind configuration |
