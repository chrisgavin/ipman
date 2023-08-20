package diff

import (
	"sort"

	"github.com/chrisgavin/ipman/internal/intermediates"
)

func CompareDNSRecords(current []intermediates.DNSRecord, desired []intermediates.DNSRecord) intermediates.DNSChanges {
	changes := intermediates.DNSChanges{
		Additions: []intermediates.DNSRecord{},
		Deletions: []intermediates.DNSRecord{},
		Updates:   map[intermediates.DNSRecord]intermediates.DNSRecord{},
	}

	currentIdentifiersToRecords := map[string]intermediates.DNSRecords{}
	for _, record := range current {
		currentIdentifiersToRecords[record.Identifier()] = append(currentIdentifiersToRecords[record.Identifier()], record)
	}

	desiredIdentifiersToRecords := map[string]intermediates.DNSRecords{}
	for _, record := range desired {
		desiredIdentifiersToRecords[record.Identifier()] = append(desiredIdentifiersToRecords[record.Identifier()], record)
	}

	for identifier, currentRecords := range currentIdentifiersToRecords {
		desiredRecords, exists := desiredIdentifiersToRecords[identifier]
		if !exists {
			changes.Deletions = append(changes.Deletions, currentRecords...)
		} else {
			sort.Sort(intermediates.DNSRecords(currentRecords))
			sort.Sort(intermediates.DNSRecords(desiredRecords))
			maxLength := len(currentRecords)
			if len(desiredRecords) > maxLength {
				maxLength = len(desiredRecords)
			}

			for i := 0; i < maxLength; i++ {
				if i >= len(currentRecords) {
					changes.Additions = append(changes.Additions, desiredRecords[i])
				} else if i >= len(desiredRecords) {
					changes.Deletions = append(changes.Deletions, currentRecords[i])
				} else if !currentRecords[i].Equals(&desiredRecords[i]) {
					changes.Updates[currentRecords[i]] = desiredRecords[i]
				}
			}
		}
	}

	for identifier, desiredRecords := range desiredIdentifiersToRecords {
		_, exists := currentIdentifiersToRecords[identifier]
		if !exists {
			changes.Additions = append(changes.Additions, desiredRecords...)
		}
	}

	return changes
}
