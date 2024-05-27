package commands

import (
	"github.com/chrisgavin/ipman/internal/input"
	"github.com/chrisgavin/ipman/internal/processor"
	"github.com/spf13/cobra"
)

type DiffCommand struct {
	*RootCommand
}

func registerDiffCommand(rootCommand *RootCommand) {
	command := &DiffCommand{
		RootCommand: rootCommand,
	}
	diffCommand := &cobra.Command{
		Use:           "diff",
		Short:         "Show the difference between the current state and the desired state.",
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

			err = processor.ProcessDNS(cmd.Context(), input, false, command.logger)
			if err != nil {
				return err
			}
			err = processor.ProcessDHCP(cmd.Context(), input, false, command.logger)
			if err != nil {
				return err
			}

			return nil
		},
	}
	command.root.AddCommand(diffCommand)
}
