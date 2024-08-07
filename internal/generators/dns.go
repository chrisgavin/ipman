package generators

import (
	"fmt"
	"strings"

	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
)

func HostsToRecords(hosts []types.Host, providerState interface{}) []intermediates.DNSRecord {
	result := []intermediates.DNSRecord{}
	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, host.Pool.Site.Name, host.Pool.Site.Network.Name)
		result = append(result, intermediates.DNSRecord{
			Name:          fullName,
			Type:          "A",
			Data:          primaryInterface.Address,
			ProviderState: providerState,
		})
		for _, networkInterface := range host.Interfaces {
			if networkInterface.Address == "" {
				continue
			}
			fullName := fmt.Sprintf("%s.%s.%s.%s", networkInterface.Name, host.Name, host.Pool.Site.Name, host.Pool.Site.Network.Name)
			result = append(result, intermediates.DNSRecord{
				Name:          fullName,
				Type:          "A",
				Data:          networkInterface.Address,
				ProviderState: providerState,
			})
		}
		for _, record := range host.Records {
			fullName := fmt.Sprintf("%s.%s.%s.%s", record.Name, host.Name, host.Pool.Site.Name, host.Pool.Site.Network.Name)
			result = append(result, intermediates.DNSRecord{
				Name:          fullName,
				Type:          record.Type,
				Data:          record.Data,
				ProviderState: providerState,
			})
		}
	}
	return result
}

func RecordsForSite(network types.Network, site types.Site, records []intermediates.DNSRecord) []intermediates.DNSRecord {
	result := []intermediates.DNSRecord{}
	for _, record := range records {
		if !strings.HasSuffix(record.Name, fmt.Sprintf(".%s.%s", site.Name, network.Name)) {
			continue
		}
		result = append(result, record)
	}
	return result
}
