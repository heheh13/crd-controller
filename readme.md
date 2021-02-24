# CRD
## goals
+ understanding custom resource
+ define a custom resource
+ create a controller for it


## generating custom code:

need to  create the types.go defining the struct for CRD
need to add to known object in register.go

need to  run generate-groups.sh

hack/update file will populated with 
directory where to gen
directory from where to gen
group and version

chmod +x `hack/update-codegen.sh`

chmod + `vendor/k8s.io/code-generator/generate-groups.sh`

`hack/update-codegen.sh`

## generated Code:

### deepcopy
    a deep copy for structured that defined in the types.go and commanded to make a deepcopy
### clientset
    creates the functionality to configure and do any REST  Operation
### informer
    offers an event based information and react based on changes
### lister
    offers readonly chaching layer for s get and list

## coments that generates codes:

+genclient

+genclient:noStatus

+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object