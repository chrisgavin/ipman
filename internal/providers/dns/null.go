package dns

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
	Name string `yaml:"-"`
}

func (provider *NullProvider) GetName(ctx context.Context) string {
	return provider.Name
}

func (provider *NullProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DNSAction, error) {
	current := []intermediates.DNSRecord{}
	desired := generators.HostsToRecords(hosts, nil)
	changes := diff.CompareDNSRecords(current, desired)
	return changes.ToActions(), nil
}

func (provider *NullProvider) ApplyAction(ctx context.Context, action actions.DNSAction) error {
	return nil
}
