package processor

import (
	"context"

	"github.com/chrisgavin/ipman/internal/types"
	"go.uber.org/zap"
)

func ProcessDNS(ctx context.Context, input *types.Input, apply bool, logger *zap.Logger) error {
	for _, provider := range input.DNSProviders {
		logger.Info("Processing changes for provider.", zap.String("provider", provider.GetName(ctx)))
		for _, network := range input.Networks {
			for _, site := range network.Sites {
				for _, pool := range site.Pools {
					if !providerIncluded(provider.GetName(ctx), network, site, pool) {
						continue
					}
					actions, err := provider.GetActions(ctx, network, site, pool, pool.Hosts)
					if err != nil {
						return err
					}
					for _, action := range actions {
						logger.Info(action.ToString())
						if apply {
							err := provider.ApplyAction(ctx, action)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}
