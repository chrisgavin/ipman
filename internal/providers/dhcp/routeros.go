package dhcp

import (
	"context"
	"fmt"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/go-routeros/routeros"
	"github.com/pkg/errors"
)

type RouterOSProvider struct {
	Type     string
	Address  string
	Username string
	Password string
	Insecure bool
}

func (provider *RouterOSProvider) client() (*routeros.Client, error) {
	var client *routeros.Client
	var err error
	username := provider.Username
	password := provider.Password
	if provider.Insecure {
		client, err = routeros.Dial(provider.Address, username, password)
	} else {
		client, err = routeros.DialTLS(provider.Address, username, password, nil)
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to RouterOS.")
	}
	return client, err
}

func (provider *RouterOSProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DHCPAction, error) {
	result := []actions.DHCPAction{}

	client, err := provider.client()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	leases, err := client.Run("/ip/dhcp-server/lease/print")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get leases.")
	}

	staticLeases := []map[string]string{}
	for _, lease := range leases.Re {
		if lease.Map["dynamic"] == "false" {
			staticLeases = append(staticLeases, lease.Map)
		}
	}

	for _, lease := range staticLeases {
		leaseFound := false
		for _, host := range hosts {
			primaryInterface := host.Interfaces[0]
			fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
			if lease["host-name"] == fullName && lease["mac-address"] == primaryInterface.MAC && lease["address"] == primaryInterface.Address {
				leaseFound = true
				break
			}
		}
		if !leaseFound {
			result = append(result, &actions.DHCPDeleteReservationAction{
				BaseDHCPAction: actions.BaseDHCPAction{
					Name: lease["host-name"],
				},
				MAC:     lease["mac-address"],
				Address: lease["address"],
			})
		}
	}

	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
		leaseFound := false
		for _, lease := range staticLeases {
			if lease["host-name"] == fullName && lease["mac-address"] == primaryInterface.MAC && lease["address"] == primaryInterface.Address {
				leaseFound = true
				break
			}
		}
		if leaseFound {
			continue
		}
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
