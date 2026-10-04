package processor

import (
	"context"
	"testing"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
	"go.uber.org/zap"
)

// siteScopedProvider mirrors how a real DNS provider reconciles: its state is
// scoped to the site, because pool names form no part of a record's name.
type siteScopedProvider struct {
	current []intermediates.DNSRecord
}

func (provider *siteScopedProvider) GetName(ctx context.Context) string {
	return "test"
}

func (provider *siteScopedProvider) GetActions(ctx context.Context, network types.Network, site types.Site) ([]actions.DNSAction, error) {
	current := generators.RecordsForSite(network, site, provider.current)
	desired := generators.HostsToRecords(site.Hosts(), nil)
	changes := diff.CompareDNSRecords(current, desired)
	return changes.ToActions(), nil
}

func (provider *siteScopedProvider) ApplyAction(ctx context.Context, action actions.DNSAction) error {
	return nil
}

func TestPoolWithNoHostsDoesNotDeleteSiblingPoolRecords(t *testing.T) {
	network := types.Network{Name: "example.com", Providers: []string{"test"}}
	site := types.Site{Name: "site-01", Network: &network}
	static := types.Pool{Name: "static", Site: &site}
	static.Hosts = []types.Host{{
		Name:       "host-01",
		Interfaces: []types.Interface{{Name: "eth0", Address: "10.0.0.1"}},
		Pool:       &static,
	}}
	empty := types.Pool{Name: "dynamic", Site: &site}
	site.Pools = []types.Pool{empty, static}
	network.Sites = []types.Site{site}

	provider := &siteScopedProvider{current: generators.HostsToRecords(static.Hosts, nil)}
	input := types.Input{
		Networks:     []types.Network{network},
		DNSProviders: []types.DNSProvider{provider},
	}

	changes, err := ProcessDNS(context.Background(), &input, false, zap.NewNop())
	if err != nil {
		t.Fatalf("ProcessDNS returned an unexpected error: %v", err)
	}
	for _, change := range changes {
		if change.Operation != actions.OperationCreate {
			t.Errorf("Expected no %s changes, but got one for %s.", change.Operation, change.Name)
		}
	}
	if len(changes) != 0 {
		t.Errorf("Expected no changes, but got %d.", len(changes))
	}
}
