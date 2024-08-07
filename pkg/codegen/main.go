package main

import (
	"os"

	longhornv1 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	controllergen "github.com/rancher/wrangler/v3/pkg/controller-gen"
	"github.com/rancher/wrangler/v3/pkg/controller-gen/args"

	netfsv1 "github.com/Vicente-Cheng/networkfs-manager/pkg/apis/harvesterhci.io/v1beta1"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/Vicente-Cheng/networkfs-manager/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"harvesterhci.io": {
				Types: []interface{}{
					netfsv1.NetworkFilesystem{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
			longhornv1.SchemeGroupVersion.Group: {
				Types: []interface{}{
					longhornv1.ShareManager{},
				},
				GenerateTypes:   false,
				GenerateClients: true,
			},
		},
	})
}
