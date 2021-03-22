#!/bin/sbash
set -eu

KUBECTL_VERSION=1.17.0
KYML_VERSION=20190906

mkdir -p "$HOME/bin"
PATH="$PATH":"$HOME/bin"
curl -sfL \
	-o "$HOME/bin/kubectl" \
	https://storage.googleapis.com/kubernetes-release/release/v$KUBECTL_VERSION/bin/linux/amd64/kubectl
chmod +x "$HOME/bin/kubectl"
curl -sfL \
	-o "$HOME/bin/kyml" \
	https://github.com/frigus02/kyml/releases/download/v$KYML_VERSION/kyml_${KYML_VERSION}_linux_amd64
chmod +x "$HOME/bin/kyml"

echo "$GCP_SVC_ACC" >"$HOME/gcloud.json"
gcloud auth activate-service-account --key-file="$HOME/gcloud.json"
gcloud container clusters get-credentials website --zone europe-west3-a --project me-kuehle

kyml cat deploy/k8s/*.yml |
	kyml tmpl -v ImageTag=$(git rev-parse HEAD) |
	kyml resolve |
	kubectl apply -f -
