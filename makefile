# Generate manifests for CRDs

.PHONY: manifests

CONTROLLER_GEN = "controller-gen"
all: gen crd
gen:
	@hack/update-codegen.sh

crd:
	$(CONTROLLER_GEN) crd:trivialVersions=true paths="./..." output:crd:artifacts:config=config/crd/bases
applycrd:

	@kubectl apply -f config/crd/bases/heheh.com_destroyments.yaml
	@kubectl apply -f config/crd/bases/destroyment.yaml

