// Package libdnstemplate implements a DNS record management client compatible
// with the libdns interfaces for LuaDNS.
package luadns

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/libdns/libdns"
	"github.com/luadns/luadns-go"
)

// Provider facilitates DNS record manipulation with LuaDNS.
type Provider struct {
	Email  string `json:"email,omitempty"`
	APIKey string `json:"api_key,omitempty"`
	mutex  sync.Mutex
}

// ListZones list available DNS zones.
func (p *Provider) ListZones(ctx context.Context) ([]libdns.Zone, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	zones, err := p.client().ListZones(ctx, &luadns.ListParams{})
	if err != nil {
		return nil, libdns.AtomicErr(err)
	}

	return toLibdnsZone(zones)
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	z, err := p.getZone(ctx, zone)
	if err != nil {
		return nil, err
	}

	records, err := p.client().ListRecords(ctx, z, &luadns.ListParams{})
	if err != nil {
		return nil, libdns.AtomicErr(err)
	}

	return toLibdnsRecord(records, zone)
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	z, err := p.getZone(ctx, zone)
	if err != nil {
		return nil, err
	}

	records, err := p.client().CreateManyRecords(ctx, z, toLuaRR(recs, zone))
	if err != nil {
		return nil, libdns.AtomicErr(err)
	}

	return toLibdnsRecord(records, zone)
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	z, err := p.getZone(ctx, zone)
	if err != nil {
		return nil, err
	}

	records, err := p.client().UpdateManyRecords(ctx, z, toLuaRR(recs, zone))
	if err != nil {
		return nil, libdns.AtomicErr(err)

	}

	return toLibdnsRecord(records, zone)
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	z, err := p.getZone(ctx, zone)
	if err != nil {
		return nil, err
	}

	records, err := p.client().DeleteManyRecords(ctx, z, toLuaRR(recs, zone))
	if err != nil {
		return nil, libdns.AtomicErr(err)
	}

	return toLibdnsRecord(records, zone)
}

// getZone issues a search request to find zone by name.
func (p *Provider) getZone(ctx context.Context, name string) (*luadns.Zone, error) {
	query := unFQDN(name)

	result, err := p.client().ListZones(ctx, &luadns.ListParams{Query: query})
	if err != nil {
		return nil, err
	}

	for _, z := range result {
		if strings.EqualFold(z.Name, query) {
			return z, nil
		}
	}

	return nil, fmt.Errorf("zone %v not found", name)
}

func (p *Provider) client() *luadns.Client {
	return luadns.NewClient(p.Email, p.APIKey)
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
	_ libdns.ZoneLister     = (*Provider)(nil)
)
