package processor

import (
	"context"

	"github.com/chrisgavin/ipman/internal/types"
	"go.uber.org/zap"
)

func ProcessDHCP(ctx context.Context, input *types.Input, apply bool, logger *zap.Logger) ([]Change, error) {
	changes := []Change{}
	for _, provider := range input.DHCPProviders {
		logger.Info("Processing changes for provider.", zap.String("provider", provider.GetName(ctx)))
		for _, network := range input.Networks {
			for _, site := range network.Sites {
				for _, pool := range site.Pools {
					if !providerIncluded(provider.GetName(ctx), network, site, pool) {
						continue
					}
					actions, err := provider.GetActions(ctx, network, site, pool, pool.Hosts)
					if err != nil {
						return nil, err
					}
					for _, action := range actions {
						changes = append(changes, Change{
							Kind:      dhcpKind,
							Operation: action.GetOperation(),
							Provider:  provider.GetName(ctx),
							Network:   network.Name,
							Site:      site.Name,
							Pool:      pool.Name,
							Name:      action.GetName(),
							Summary:   action.ToString(),
						})
						logger.Info(action.ToString())
						if apply {
							err := provider.ApplyAction(ctx, action)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
	}
	return changes, nil
}
