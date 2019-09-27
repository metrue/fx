#!/bin/bash

set -e

has() {
  type "$1" > /dev/null 2>&1
}

CLUSTER_FOR_TEST="fx-test"

kind
## install kind
if has "kind"; then
  echo "====== kind is ready, skip installation"
else
  GO111MODULE="on" go get sigs.k8s.io/kind@v0.5.1
fi

clusters=$(kind get clusters)
if [[ $clusters == *${CLUSTER_FOR_TEST}* ]]; then
  echo "====== Skip create cluster since ${CLUSTER_FOR_TEST} is created"
else
  kind create cluster --name fx-test --wait 300s
fi

export KUBECONFIG="$(kind get kubeconfig-path --name ${CLUSTER_FOR_TEST})"
echo "====== KUBECONFIG is $KUBECONFIG"
