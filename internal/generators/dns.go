package generators

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
)

func HostsToRecords(network types.Network, site types.Site, pool types.Pool, hosts []types.Host) []intermediates.DNSRecord {
	result := []intermediates.DNSRecord{}
	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
		result = append(result, intermediates.DNSRecord{
			Name: fullName,
			Type: "A",
			Data: primaryInterface.Address,
		})
		for _, networkInterface := range host.Interfaces {
			if networkInterface.Address == "" {
				continue
			}
			fullName := fmt.Sprintf("%s.%s.%s.%s", networkInterface.Name, host.Name, site.Name, network.Name)
			result = append(result, intermediates.DNSRecord{
				Name: fullName,
				Type: "A",
				Data: networkInterface.Address,
			})
		}
	}
	return result
}
