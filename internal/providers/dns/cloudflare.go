package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
)

type CloudflareProvider struct {
	Type     string
	APIKey   string `yaml:"api_key"`
	APIEmail string `yaml:"api_email"`
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
	result := []actions.DNSAction{}
	for _, host := range hosts {
		primaryInterface := host.Interfaces[0]
		fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
		found := false
		for _, record := range records {
			if record.Name == fullName && record.Type == "A" && record.Content == primaryInterface.Address {
				found = true
				continue
			}
			if record.Name == fullName && record.Type == "A" {
				found = true
				result = append(result, &actions.DNSChangeRecordAction{
					BaseDNSAction: actions.BaseDNSAction{
						Name: fullName,
						Type: "A",
					},
					OldData: record.Content,
					NewData: primaryInterface.Address,
				})
				continue
			}
			if record.Name == fullName && record.Type != "A" {
				result = append(result, &actions.DNSDeleteRecordAction{
					BaseDNSAction: actions.BaseDNSAction{
						Name: fullName,
						Type: record.Type,
					},
				})
				continue
			}
		}
		if !found {
			result = append(result, &actions.DNSCreateRecordAction{
				BaseDNSAction: actions.BaseDNSAction{
					Name: fullName,
					Type: "A",
				},
				Data: primaryInterface.Address,
			})
		}
	}
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
		found := false
		for _, host := range hosts {
			fullName := fmt.Sprintf("%s.%s.%s", host.Name, site.Name, network.Name)
			if record.Name == fullName {
				found = true
				continue
			}
		}
		if !found {
			result = append(result, &actions.DNSDeleteRecordAction{
				BaseDNSAction: actions.BaseDNSAction{
					Name: record.Name,
					Type: record.Type,
				},
			})
		}
	}
	return result, nil
}
