package satellite

import (
	"context"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const SatelliteClientKey = "satellite_client"

func getClient(ctx context.Context, d *plugin.QueryData) (*resty.Client, error) {
	plugin.Logger(ctx).Debug("retrieving client")

	// load connection from cache, which preserves throttling protection etc
	if cachedData, ok := d.ConnectionManager.Cache.Get(SatelliteClientKey); ok {
		plugin.Logger(ctx).Debug("returning satellite client from cache")
		return cachedData.(*resty.Client), nil
	}

	client := resty.
		New().
		SetHeader("Accept", "application/json,version=2").
		SetHeader("Content-Type", "application/json")

	satelliteConfig := GetConfig(d.Connection)

	if satelliteConfig.Endpoint != nil {
		client.SetBaseURL(*satelliteConfig.Endpoint)
	} else {
		plugin.Logger(ctx).Error("no API endpoint available")
		return nil, errors.New("no API endpoint available")
	}

	if satelliteConfig.Username != nil && satelliteConfig.Password != nil {
		client.SetBasicAuth(*satelliteConfig.Username, *satelliteConfig.Password)
	} else {
		plugin.Logger(ctx).Error("no authentication info available")
		return nil, errors.New("no authentication info available")
	}

	if satelliteConfig.Organisation != nil {
		client.SetQueryParam("organization_id", *satelliteConfig.Organisation)
	}

	if satelliteConfig.Location != nil {
		client.SetQueryParam("location_id", *satelliteConfig.Location)
	}

	// save to cache
	plugin.Logger(ctx).Debug("saving satellite client to cache")
	d.ConnectionManager.Cache.Set(SatelliteClientKey, client)

	return client, nil
}
