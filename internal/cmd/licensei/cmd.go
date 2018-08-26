package licensei

import "github.com/spf13/cobra"

// AddCommands adds licensei commands to a Cobra command.
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListCommand(),
		NewCheckCommand(),
		NewCacheCommand(),
	)
}
