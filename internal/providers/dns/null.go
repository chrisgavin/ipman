package dns

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/types"
)

type NullProvider struct {
	Type string
}

func (provider *NullProvider) GetActions(network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DNSAction, error) {
	result := []actions.DNSAction{}
	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
		result = append(result, &actions.DNSCreateRecordAction{
			BaseDNSAction: actions.BaseDNSAction{
				Name: fullName,
				Type: "A",
			},
			Data: primaryInterface.Address,
		})
	}
	return result, nil
}
