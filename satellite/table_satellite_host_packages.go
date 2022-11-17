package satellite

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
				Description: "The id of the host having the package.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the host.",
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
		},
		List: &plugin.ListConfig{
			Hydrate: listSatelliteHostPackage,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "id",
					Require: plugin.AnyOf,
				},
				&plugin.KeyColumn{
					Name:    "name",
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
	value := d.KeyColumnQuals["id"].GetInt64Value()
	if value != 0 {
		id = fmt.Sprintf("%d", value)
	} else {
		id = d.KeyColumnQuals["name"].GetStringValue()
	}
	if id == "" {
		plugin.Logger(ctx).Error("no valid host id provided")
		return nil, errors.New("no valid host id or name provided")
	}

	request := client.
		R().
		SetContext(ctx).
		SetPathParam("id", id)

	request.SetHeaders(map[string]string{
		"Accept-Encoding": "gzip",
		"Accept":          "text/html",
	})
	result := &struct {
		Total    int         `json:"total"`
		Subtotal int         `json:"subtotal"`
		Page     int         `json:"page"`
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
	_, err = request.Get("/hosts/{id}/packages")
	if err != nil {
		plugin.Logger(ctx).Error("error performing request", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Error("request successful", "total", result.Total, "subtotal", result.Subtotal, "page", result.Page, "per page", result.PerPage) //, "response", toJSON(response.Body()))

	for _, pkg := range result.Packages {
		pkg := pkg
		plugin.Logger(ctx).Error("package", "contents", toJSON(pkg))
		d.StreamListItem(ctx, pkg)
	}

	return nil, nil
}

type apiPackage struct {
	ID    int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	NVREA string `json:"nvrea,omitempty" yaml:"nvrea,omitempty"`
	NVRA  string `json:"nvra,omitempty" yaml:"nvra,omitempty"`
}
