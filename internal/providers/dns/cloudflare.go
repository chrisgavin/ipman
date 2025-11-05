package dns

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
	"github.com/chrisgavin/ipman/internal/diff"
	"github.com/chrisgavin/ipman/internal/generators"
	"github.com/chrisgavin/ipman/internal/intermediates"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
	"github.com/pkg/errors"
)

type CloudflareProvider struct {
	Type     string
	Name     string `yaml:"-"`
	APIKey   string `yaml:"api_key"`
	APIEmail string `yaml:"api_email"`
}

type CloudflareProviderState struct {
	ZoneID   string
	RecordID string
}

func (provider *CloudflareProvider) apiClient() (*cloudflare.Client, error) {
	apiKey := provider.APIKey
	apiEmail := provider.APIEmail
	api := cloudflare.NewClient(
		option.WithAPIKey(apiKey),
		option.WithAPIEmail(apiEmail),
	)
	return api, nil
}

func (provider *CloudflareProvider) GetName(ctx context.Context) string {
	return provider.Name
}

func (provider *CloudflareProvider) GetActions(ctx context.Context, network types.Network, site types.Site, pool types.Pool, hosts []types.Host) ([]actions.DNSAction, error) {
	api, err := provider.apiClient()
	if err != nil {
		return nil, err
	}

	zoneList, err := api.Zones.List(ctx, zones.ZoneListParams{
		Name: cloudflare.F(network.Name),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find zone for network %s.", network.Name)
	}
	if len(zoneList.Result) == 0 {
		return nil, errors.Errorf("No zone found for network %s.", network.Name)
	}
	zoneID := zoneList.Result[0].ID

	records, err := api.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find DNS records for network %s.", network.Name)
	}

	current := []intermediates.DNSRecord{}
	for _, record := range records.Result {
		current = append(current, intermediates.DNSRecord{
			ProviderState: CloudflareProviderState{
				ZoneID:   zoneID,
				RecordID: record.ID,
			},
			Name: record.Name,
			Type: string(record.Type),
			Data: record.Content,
		})
	}

	// Filter out records that are not in the site.
	current = generators.RecordsForSite(network, site, current)

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
		params := dns.RecordNewParams{
			ZoneID: cloudflare.F(providerState.ZoneID),
			Body: dns.RecordNewParamsBody{
				Type:    cloudflare.F(dns.RecordNewParamsBodyType(typedAction.Type)),
				Name:    cloudflare.F(typedAction.Name),
				Content: cloudflare.F(typedAction.Data),
			}}
		_, err = api.DNS.Records.New(ctx, params)
		return errors.Wrapf(err, "Failed to create DNS record %s.", typedAction.Name)
	case *actions.DNSUpdateRecordAction:
		params := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(providerState.ZoneID),
			Body: dns.RecordUpdateParamsBody{
				Type:    cloudflare.F(dns.RecordUpdateParamsBodyType(typedAction.Type)),
				Name:    cloudflare.F(typedAction.Name),
				Content: cloudflare.F(typedAction.NewData),
			},
		}
		_, err = api.DNS.Records.Update(ctx, providerState.RecordID, params)
		return errors.Wrapf(err, "Failed to update DNS record %s.", typedAction.Name)
	case *actions.DNSDeleteRecordAction:
		params := dns.RecordDeleteParams{
			ZoneID: cloudflare.F(providerState.ZoneID),
		}
		_, err := api.DNS.Records.Delete(ctx, providerState.RecordID, params)
		return errors.Wrapf(err, "Failed to delete DNS record %s.", typedAction.Name)
	default:
		return errors.Errorf("Unknown action type %T.", action)
	}
}
