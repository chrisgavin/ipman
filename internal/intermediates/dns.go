package intermediates

import "github.com/chrisgavin/ipman/internal/actions"

type DNSRecord struct {
	Name string
	Type string
	Data string
}

func (record *DNSRecord) Identifier() string {
	return record.Name + ":" + record.Type
}

func (record *DNSRecord) Equals(other *DNSRecord) bool {
	return record.Data == other.Data
}

type DNSRecords []DNSRecord

func (records DNSRecords) Len() int {
	return len(records)
}

func (records DNSRecords) Less(i, j int) bool {
	return records[i].Identifier() < records[j].Identifier()
}

func (records DNSRecords) Swap(i, j int) {
	records[i], records[j] = records[j], records[i]
}

type DNSChanges struct {
	Additions []DNSRecord
	Deletions []DNSRecord
	Updates   map[DNSRecord]DNSRecord
}

func (changes *DNSChanges) ToActions() []actions.DNSAction {
	result := []actions.DNSAction{}
	for _, deletion := range changes.Deletions {
		result = append(result, &actions.DNSDeleteRecordAction{
			BaseDNSAction: actions.BaseDNSAction{
				Name: deletion.Name,
				Type: deletion.Type,
			},
		})
	}
	for current, desired := range changes.Updates {
		result = append(result, &actions.DNSUpdateRecordAction{
			BaseDNSAction: actions.BaseDNSAction{
				Name: current.Name,
				Type: current.Type,
			},
			OldData: current.Data,
			NewData: desired.Data,
		})
	}
	for _, addition := range changes.Additions {
		result = append(result, &actions.DNSCreateRecordAction{
			BaseDNSAction: actions.BaseDNSAction{
				Name: addition.Name,
				Type: addition.Type,
			},
			Data: addition.Data,
		})
	}
	return result
}
