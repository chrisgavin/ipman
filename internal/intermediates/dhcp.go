package intermediates

import (
	"github.com/chrisgavin/ipman/internal/actions"
)

type DHCPReservation struct {
	ProviderState interface{}
	Name          string
	MAC           string
	Address       string
	Disabled      bool
}

func (reservation *DHCPReservation) Identifier() string {
	if reservation.Name != "" {
		return reservation.Name
	}
	return reservation.MAC
}

func (reservation *DHCPReservation) Equals(other *DHCPReservation) bool {
	return reservation.MAC == other.MAC && reservation.Address == other.Address && reservation.Disabled == other.Disabled
}

type DHCPChanges struct {
	Additions []DHCPReservation
	Deletions []DHCPReservation
	Updates   map[DHCPReservation]DHCPReservation
}

func (changes *DHCPChanges) ToActions() []actions.DHCPAction {
	result := []actions.DHCPAction{}
	for _, deletion := range changes.Deletions {
		result = append(result, &actions.DHCPDeleteReservationAction{
			BaseDHCPAction: actions.BaseDHCPAction{
				Name:          deletion.Name,
				ProviderState: deletion.ProviderState,
			},
		})
	}
	for current, desired := range changes.Updates {
		result = append(result, &actions.DHCPUpdateReservationAction{
			BaseDHCPAction: actions.BaseDHCPAction{
				Name:          current.Name,
				ProviderState: current.ProviderState,
			},
			OldMAC:     current.MAC,
			NewMAC:     desired.MAC,
			OldAddress: current.Address,
			NewAddress: desired.Address,
		})
	}
	for _, addition := range changes.Additions {
		result = append(result, &actions.DHCPCreateReservationAction{
			BaseDHCPAction: actions.BaseDHCPAction{
				Name: addition.Name,
			},
			MAC:     addition.MAC,
			Address: addition.Address,
		})
	}
	return result
}
