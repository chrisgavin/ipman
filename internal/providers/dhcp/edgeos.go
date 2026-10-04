package dhcp

import (
	"context"
	"sort"

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

func (provider *EdgeOSProvider) GetActions(ctx context.Context, network types.Network, site types.Site) ([]actions.DHCPAction, error) {
	client := provider.client()

	configuration, err := client.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get DHCP leases.")
	}

	subnets := configuration.Service.DHCPServer.SharedNetworkName["Local"].Subnet

	subnetHosts := map[string][]types.Host{}
	for _, pool := range site.Pools {
		poolRange := ipaddr.NewIPAddressString(pool.Range)
		matchingSubnet := ""
		for subnetRange := range subnets {
			parsedSubnetRange := ipaddr.NewIPAddressString(subnetRange)
			if parsedSubnetRange.Contains(poolRange) {
				if matchingSubnet != "" {
					return nil, errors.New("Multiple subnets contain the pool range " + pool.Range + ".")
				}
				matchingSubnet = subnetRange
			}
		}
		if matchingSubnet == "" {
			return nil, errors.New("No subnet contains the pool range " + pool.Range + ".")
		}
		subnetHosts[matchingSubnet] = append(subnetHosts[matchingSubnet], pool.Hosts...)
	}

	matchingSubnets := []string{}
	for subnet := range subnetHosts {
		matchingSubnets = append(matchingSubnets, subnet)
	}
	sort.Strings(matchingSubnets)

	result := []actions.DHCPAction{}
	for _, subnet := range matchingSubnets {
		current := []intermediates.DHCPReservation{}
		for leaseName, lease := range subnets[subnet].StaticMapping {
			current = append(current, intermediates.DHCPReservation{
				ProviderState: EdgeOSProviderState{ReservationID: leaseName, Subnet: subnet},
				Name:          leaseName,
				Address:       lease.IPAddress,
				MAC:           lease.MACAddress,
			})
		}

		desired := generators.HostsToReservations(subnetHosts[subnet], EdgeOSProviderState{Subnet: subnet})
		changes := diff.CompareDHCPReservations(current, desired)
		result = append(result, changes.ToActions()...)
	}
	return result, nil
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
