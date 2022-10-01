package commands

import (
	"github.com/chrisgavin/ipman/internal/input"
	"github.com/spf13/cobra"
)

type CheckCommand struct {
	*RootCommand
}

func registerCheckCommand(rootCommand *RootCommand) {
	command := &CheckCommand{
		RootCommand: rootCommand,
	}
	checkCommand := &cobra.Command{
		Use:           "check",
		Short:         "Validate whether a given configuration is valid.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := input.ReadInput(command.input)
			if err != nil {
				return err
			}
			err = input.Validate()
			if err != nil {
				return err
			}
			command.logger.Info("Configuration is valid.")
			return nil
		},
	}
	command.root.AddCommand(checkCommand)
}
