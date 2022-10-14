package main

import (
	"github.com/cloudbit-ch/cli/v2/internal/commands"
	"github.com/cloudbit-ch/cli/v2/internal/commands/common"
	"github.com/cloudbit-ch/cli/v2/internal/commands/compute"
	"github.com/cloudbit-ch/cli/v2/internal/commands/kubernetes"
	"github.com/cloudbit-ch/cli/v2/internal/commands/objectstorage"
)

var Version string

func main() {
	app := commands.Application{
		Name:        "cloudbit",
		Description: "cloudbit is a command-line interface for managing the Cloudbit cloud platform.",
		Version:     Version,
		Endpoint:    "https://api.cloudbit.ch/",

		Modules: []commands.ModuleFactory{
			common.Location,
			common.Module,
			common.Product,

			compute.Module,
			kubernetes.Module,
			objectstorage.Module,
		},
	}

	commands.Run(app)
}
