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
	"github.com/seancfoley/ipaddress-go/ipaddr"
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
	Subnet        string
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

	configuration, err := client.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get DHCP leases.")
	}

	poolRange := ipaddr.NewIPAddressString(pool.Range)
	subnets := configuration.Service.DHCPServer.SharedNetworkName["Local"].Subnet
	var matchingSubnet *string
	for subnetRange := range subnets {
		parsedSubnetRange := ipaddr.NewIPAddressString(subnetRange)
		if parsedSubnetRange.Contains(poolRange) {
			if matchingSubnet != nil {
				return nil, errors.New("Multiple subnets contain the pool range " + pool.Range + ".")
			}
			matchingSubnet = &subnetRange
		}
	}
	if matchingSubnet == nil {
		return nil, errors.New("No subnet contains the pool range " + pool.Range + ".")
	}

	for leaseName, lease := range configuration.Service.DHCPServer.SharedNetworkName["Local"].Subnet[*matchingSubnet].StaticMapping {
		current = append(current, intermediates.DHCPReservation{
			ProviderState: EdgeOSProviderState{ReservationID: leaseName, Subnet: *matchingSubnet},
			Name:          leaseName,
			Address:       lease.IPAddress,
			MAC:           lease.MACAddress,
		})
	}

	desired := generators.HostsToReservations(hosts, EdgeOSProviderState{Subnet: *matchingSubnet})
	changes := diff.CompareDHCPReservations(current, desired)
	return changes.ToActions(), nil
}

func (provider *EdgeOSProvider) ApplyAction(ctx context.Context, action actions.DHCPAction) error {
	client := provider.client()
	providerState := action.GetProviderState().(EdgeOSProviderState)

	switch typedAction := action.(type) {
	case *actions.DHCPCreateReservationAction:
		err := client.Set(edgeosclient.ConfigurationRoot{
			Service: edgeosclient.Service{
				DHCPServer: edgeosclient.DHCPServer{
					SharedNetworkName: map[string]edgeosclient.SharedNetwork{
						"Local": {
							Subnet: map[string]edgeosclient.Subnet{
								providerState.Subnet: {
									StaticMapping: map[string]*edgeosclient.StaticMapping{
										typedAction.GetName(): {
											IPAddress:  typedAction.Address,
											MACAddress: typedAction.MAC,
										},
									},
								},
							},
						},
					},
				},
			},
		})
		return errors.Wrapf(err, "Failed to create DHCP reservation %s.", typedAction.GetName())
	case *actions.DHCPDeleteReservationAction:
		err := client.Delete(edgeosclient.ConfigurationRoot{
			Service: edgeosclient.Service{
				DHCPServer: edgeosclient.DHCPServer{
					SharedNetworkName: map[string]edgeosclient.SharedNetwork{
						"Local": {
							Subnet: map[string]edgeosclient.Subnet{
								providerState.Subnet: {
									StaticMapping: map[string]*edgeosclient.StaticMapping{
										typedAction.GetName(): nil,
									},
								},
							},
						},
					},
				},
			},
		})
		return errors.Wrapf(err, "Failed to delete DHCP reservation %s.", typedAction.GetName())
	case *actions.DHCPUpdateReservationAction:
		err := client.Set(edgeosclient.ConfigurationRoot{
			Service: edgeosclient.Service{
				DHCPServer: edgeosclient.DHCPServer{
					SharedNetworkName: map[string]edgeosclient.SharedNetwork{
						"Local": {
							Subnet: map[string]edgeosclient.Subnet{
								providerState.Subnet: {
									StaticMapping: map[string]*edgeosclient.StaticMapping{
										typedAction.GetName(): {
											IPAddress:  typedAction.NewAddress,
											MACAddress: typedAction.NewMAC,
										},
									},
								},
							},
						},
					},
				},
			},
		})
		return errors.Wrapf(err, "Failed to update DHCP reservation %s.", typedAction.GetName())
	default:
		return errors.Errorf("Unknown action type %T.", action)
	}
}
