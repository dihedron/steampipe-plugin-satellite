package satellite

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableSatelliteHostPackage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "satellite_host_package",
		Description: "Red Hat Satellite Host Packages",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "The id of the package.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the package.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "nvrea",
				Type:        proto.ColumnType_STRING,
				Description: "The Name, Version, Release, Environment and Architecture (NVREA) of the package.",
				Transform:   transform.FromField("NVREA"),
			},
			{
				Name:        "nvra",
				Type:        proto.ColumnType_STRING,
				Description: "The Name, Version, Release and Architecture (NVRA) of the package.",
				Transform:   transform.FromField("NVRA"),
			},
			{
				Name:        "host_id",
				Type:        proto.ColumnType_INT,
				Description: "The id of the host having the package.",
				Transform:   transform.FromField("HostID"),
			},
			{
				Name:        "host_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the host having the package.",
				Transform:   transform.FromField("HostName"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSatelliteHostPackage,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "host_id",
					Require: plugin.AnyOf,
				},
				&plugin.KeyColumn{
					Name:    "host_name",
					Require: plugin.AnyOf,
				},
			},
		},
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.KeyColumnSlice{
		// 		&plugin.KeyColumn{
		// 			Name:    "id",
		// 			Require: plugin.Optional,
		// 		},
		// 		&plugin.KeyColumn{
		// 			Name:    "name",
		// 			Require: plugin.Optional,
		// 		},
		// 	},
		// 	// plugin.SingleColumn("id"),
		// 	Hydrate: getSatelliteHostPackage,
		// },
	}
}

//// LIST FUNCTIONS

func listSatelliteHostPackage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	setLogLevel(ctx, d)
	plugin.Logger(ctx).Debug("retrieving satellite package list for host", "query data", toJSON(d))

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving satellite client", "error", err)
		return nil, err
	}

	id := ""
	hostid := d.KeyColumnQuals["host_id"].GetInt64Value()
	hostname := d.KeyColumnQuals["host_name"].GetStringValue()
	if hostid != 0 {
		id = fmt.Sprintf("%d", hostid)
	} else {
		id = hostname
	}
	if id == "" {
		plugin.Logger(ctx).Error("no valid host id or name provided")
		return nil, errors.New("no valid host id or name provided")
	}

	page := 1
	for {

		request := client.
			R().
			SetContext(ctx).
			SetPathParam("id", id).
			SetQueryParam("page", fmt.Sprintf("%d", page))

		request.SetHeaders(map[string]string{
			"Accept-Encoding": "gzip",
			"Accept":          "text/html",
		})
		result := &struct {
			Total    int         `json:"total"`
			Subtotal int         `json:"subtotal"`
			Page     interface{} `json:"page"`
			PerPage  int         `json:"per_page"`
			Error    interface{} `json:"error"`
			Search   interface{} `json:"search"`
			Sort     struct {
				By    string `json:"by"`
				Order string `json:"order"`
			} `json:"sort"`
			Packages []apiPackage `json:"results"`
		}{}
		request.SetResult(result)
		// response, err := request.Get(fmt.Sprintf("/hosts/{id}/packages?page=%d", page))
		response, err := request.Get("/hosts/{id}/packages")
		if err != nil {
			plugin.Logger(ctx).Error("error performing request", "error", err, "response", toPrettyJSON(response.Body()))
			return nil, err
		}
		plugin.Logger(ctx).Debug("request successful", "total", result.Total, "subtotal", result.Subtotal, "page", result.Page, "per page", result.PerPage, "response", toJSON(response.Body()))

		for _, pkg := range result.Packages {
			pkg := pkg
			//plugin.Logger(ctx).Debug("package", "contents", toJSON(pkg))
			d.StreamListItem(ctx, &struct {
				HostID   int    `json:"host_id,omitempty" yaml:"host_id,omitempty"`
				HostName string `json:"host_name,omitempty" yaml:"host_name,omitempty"`
				apiPackage
			}{
				HostID:     int(hostid),
				HostName:   hostname,
				apiPackage: pkg,
			})
		}

		// handle pagination. Note that the Satellite API returns results.Page as an
		// integer if there is no page?{page} query  parameter, and as a string if you
		// set one; thus we need to handle both cases
		resultPage := 0
		switch v := result.Page.(type) {
		case int:
			resultPage = v
		case int32:
			resultPage = int(v)
		case int64:
			resultPage = int(v)
		case float32:
			resultPage = int(v)
		case float64:
			resultPage = int(v)
		case string:
			resultPage, _ = strconv.Atoi(v)
		default:
			plugin.Logger(ctx).Debug("unsupported type in pagination", "type", fmt.Sprintf("%T", result.Page))
			return nil, fmt.Errorf("unexpected type in pagination API result: %T", result.Page)
		}
		if result.PerPage*resultPage < result.Total {
			page++
			plugin.Logger(ctx).Debug("retrieving next page", "page", page)
		} else {
			plugin.Logger(ctx).Debug("all pages retrieved", "subtotal", result.Subtotal, "total", result.Total)
			break
		}
	}

	return nil, nil
}

type apiPackage struct {
	ID    int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	NVREA string `json:"nvrea,omitempty" yaml:"nvrea,omitempty"`
	NVRA  string `json:"nvra,omitempty" yaml:"nvra,omitempty"`
}
