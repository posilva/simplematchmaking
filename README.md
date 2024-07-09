# Matchmaking for games

## TODO


## Infra needed for AWS 


## Test locally with Kubernetes
### install k9s
`brew install derailed/k9s/k9s`
### Setup minikube

Ref: https://minikube.sigs.k8s.io/docs/start

- `curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-arm64`
- `sudo install minikube-darwin-arm64 /usr/local/bin/minikube`
- start the cluster: `minikube start -p simplematchmaking`
- install helm `minikube addons enable helm-tiller -p default`
- `minikube addons enable yakd -p default `
- `minikube -p default service yakd-dashboard -n yakd-dashboard`
- ` minikube -p default addons enable metrics-server`
- `minikube -p default dashboard --url`
