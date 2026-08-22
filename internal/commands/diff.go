package commands

import (
	"encoding/json"
	"os"

	"github.com/chrisgavin/ipman/internal/input"
	"github.com/chrisgavin/ipman/internal/processor"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	textFormat = "text"
	jsonFormat = "json"
)

type Difference struct {
	Changes []processor.Change `json:"changes"`
}

type DiffCommand struct {
	*RootCommand
	format string
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
			if command.format != textFormat && command.format != jsonFormat {
				return errors.Errorf("Unknown format %s.", command.format)
			}

			input, err := input.ReadInput(command.input)
			if err != nil {
				return err
			}

			err = input.Validate()
			if err != nil {
				return err
			}

			dnsChanges, err := processor.ProcessDNS(cmd.Context(), input, false, command.logger)
			if err != nil {
				return err
			}
			dhcpChanges, err := processor.ProcessDHCP(cmd.Context(), input, false, command.logger)
			if err != nil {
				return err
			}

			difference := Difference{Changes: append(dnsChanges, dhcpChanges...)}
			if len(difference.Changes) == 0 {
				command.logger.Info("No changes required.")
			} else {
				command.logger.Info("Changes required.", zap.Int("changes", len(difference.Changes)))
			}

			if command.format == jsonFormat {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "\t")
				if err := encoder.Encode(difference); err != nil {
					return errors.Wrap(err, "Error writing difference.")
				}
			}

			return nil
		},
	}
	diffCommand.Flags().StringVar(&command.format, "format", textFormat, "The format to output the difference in. One of \"text\" or \"json\".")
	command.root.AddCommand(diffCommand)
}
