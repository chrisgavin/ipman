package generators

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
)

func HostsToReservations(hosts []types.Host, providerState interface{}) []intermediates.DHCPReservation {
	result := []intermediates.DHCPReservation{}
	for _, host := range hosts {
		for _, networkInterface := range host.Interfaces {
			if networkInterface.MAC == "" {
				continue
			}
			fullName := fmt.Sprintf("%s.%s.%s.%s", networkInterface.Name, host.Name, host.Pool.Site.Name, host.Pool.Site.Network.Name)
			result = append(result, intermediates.DHCPReservation{
				Name:          fullName,
				MAC:           networkInterface.MAC,
				Address:       networkInterface.Address,
				ProviderState: providerState,
			})
		}
	}
	return result
}
