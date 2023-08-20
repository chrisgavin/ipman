package generators

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
)

func HostsToReservations(network types.Network, site types.Site, pool types.Pool, hosts []types.Host) []intermediates.DHCPReservation {
	result := []intermediates.DHCPReservation{}
	for _, host := range hosts {
		for _, networkInterface := range host.Interfaces {
			if networkInterface.MAC == "" {
				continue
			}
			fullName := fmt.Sprintf("%s.%s.%s.%s", networkInterface.Name, host.Name, site.Name, network.Name)
			result = append(result, intermediates.DHCPReservation{
				Name:    fullName,
				MAC:     networkInterface.MAC,
				Address: networkInterface.Address,
			})
		}
	}
	return result
}
