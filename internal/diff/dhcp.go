package diff

import "github.com/chrisgavin/ipman/internal/intermediates"

func CompareDHCPReservations(current []intermediates.DHCPReservation, desired []intermediates.DHCPReservation) intermediates.DHCPChanges {
	// TODO: Handle duplicate stuff better
	changes := intermediates.DHCPChanges{
		Additions: []intermediates.DHCPReservation{},
		Deletions: []intermediates.DHCPReservation{},
		Updates:   map[intermediates.DHCPReservation]intermediates.DHCPReservation{},
	}

	currentIdentifiersToReservation := map[string]intermediates.DHCPReservation{}
	for _, reservation := range current {
		currentIdentifiersToReservation[reservation.Identifier()] = reservation
	}

	desiredIdentifiersToReservation := map[string]intermediates.DHCPReservation{}
	for _, reservation := range desired {
		desiredIdentifiersToReservation[reservation.Identifier()] = reservation
	}

	for identifier, currentReservation := range currentIdentifiersToReservation {
		desiredReservation, exists := desiredIdentifiersToReservation[identifier]
		if !exists {
			changes.Deletions = append(changes.Deletions, currentReservation)
		} else {
			if !currentReservation.Equals(&desiredReservation) {
				changes.Updates[currentReservation] = desiredReservation
			}
		}
	}

	for identifier, desiredReservation := range desiredIdentifiersToReservation {
		_, exists := currentIdentifiersToReservation[identifier]
		if !exists {
			changes.Additions = append(changes.Additions, desiredReservation)
		}
	}

	return changes
}
