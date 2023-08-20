package dhcp

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
)

type NullProvider struct {
	Type string
}

func (provider *NullProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DHCPAction, error) {
	current := []intermediates.DHCPReservation{}
	desired := generators.HostsToReservations(network, site, pool, hosts)
	changes := diff.CompareDHCPReservations(current, desired)
	return changes.ToActions(), nil
}
