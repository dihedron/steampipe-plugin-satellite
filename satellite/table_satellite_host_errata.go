package satellite

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSatelliteHostErrata(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "satellite_host_errata",
		Description: "Red Hat Satellite Host Errata",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "The id of the errata.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "pulp_id",
				Type:        proto.ColumnType_STRING,
				Description: "The pulp ID of the errata.",
				Transform:   transform.FromField("PulpID"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the errata.",
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "errata_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the errata.",
				Transform:   transform.FromField("ErrataID"),
			},
			{
				Name:        "issued_at",
				Type:        proto.ColumnType_STRING,
				Description: "The time when the errata was issued.",
				Transform:   transform.FromField("Issued"),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The time when the errata updated.",
				Transform:   transform.FromField("Updated"),
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_STRING,
				Description: "The severity of the errata.",
				Transform:   transform.FromField("Severity"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the errata.",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "solution",
				Type:        proto.ColumnType_STRING,
				Description: "The solution of the errata.",
				Transform:   transform.FromField("Solution"),
			},
			{
				Name:        "summary",
				Type:        proto.ColumnType_STRING,
				Description: "The summary of the errata.",
				Transform:   transform.FromField("Summary"),
			},
			{
				Name:        "reboot_suggested",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether a reboot is suggested to fix the errata.",
				Transform:   transform.FromField("RebootSuggested"),
			},
			{
				Name:        "uuid",
				Type:        proto.ColumnType_STRING,
				Description: "The UUID of the errata.",
				Transform:   transform.FromField("UUID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the errata.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the errata.",
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "hosts_available_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of hosts on which the errata is available.",
				Transform:   transform.FromField("HostsAvailableCount"),
			},
			{
				Name:        "hosts_applicable_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of hosts on which the errata is applicable.",
				Transform:   transform.FromField("HostsApplicableCount"),
			},
			{
				Name:        "packages",
				Type:        proto.ColumnType_JSON,
				Description: "The packages related to the errata.",
				Transform:   transform.FromField("Packages"),
			},
			{
				Name:        "installable",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the errata is installable.",
				Transform:   transform.FromField("Installable"),
			},
			{
				Name:        "cves",
				Type:        proto.ColumnType_JSON,
				Description: "The CVE applicable to this host.",
				Transform:   transform.FromField("CVEs"),
			},
			{
				Name:        "bugs",
				Type:        proto.ColumnType_JSON,
				Description: "The bugs applicable to this host.",
				Transform:   transform.FromField("Bugs"),
			},
			{
				Name:        "module_streams",
				Type:        proto.ColumnType_JSON,
				Description: "The module streams applicable to this host.",
				Transform:   transform.FromField("ModuleStreams"),
			},
			// join columns
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
			Hydrate: listSatelliteHostErrata,
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

func listSatelliteHostErrata(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	setLogLevel(ctx, d)
	plugin.Logger(ctx).Debug("retrieving satellite errata list")

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving satellite client", "error", err)
		return nil, err
	}

	id := ""
	hostid := int(d.EqualsQuals["host_id"].GetInt64Value())
	hostname := d.EqualsQuals["host_name"].GetStringValue()
	if hostid == 0 {
		hostid, err = resolveHostID(ctx, d, hostname)
		if err != nil {
			plugin.Logger(ctx).Error("error resolving host by name", "name", hostname)
			return nil, err
		}
	}
	id = fmt.Sprintf("%d", hostid)
	if id == "" {
		plugin.Logger(ctx).Error("no valid host id or name provided")
		return nil, errors.New("no valid host id or name provided")
	}

	page := 1
loop:
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
			Errata []apiErrata `json:"results"`
		}{}
		request.SetResult(result)
		response, err := request.Get("/api/hosts/{id}/errata")
		if err != nil || response.IsError() {
			plugin.Logger(ctx).Error("error performing request", "url", response.Request.URL, "status", response, response.Status(), "error", err, "response", utils.ToPrettyJSON(response.Body()))
			return nil, fmt.Errorf("request %q failed with status %d (%s): %w", response.Request.URL, response.StatusCode(), response.Status(), err)
		}
		plugin.Logger(ctx).Debug("request successful", "total", result.Total, "subtotal", result.Subtotal, "page", result.Page, "per page", result.PerPage, "response", utils.ToJSON(response.Body()))

		for _, errata := range result.Errata {
			if ctx.Err() != nil {
				plugin.Logger(ctx).Debug("context done, exit")
				break loop
			}

			errata := errata
			//plugin.Logger(ctx).Debug("package", "contents", toJSON(pkg))
			d.StreamListItem(ctx, &struct {
				HostID   int    `json:"host_id,omitempty" yaml:"host_id,omitempty"`
				HostName string `json:"host_name,omitempty" yaml:"host_name,omitempty"`
				apiErrata
			}{
				HostID:    int(hostid),
				HostName:  hostname,
				apiErrata: errata,
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

type apiErrata struct {
	ID              int    `json:"id"`
	PulpID          string `json:"pulp_id"`
	Title           string `json:"title"`
	ErrataID        string `json:"errata_id"`
	Issued          string `json:"issued"`
	Updated         string `json:"updated"`
	Severity        string `json:"severity"`
	Description     string `json:"description"`
	Solution        string `json:"solution"`
	Summary         string `json:"summary"`
	RebootSuggested bool   `json:"reboot_suggested"`
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	CVEs            []struct {
		ID   string `json:"bug_id"`
		Href string `json:"href"`
	} `json:"cves"`
	Bugs []struct {
		ID   string `json:"bug_id"`
		Href string `json:"href"`
	} `json:"bugs"`
	HostsAvailableCount  int           `json:"hosts_available_count"`
	HostsApplicableCount int           `json:"hosts_applicable_count"`
	Packages             []string      `json:"packages"`
	ModuleStreams        []interface{} `json:"module_streams"`
	Installable          bool          `json:"installable"`
}
