package commands

import (
	"github.com/chrisgavin/ipman/internal/input"
	"github.com/chrisgavin/ipman/internal/processor"
	"github.com/spf13/cobra"
)

type ApplyCommand struct {
	*RootCommand
}

func registerApplyCommand(rootCommand *RootCommand) {
	command := &ApplyCommand{
		RootCommand: rootCommand,
	}
	applyCommand := &cobra.Command{
		Use:           "apply",
		Short:         "Apply the desired state.",
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

			err = processor.ProcessDNS(cmd.Context(), input, true, command.logger)
			if err != nil {
				return err
			}
			err = processor.ProcessDHCP(cmd.Context(), input, true, command.logger)
			if err != nil {
				return err
			}

			return nil
		},
	}
	command.root.AddCommand(applyCommand)
}
