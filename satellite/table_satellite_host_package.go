package satellite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the package.",
				Transform: transform.FromField("NVRA").Transform(func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					_, ver, _, _, err := ParseNVRA(d.Value.(string))
					return ver, err
				}),
			},
			{
				Name:        "release",
				Type:        proto.ColumnType_STRING,
				Description: "The release of the package.",
				Transform: transform.FromField("NVRA").Transform(func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					_, _, rel, _, err := ParseNVRA(d.Value.(string))
					return rel, err
				}),
			},
			{
				Name:        "architecture",
				Type:        proto.ColumnType_STRING,
				Description: "The architecture of the package.",
				Transform: transform.FromField("NVRA").Transform(func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					_, _, _, arch, err := ParseNVRA(d.Value.(string))
					return arch, err
				}),
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
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "host_name",
					Require: plugin.Optional,
				},
			},
		},
	}
}

//// LIST FUNCTIONS

func listSatelliteHostPackage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	setLogLevel(ctx, d)
	plugin.Logger(ctx).Debug("retrieving satellite package list for host", "query data", utils.ToJSON(d))

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving satellite client", "error", err)
		return nil, err
	}

	hosts := []apiHost{}

	// first, check if running against a single host
	single := apiHost{
		ID:   int(d.EqualsQuals["host_id"].GetInt64Value()),
		Name: d.EqualsQuals["host_name"].GetStringValue(),
	}
	if single.ID != 0 {
		plugin.Logger(ctx).Debug("running query against single host", "host", utils.ToPrettyJSON(single))
		hosts = append(hosts, single)
	} else if single.Name != "" {
		id, err := resolveHostID(ctx, d, single.Name)
		if err != nil {
			plugin.Logger(ctx).Error("error resolving host by name", "name", single.Name)
			return nil, err
		}
		single.ID = id
		hosts = append(hosts, single)
	} else {
		hosts, err = listSatelliteHostImpl(ctx, client, true)
		if err != nil {
			plugin.Logger(ctx).Error("error getting list of hosts", "error", err)
			return nil, err
		}
	}

	plugin.Logger(ctx).Debug("retrieving packages from hosts", "hosts", utils.ToPrettyJSON(hosts))

outer:
	for _, host := range hosts {
		host := host

		id := fmt.Sprintf("%d", host.ID)

		plugin.Logger(ctx).Debug("running query against host", "id", id, "name", host.Name)

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
				Packages []apiHostPackage `json:"results"`
			}{}
			request.SetResult(result)
			plugin.Logger(ctx).Debug("running request", "id", id)
			// response, err := request.Get(fmt.Sprintf("/hosts/{id}/packages?page=%d", page))
			response, err := request.Get("/api/hosts/{id}/packages")

			if err != nil || response.IsError() {
				plugin.Logger(ctx).Error("error performing request", "status", response.Status(), "error", err, "response", utils.ToPrettyJSON(response.Body()))
				return nil, fmt.Errorf("request %q error %d (%s): %w", request.URL, response.StatusCode(), response.Status(), err)
			}
			plugin.Logger(ctx).Debug("request successful", "status", response.Status(), "total", result.Total, "subtotal", result.Subtotal, "page", result.Page, "per page", result.PerPage, "response", utils.ToJSON(response.Body()))

			for _, pkg := range result.Packages {
				if ctx.Err() != nil {
					plugin.Logger(ctx).Debug("context done, exit")
					break outer
				}
				pkg := pkg
				plugin.Logger(ctx).Debug("package", "contents", utils.ToPrettyJSON(pkg))
				d.StreamListItem(ctx, &struct {
					HostID   int    `json:"host_id,omitempty" yaml:"host_id,omitempty"`
					HostName string `json:"host_name,omitempty" yaml:"host_name,omitempty"`
					apiHostPackage
				}{
					HostID:         host.ID,
					HostName:       host.Name,
					apiHostPackage: pkg,
				})
			}

			// handle pagination; note that the Satellite API returns results.Page as an
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
	}

	return nil, nil
}

type apiHostPackage struct {
	ID    int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	NVREA string `json:"nvrea,omitempty" yaml:"nvrea,omitempty"`
	NVRA  string `json:"nvra,omitempty" yaml:"nvra,omitempty"`
}
