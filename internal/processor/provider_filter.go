package processor

import "github.com/chrisgavin/ipman/internal/types"

func providerIncluded(provider string, network types.Network, site types.Site, pool types.Pool) bool {
	for _, p := range network.Providers {
		if p == provider {
			return true
		}
	}
	for _, p := range site.Providers {
		if p == provider {
			return true
		}
	}
	for _, p := range pool.Providers {
		if p == provider {
			return true
		}
	}
	return false
}
