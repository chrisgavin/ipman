package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
)

type CloudflareProvider struct {
	Type     string
	APIKey   string `yaml:"api_key"`
	APIEmail string `yaml:"api_email"`
}

type CloudflareProviderState struct {
	ZoneID   string
	RecordID string
}

func (provider *CloudflareProvider) apiClient() (*cloudflare.API, error) {
	apiKey := provider.APIKey
	apiEmail := provider.APIEmail
	api, err := cloudflare.New(apiKey, apiEmail)
	return api, errors.Wrap(err, "Failed to create Cloudflare API client.")
}

func (provider *CloudflareProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DNSAction, error) {
	api, err := provider.apiClient()
	if err != nil {
		return nil, err
	}
	zoneID, err := api.ZoneIDByName(network.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find zone ID for network %s.", network.Name)
	}
	records, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find DNS records for network %s.", network.Name)
	}

	current := []intermediates.DNSRecord{}
	for _, record := range records {
		recordPartOfSite := false
		for _, site := range network.Sites {
			if strings.HasSuffix(record.Name, fmt.Sprintf(".%s.%s", site.Name, network.Name)) {
				recordPartOfSite = true
				continue
			}
		}
		if !recordPartOfSite {
			continue
		}
		current = append(current, intermediates.DNSRecord{
			ProviderState: CloudflareProviderState{
				ZoneID:   zoneID,
				RecordID: record.ID,
			},
			Name: record.Name,
			Type: record.Type,
			Data: record.Content,
		})
	}

	desired := generators.HostsToRecords(hosts, CloudflareProviderState{ZoneID: zoneID})
	changes := diff.CompareDNSRecords(current, desired)
	return changes.ToActions(), nil
}

func (provider *CloudflareProvider) ApplyAction(ctx context.Context, action actions.DNSAction) error {
	api, err := provider.apiClient()
	if err != nil {
		return err
	}
	providerState := action.GetProviderState().(CloudflareProviderState)
	switch typedAction := action.(type) {
	case *actions.DNSCreateRecordAction:
		record := cloudflare.DNSRecord{
			Type:    typedAction.Type,
			Name:    typedAction.Name,
			Content: typedAction.Data,
		}
		_, err = api.CreateDNSRecord(ctx, providerState.ZoneID, record)
		return errors.Wrapf(err, "Failed to create DNS record %s.", typedAction.Name)
	case *actions.DNSUpdateRecordAction:
		record := cloudflare.DNSRecord{
			Type:    typedAction.Type,
			Name:    typedAction.Name,
			Content: typedAction.NewData,
		}
		err = api.UpdateDNSRecord(ctx, providerState.ZoneID, providerState.RecordID, record)
		return errors.Wrapf(err, "Failed to update DNS record %s.", typedAction.Name)
	case *actions.DNSDeleteRecordAction:
		err := api.DeleteDNSRecord(ctx, providerState.ZoneID, providerState.RecordID)
		return errors.Wrapf(err, "Failed to delete DNS record %s.", typedAction.Name)
	default:
		return errors.Errorf("Unknown action type %T.", action)
	}
}
