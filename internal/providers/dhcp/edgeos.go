package dhcp

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/clients/edgeosclient"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/pkg/errors"
)

type EdgeOSProvider struct {
	Type     string
	Name     string `yaml:"-"`
	Address  string
	Username string
	Password string
}

type EdgeOSProviderState struct {
	ReservationID string
}

func (provider *EdgeOSProvider) client() *edgeosclient.EdgeOSClient {
	return edgeosclient.NewEdgeOSClient(provider.Address, provider.Username, provider.Password)
}

func (provider *EdgeOSProvider) GetName(ctx context.Context) string {
	return provider.Name
}

func (provider *EdgeOSProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DHCPAction, error) {
	current := []intermediates.DHCPReservation{}

	client := provider.client()

	configuration := edgeosclient.EdgeOSConfiguration{}

	err := client.Get(&configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get DHCP leases.")
	}

	for leaseName, lease := range configuration.Get.Service.DHCPServer.SharedNetworkName["Local"].Subnet[pool.Range].StaticMapping {
		current = append(current, intermediates.DHCPReservation{
			ProviderState: EdgeOSProviderState{ReservationID: leaseName},
			Name:          leaseName,
			Address:       lease.IPAddress,
			MAC:           lease.MACAddress,
		})
	}

	desired := generators.HostsToReservations(hosts, nil)
	changes := diff.CompareDHCPReservations(current, desired)
	return changes.ToActions(), nil
}

func (provider *EdgeOSProvider) ApplyAction(ctx context.Context, action actions.DHCPAction) error {
	return errors.New("Not implemented.")
}
