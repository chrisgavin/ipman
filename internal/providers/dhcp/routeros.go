package dhcp

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/go-routeros/routeros"
	"github.com/pkg/errors"
)

type RouterOSProvider struct {
	Type     string
	Name     string `yaml:"-"`
	Address  string
	Username string
	Password string
	Insecure bool
}

type RouterOSProviderState struct {
	ReservationID string
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

func (provider *RouterOSProvider) GetName(ctx context.Context) string {
	return provider.Name
}

func (provider *RouterOSProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DHCPAction, error) {
	current := []intermediates.DHCPReservation{}

	client, err := provider.client()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	leases, err := client.Run("/ip/dhcp-server/lease/print")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get leases.")
	}

	for _, lease := range leases.Re {
		if lease.Map["dynamic"] == "false" {
			current = append(current, intermediates.DHCPReservation{
				ProviderState: RouterOSProviderState{ReservationID: lease.Map[".id"]},
				Name:          lease.Map["comment"],
				MAC:           lease.Map["mac-address"],
				Address:       lease.Map["address"],
				Disabled:      lease.Map["disabled"] == "true",
			})
		}
	}

	desired := generators.HostsToReservations(hosts, nil)
	changes := diff.CompareDHCPReservations(current, desired)
	return changes.ToActions(), nil
}

func (provider *RouterOSProvider) ApplyAction(ctx context.Context, action actions.DHCPAction) error {
	client, err := provider.client()
	if err != nil {
		return err
	}
	defer client.Close()

	switch typedAction := action.(type) {
	case *actions.DHCPCreateReservationAction:
		_, err := client.Run("/ip/dhcp-server/lease/add", "=address="+typedAction.Address, "=mac-address="+typedAction.MAC, "=comment="+typedAction.GetName())
		return errors.Wrapf(err, "Failed to create DHCP reservation %s.", typedAction.GetName())
	case *actions.DHCPDeleteReservationAction:
		providerState := action.GetProviderState().(RouterOSProviderState)
		_, err := client.Run("/ip/dhcp-server/lease/remove", "=.id="+providerState.ReservationID)
		return errors.Wrapf(err, "Failed to delete DHCP reservation %s.", typedAction.GetName())
	case *actions.DHCPUpdateReservationAction:
		providerState := action.GetProviderState().(RouterOSProviderState)
		_, err := client.Run("/ip/dhcp-server/lease/set", "=.id="+providerState.ReservationID, "=address="+typedAction.NewAddress, "=mac-address="+typedAction.NewMAC)
		return errors.Wrapf(err, "Failed to update DHCP reservation %s.", typedAction.GetName())
	default:
		return errors.Errorf("Unknown action type %T.", action)
	}
}
