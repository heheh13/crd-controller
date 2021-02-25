#!/bin/bash


 vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/heheh13/crd/custom/client github.com/heheh13/crd/custom/apis \
  heheh.com:v1 \
  --go-header-file hack/boilerplate.go.txt
