package commands

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/input"
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

			for _, provider := range input.DNSProviders {
				command.logger.Info(fmt.Sprintf("Processing changes for provider %T.", provider))
				for _, network := range input.Networks {
					for _, site := range network.Sites {
						for _, pool := range site.Pools {
							actions, err := provider.GetActions(cmd.Context(), network, site, pool, pool.Hosts)
							if err != nil {
								return err
							}
							for _, action := range actions {
								command.logger.Info(action.ToString())
								err := provider.ApplyAction(cmd.Context(), action)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}

			for _, provider := range input.DHCPProviders {
				command.logger.Info(fmt.Sprintf("Processing changes for provider %T.", provider))
				for _, network := range input.Networks {
					for _, site := range network.Sites {
						for _, pool := range site.Pools {
							actions, err := provider.GetActions(cmd.Context(), network, site, pool, pool.Hosts)
							if err != nil {
								return err
							}
							for _, action := range actions {
								command.logger.Info(action.ToString())
								err := provider.ApplyAction(cmd.Context(), action)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}

			return nil
		},
	}
	command.root.AddCommand(applyCommand)
}
