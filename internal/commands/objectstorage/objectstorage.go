package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func AddCommands(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "object-storage",
		Aliases: []string{"objectstorage"},
		Short:   "Manage your object storage",
	}

	cmd.AddCommand(
		InstanceCommand(),
	)

	parent.AddCommand(cmd)
}

func init() {
	AddCommands(&commands.Root)
}
