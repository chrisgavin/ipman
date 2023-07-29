package dhcp

import (
	"context"
	"fmt"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/types"
)

type NullProvider struct {
	Type string
}

func (provider *NullProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DHCPAction, error) {
	result := []actions.DHCPAction{}
	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
		result = append(result, &actions.DHCPCreateReservationAction{
			BaseDHCPAction: actions.BaseDHCPAction{
				Name: fullName,
			},
			MAC:     primaryInterface.MAC,
			Address: primaryInterface.Address,
		})
	}
	return result, nil
}
