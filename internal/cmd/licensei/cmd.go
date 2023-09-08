package licensei

import (
	"github.com/spf13/cobra"
)

// AddCommands adds licensei commands to a Cobra command.
func AddCommands(cmd *cobra.Command, globalOptions *GlobalOptions) {
	cmd.AddCommand(
		NewListCommand(globalOptions),
		NewCheckCommand(globalOptions),
		NewCacheCommand(globalOptions),
		NewHeaderCommand(globalOptions),
		NewStatCommand(globalOptions),
	)
}
