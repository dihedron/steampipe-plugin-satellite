package satellite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableSatelliteHost(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "satellite_host",
		Description: "Red Hat Satellite Host",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "The host id",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the host.",
			},
			{
				Name:        "organization",
				Type:        proto.ColumnType_STRING,
				Description: "The organisation managing the host.",
				Transform:   transform.FromField("OrganizationName"),
			},
			{
				Name:        "model",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's model.",
				Transform:   transform.FromField("ModelName"),
			},
			{
				Name:        "ipv4_address",
				Type:        proto.ColumnType_STRING,
				Description: "The IPv4 address of the host.",
				Transform:   transform.FromField("IPv4"),
			},
			{
				Name:        "ipv6_address",
				Type:        proto.ColumnType_STRING,
				Description: "The IPv6 address of the host.",
				Transform:   transform.FromField("IPv6"),
			},
			{
				Name:        "mac_address",
				Type:        proto.ColumnType_STRING,
				Description: "The MAC address of the host.",
				Transform:   transform.FromField("MACAddress"),
			},
			{
				Name:        "architecture",
				Type:        proto.ColumnType_STRING,
				Description: "The machine architecture of the host.",
				Transform:   transform.FromField("ArchitectureName"),
			},
			{
				Name:        "operating_system",
				Type:        proto.ColumnType_STRING,
				Description: "The operating system of the host.",
				Transform:   transform.FromField("OperatingSystemName"),
			},
			{
				Name:        "environment",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the machine's environment.",
				Transform:   transform.FromField("EnvironmentName"),
			},
			{
				Name:        "host_group_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the machine's host group.",
				Transform:   transform.FromField("HostGroupName"),
			},
			{
				Name:        "host_group_title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the machine's host group.",
				Transform:   transform.FromField("HostGroupTitle"),
			},
			{
				Name:        "compute_resource",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's compute resource.",
				Transform:   transform.FromField("ComputeResourceName"),
			},
			{
				Name:        "compute_profile",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's compute profile.",
				Transform:   transform.FromField("ComputeProfileName"),
			},
			{
				Name:        "realm",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the machine's realm.",
				Transform:   transform.FromField("RealmName"),
			},
			{
				Name:        "image_file",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's image file.",
				Transform:   transform.FromField("ImageFile"),
			},
			{
				Name:        "provision_method",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's provision method.",
				Transform:   transform.FromField("ProvisionMethod"),
			},
			{
				Name:        "pxe_loader",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's PXE loader info.",
				Transform:   transform.FromField("PXELoader"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's creation time.",
				Transform: &transform.ColumnTransforms{
					Transforms: []*transform.TransformCall{
						{
							Transform: transform.FieldValueGo,
							Param:     "CreatedAt",
						},
						{
							Transform: func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
								if t, ok := d.Value.(fmt.Stringer); ok {
									return t.String(), nil
								}
								return nil, fmt.Errorf("invalid data type: %T", d.Value)
							},
							Param: nil,
						},
					},
				},
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's update time.",
				Transform:   FromStringerField("UpdatedAt"),
			},
			{
				Name:        "installed_at",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's installation time.",
				Transform:   FromStringerField("InstalledAt"),
			},
			{
				Name:        "enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "The machine's enablement status.",
				Transform:   transform.FromField("Enabled"),
			},
			{
				Name:        "managed",
				Type:        proto.ColumnType_STRING,
				Description: "Whether the machine is managed.",
				Transform:   transform.FromField("Managed"),
			},
			{
				Name:        "uptime_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The machine's uptime in seconds.",
				Transform:   transform.FromField("UptimeSeconds"),
			},
			{
				Name:        "uptime_duration",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's uptime in seconds.",
				Transform: &transform.ColumnTransforms{
					Transforms: []*transform.TransformCall{
						{
							Transform: transform.FieldValue,
							Param:     "UptimeSeconds",
						},
						{
							Transform: func(ctx context.Context, d *transform.TransformData) (any, error) {
								return time.Duration(d.Value.(int) * 1_000_000_000).Round(time.Second).String(), nil
							},
							Param: nil,
						},
					},
				},
			},
			{
				Name:        "global_status",
				Type:        proto.ColumnType_STRING,
				Description: "The global machine status.",
				Transform:   transform.FromField("GlobalStatusLabel"),
			},
			{
				Name:        "errata_status",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's errata status.",
				Transform:   transform.FromField("ErrataStatusLabel"),
			},
			{
				Name:        "purpose_status",
				Type:        proto.ColumnType_STRING,
				Description: "The machine's purpose status.",
				Transform:   transform.FromField("PurposeStatusLabel"),
			},
			{
				Name:        "subscription_status",
				Type:        proto.ColumnType_STRING,
				Description: "The machine subscription status.",
				Transform:   transform.FromField("SubscriptionStatusLabel"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSatelliteHost,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
			// plugin.SingleColumn("id"),
			Hydrate: getSatelliteHost,
		},
	}
}

//// LIST FUNCTIONS

func listSatelliteHost(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	setLogLevel(ctx, d)
	plugin.Logger(ctx).Debug("retrieving satellite host list", "query data", toPrettyJSON(d))

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving satellite client", "error", err)
		return nil, err
	}

	request := client.
		R().
		SetContext(ctx).
		SetQueryParam("thin", "true") // TODO: temporarily gather only part of the data

	result := &struct {
		Total    int    `json:"total"`
		Subtotal int    `json:"subtotal"`
		Page     int    `json:"page"`
		PerPage  int    `json:"per_page"`
		Search   string `json:"search"`
		Sort     struct {
			By    string `json:"by"`
			Order string `json:"order"`
		} `json:"sort"`
		Hosts []apiHost `json:"results"`
	}{}
	request.SetResult(result)

	_, err = request.Get("/hosts")
	if err != nil {
		plugin.Logger(ctx).Error("error performing request", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("request successful", "total", result.Total, "subtotal", result.Subtotal, "page", result.Page, "per page", result.PerPage)

	for _, host := range result.Hosts {
		d.StreamListItem(ctx, &host)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSatelliteHost(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	setLogLevel(ctx, d)

	id := ""
	if val, ok := d.KeyColumnQuals["id"]; ok {
		plugin.Logger(ctx).Debug("retrieving satellite host by id", "id", id)
		id = fmt.Sprintf("%d", val.GetInt64Value())
	} else if val, ok := d.KeyColumnQuals["name"]; ok {
		plugin.Logger(ctx).Debug("retrieving satellite host by name", "name", id)
		id = val.GetStringValue()
	} else {
		plugin.Logger(ctx).Error("no valid key provided")
		return nil, errors.New("no valid key provided")
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving satellite client", "error", err)
		return nil, err
	}

	request := client.
		R().
		SetContext(ctx).
		SetQueryParam("thin", "true") // TODO: temporarily gather only part of the data

	request = request.SetPathParam("id", id)

	host := &apiHost{}
	request.SetResult(host)
	_, err = request.Get("/hosts/{id}")
	if err != nil {
		plugin.Logger(ctx).Error("error performing request", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("request successful", "host", toPrettyJSON(host))

	return host, nil
}

type apiHost struct {
	IPv4                     string      `json:"ip,omitempty" yaml:"ip,omitempty"`
	IPv6                     string      `json:"ip6,omitempty" yaml:"ip6,omitempty"`
	EnvironmentID            int         `json:"environment_id,omitempty" yaml:"environment_id,omitempty"`
	EnvironmentName          string      `json:"environment_name,omitempty" yaml:"environment_name,omitempty"`
	LastReport               interface{} `json:"last_report,omitempty" yaml:"last_report,omitempty"`
	MACAddress               string      `json:"mac,omitempty" yaml:"mac,omitempty"`
	RealmID                  int         `json:"realm_id,omitempty" yaml:"realm_id,omitempty"`
	RealmName                string      `json:"realm_name,omitempty" yaml:"realm_name,omitempty"`
	SpMAC                    interface{} `json:"sp_mac,omitempty" yaml:"sp_mac,omitempty"`
	SpIP                     interface{} `json:"sp_ip,omitempty" yaml:"sp_ip,omitempty"`
	SpName                   interface{} `json:"sp_name,omitempty" yaml:"sp_name,omitempty"`
	DomainID                 int         `json:"domain_id,omitempty" yaml:"domain_id,omitempty"`
	DomainName               string      `json:"domain_name,omitempty" yaml:"domain_name,omitempty"`
	ArchitectureID           int         `json:"architecture_id,omitempty" yaml:"architecture_id,omitempty"`
	ArchitectureName         string      `json:"architecture_name,omitempty" yaml:"architecture_name,omitempty"`
	OperatingSystemID        int         `json:"operatingsystem_id,omitempty" yaml:"operatingsystem_id,omitempty"`
	OperatingSystemName      string      `json:"operatingsystem_name,omitempty" yaml:"operatingsystem_name,omitempty"`
	SubnetID                 interface{} `json:"subnet_id,omitempty" yaml:"subnet_id,omitempty"`
	SubnetName               interface{} `json:"subnet_name,omitempty" yaml:"subnet_name,omitempty"`
	Subnet6ID                interface{} `json:"subnet6_id,omitempty" yaml:"subnet6_id,omitempty"`
	Subnet6Name              interface{} `json:"subnet6_name,omitempty" yaml:"subnet6_name,omitempty"`
	SpSubnetID               int         `json:"sp_subnet_id,omitempty" yaml:"sp_subnet_id,omitempty"`
	PTableID                 int         `json:"ptable_id,omitempty" yaml:"ptable_id,omitempty"`
	PTableName               string      `json:"ptable_name,omitempty" yaml:"ptable_name,omitempty"`
	MediumID                 interface{} `json:"medium_id,omitempty" yaml:"medium_id,omitempty"`
	MediumName               interface{} `json:"medium_name,omitempty" yaml:"medium_name,omitempty"`
	PXELoader                string      `json:"pxe_loader,omitempty" yaml:"pxe_loader,omitempty"`
	Build                    bool        `json:"build,omitempty" yaml:"build,omitempty"`
	Comment                  interface{} `json:"comment,omitempty" yaml:"comment,omitempty"`
	Disk                     interface{} `json:"disk,omitempty" yaml:"disk,omitempty"`
	InstalledAt              *Time       `json:"installed_at,omitempty" yaml:"installed_at,omitempty"`
	ModelID                  int         `json:"model_id,omitempty" yaml:"model_id,omitempty"`
	OwnerID                  int         `json:"owner_id,omitempty" yaml:"owner_id,omitempty"`
	OwnerName                string      `json:"owner_name,omitempty" yaml:"owner_name,omitempty"`
	OwnerType                string      `json:"owner_type,omitempty" yaml:"owner_type,omitempty"`
	Enabled                  bool        `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Managed                  bool        `json:"managed,omitempty" yaml:"managed,omitempty"`
	UseImage                 interface{} `json:"use_image,omitempty" yaml:"use_image,omitempty"`
	ImageFile                string      `json:"image_file,omitempty" yaml:"image_file,omitempty"`
	UUID                     interface{} `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	ComputeResourceID        int         `json:"compute_resource_id,omitempty" yaml:"compute_resource_id,omitempty"`
	ComputeResourceName      string      `json:"compute_resource_name,omitempty" yaml:"compute_resource_name,omitempty"`
	ComputeProfileID         int         `json:"compute_profile_id,omitempty" yaml:"compute_profile_id,omitempty"`
	ComputeProfileName       string      `json:"compute_profile_name,omitempty" yaml:"compute_profile_name,omitempty"`
	Capabilities             []string    `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	ProvisionMethod          string      `json:"provision_method,omitempty" yaml:"provision_method,omitempty"`
	CertName                 string      `json:"certname,omitempty" yaml:"certname,omitempty"`
	ImageID                  interface{} `json:"image_id,omitempty" yaml:"image_id,omitempty"`
	ImageName                interface{} `json:"image_name,omitempty" yaml:"image_name,omitempty"`
	CreatedAt                *Time       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt                *Time       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	LastCompile              string      `json:"last_compile,omitempty" yaml:"last_compile,omitempty"`
	GlobalStatus             int         `json:"global_status,omitempty" yaml:"global_status,omitempty"`
	GlobalStatusLabel        string      `json:"global_status_label,omitempty" yaml:"global_status_label,omitempty"`
	UptimeSeconds            int         `json:"uptime_seconds,omitempty" yaml:"uptime_seconds,omitempty"`
	OrganizationID           int         `json:"organization_id,omitempty" yaml:"organization_id,omitempty"`
	OrganizationName         string      `json:"organization_name,omitempty" yaml:"organization_name,omitempty"`
	LocationID               int         `json:"location_id,omitempty" yaml:"location_id,omitempty"`
	LocationName             string      `json:"location_name,omitempty" yaml:"location_name,omitempty"`
	PuppetStatus             int         `json:"puppet_status,omitempty" yaml:"puppet_status,omitempty"`
	ModelName                string      `json:"model_name,omitempty" yaml:"model_name,omitempty"`
	ErrataStatus             int         `json:"errata_status,omitempty" yaml:"errata_status,omitempty"`
	ErrataStatusLabel        string      `json:"errata_status_label,omitempty" yaml:"errata_status_label,omitempty"`
	SubscriptionStatus       int         `json:"subscription_status,omitempty" yaml:"subscription_status,omitempty"`
	SubscriptionStatusLabel  string      `json:"subscription_status_label,omitempty" yaml:"subscription_status_label,omitempty"`
	PurposeSLAStatus         int         `json:"purpose_sla_status,omitempty" yaml:"purpose_sla_status,omitempty"`
	PurposeSLAStatusLabel    string      `json:"purpose_sla_status_label,omitempty" yaml:"purpose_sla_status_label,omitempty"`
	PurposeRoleStatus        int         `json:"purpose_role_status,omitempty" yaml:"purpose_role_status,omitempty"`
	PurposeRoleStatusLabel   string      `json:"purpose_role_status_label,omitempty" yaml:"purpose_role_status_label,omitempty"`
	PurposeUsageStatus       int         `json:"purpose_usage_status,omitempty" yaml:"purpose_usage_status,omitempty"`
	PurposeUsageStatusLabel  string      `json:"purpose_usage_status_label,omitempty" yaml:"purpose_usage_status_label,omitempty"`
	PurposeAddonsStatus      int         `json:"purpose_addons_status,omitempty" yaml:"purpose_addons_status,omitempty"`
	PurposeAddonsStatusLabel string      `json:"purpose_addons_status_label,omitempty" yaml:"purpose_addons_status_label,omitempty"`
	PurposeStatus            int         `json:"purpose_status,omitempty" yaml:"purpose_status,omitempty"`
	PurposeStatusLabel       string      `json:"purpose_status_label,omitempty" yaml:"purpose_status_label,omitempty"`
	Name                     string      `json:"name,omitempty" yaml:"name,omitempty"`
	ID                       int         `json:"id,omitempty" yaml:"id,omitempty"`
	PuppetProxyID            interface{} `json:"puppet_proxy_id,omitempty" yaml:"puppet_proxy_id,omitempty"`
	PuppetProxyName          string      `json:"puppet_proxy_name,omitempty" yaml:"puppet_proxy_name,omitempty"`
	PuppetCaProxyID          interface{} `json:"puppet_ca_proxy_id,omitempty" yaml:"puppet_ca_proxy_id,omitempty"`
	PuppetCaProxyName        string      `json:"puppet_ca_proxy_name,omitempty" yaml:"puppet_ca_proxy_name,omitempty"`
	OpenSCAPProxyID          interface{} `json:"openscap_proxy_id,omitempty" yaml:"openscap_proxy_id,omitempty"`
	OpenSCAPProxyName        string      `json:"openscap_proxy_name,omitempty" yaml:"openscap_proxy_name,omitempty"`
	PuppetProxy              interface{} `json:"puppet_proxy,omitempty" yaml:"puppet_proxy,omitempty"`
	PuppetCaProxy            interface{} `json:"puppet_ca_proxy,omitempty" yaml:"puppet_ca_proxy,omitempty"`
	OpenSCAPProxy            interface{} `json:"openscap_proxy,omitempty" yaml:"openscap_proxy,omitempty"`
	HostGroupID              int         `json:"hostgroup_id,omitempty" yaml:"hostgroup_id,omitempty"`
	HostGroupName            string      `json:"hostgroup_name,omitempty" yaml:"hostgroup_name,omitempty"`
	HostGroupTitle           string      `json:"hostgroup_title,omitempty" yaml:"hostgroup_title,omitempty"`
	Parameters               []struct {
		Priority      int         `json:"priority,omitempty" yaml:"priority,omitempty"`
		CreatedAt     string      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		UpdatedAt     string      `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
		ID            int         `json:"id,omitempty" yaml:"id,omitempty"`
		Name          string      `json:"name,omitempty" yaml:"name,omitempty"`
		ParameterType string      `json:"parameter_type" yaml:"parameter_type"`
		Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	AllParameters []struct {
		Priority      int         `json:"priority,omitempty" yaml:"priority,omitempty"`
		CreatedAt     string      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		UpdatedAt     string      `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
		ID            int         `json:"id,omitempty" yaml:"id,omitempty"`
		Name          string      `json:"name,omitempty" yaml:"name,omitempty"`
		ParameterType string      `json:"parameter_type" yaml:"parameter_type"`
		Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	} `json:"all_parameters,omitempty" yaml:"all_parameters,omitempty"`
	ContentFacetAttributes struct {
		ID                       int    `json:"id,omitempty" yaml:"id,omitempty"`
		UUID                     string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
		ContentViewID            int    `json:"content_view_id,omitempty" yaml:"content_view_id,omitempty"`
		ContentViewName          string `json:"content_view_name,omitempty" yaml:"content_view_name,omitempty"`
		LifecycleEnvironmentID   int    `json:"lifecycle_environment_id,omitempty" yaml:"lifecycle_environment_id,omitempty"`
		LifecycleEnvironmentName string `json:"lifecycle_environment_name,omitempty" yaml:"lifecycle_environment_name,omitempty"`
		ContentSourceID          int    `json:"content_source_id,omitempty" yaml:"content_source_id,omitempty"`
		ContentSourceName        string `json:"content_source_name,omitempty" yaml:"content_source_name,omitempty"`
		KickstartRepositoryID    int    `json:"kickstart_repository_id,omitempty" yaml:"kickstart_repository_id,omitempty"`
		KickstartRepositoryName  string `json:"kickstart_repository_name,omitempty" yaml:"kickstart_repository_name,omitempty"`
		ErrataCounts             struct {
			Security    int `json:"security,omitempty" yaml:"security,omitempty"`
			Bugfix      int `json:"bugfix,omitempty" yaml:"bugfix,omitempty"`
			Enhancement int `json:"enhancement,omitempty" yaml:"enhancement,omitempty"`
			Total       int `json:"total,omitempty" yaml:"total,omitempty"`
		} `json:"errata_counts,omitempty" yaml:"errata_counts,omitempty"`
		ApplicablePackageCount      int `json:"applicable_package_count,omitempty" yaml:"applicable_package_count,omitempty"`
		UpgradablePackageCount      int `json:"upgradable_package_count,omitempty" yaml:"upgradable_package_count,omitempty"`
		ApplicableModuleStreamCount int `json:"applicable_module_stream_count,omitempty" yaml:"applicable_module_stream_count,omitempty"`
		UpgradableModuleStreamCount int `json:"upgradable_module_stream_count,omitempty" yaml:"upgradable_module_stream_count,omitempty"`
		ContentView                 struct {
			ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
			Name string `json:"name,omitempty" yaml:"name,omitempty"`
		} `json:"content_view,omitempty" yaml:"content_view,omitempty"`
		LifecycleEnvironment struct {
			ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
			Name string `json:"name,omitempty" yaml:"name,omitempty"`
		} `json:"lifecycle_environment,omitempty" yaml:"lifecycle_environment,omitempty"`
		ContentSource struct {
			ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
			Name string `json:"name,omitempty" yaml:"name,omitempty"`
			URL  string `json:"url,omitempty" yaml:"url,omitempty"`
		} `json:"content_source,omitempty" yaml:"content_source,omitempty"`
		KickstartRepository         interface{} `json:"kickstart_repository,omitempty" yaml:"kickstart_repository,omitempty"`
		ContentViewVersion          string      `json:"content_view_version,omitempty" yaml:"content_view_version,omitempty"`
		ContentViewVersionID        int         `json:"content_view_version_id,omitempty" yaml:"content_view_version_id,omitempty"`
		ContentViewDefault          bool        `json:"content_view_default?,omitempty" yaml:"content_view_default?,omitempty"`
		LifecycleEnvironmentLibrary bool        `json:"lifecycle_environment_library?,omitempty" yaml:"lifecycle_environment_library?,omitempty"`
		KatelloAgentInstalled       bool        `json:"katello_agent_installed,omitempty" yaml:"katello_agent_installed,omitempty"`
		KatelloTracerInstalled      bool        `json:"katello_tracer_installed,omitempty" yaml:"katello_tracer_installed,omitempty"`
	} `json:"content_facet_attributes,omitempty" yaml:"content_facet_attributes,omitempty"`
	SubscriptionGlobalStatus    int `json:"subscription_global_status,omitempty" yaml:"subscription_global_status,omitempty"`
	SubscriptionFacetAttributes struct {
		HostType          string        `json:"host_type,omitempty" yaml:"host_type,omitempty"`
		DmiUUID           string        `json:"dmi_uuid,omitempty" yaml:"dmi_uuid,omitempty"`
		ID                int           `json:"id,omitempty" yaml:"id,omitempty"`
		UUID              string        `json:"uuid,omitempty" yaml:"uuid,omitempty"`
		LastCheckin       string        `json:"last_checkin,omitempty" yaml:"last_checkin,omitempty"`
		ServiceLevel      string        `json:"service_level,omitempty" yaml:"service_level,omitempty"`
		ReleaseVersion    interface{}   `json:"release_version,omitempty" yaml:"release_version,omitempty"`
		Autoheal          bool          `json:"autoheal,omitempty" yaml:"autoheal,omitempty"`
		RegisteredAt      *Time         `json:"registered_at,omitempty" yaml:"registered_at,omitempty"`
		RegisteredThrough string        `json:"registered_through,omitempty" yaml:"registered_through,omitempty"`
		PurposeRole       string        `json:"purpose_role,omitempty" yaml:"purpose_role,omitempty"`
		PurposeUsage      string        `json:"purpose_usage,omitempty" yaml:"purpose_usage,omitempty"`
		Hypervisor        bool          `json:"hypervisor,omitempty" yaml:"hypervisor,omitempty"`
		User              interface{}   `json:"user,omitempty" yaml:"user,omitempty"`
		PurposeAddons     []interface{} `json:"purpose_addons,omitempty" yaml:"purpose_addons,omitempty"`
		VirtualHost       interface{}   `json:"virtual_host,omitempty" yaml:"virtual_host,omitempty"`
		VirtualGuests     []struct {
			ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
			Name string `json:"name,omitempty" yaml:"name,omitempty"`
		} `json:"virtual_guests,omitempty" yaml:"virtual_guests,omitempty"`
		InstalledProducts []struct {
			ProductName string `json:"productName,omitempty" yaml:"productName,omitempty"`
			ProductID   string `json:"productId,omitempty" yaml:"productId,omitempty"`
			Arch        string `json:"arch,omitempty" yaml:"arch,omitempty"`
			Version     string `json:"version,omitempty" yaml:"version,omitempty"`
		} `json:"installed_products,omitempty" yaml:"installed_products,omitempty"`
		ActivationKeys []struct {
			ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
			Name string `json:"name,omitempty" yaml:"name,omitempty"`
		} `json:"activation_keys,omitempty" yaml:"activation_keys,omitempty"`
		ComplianceReasons []interface{} `json:"compliance_reasons,omitempty" yaml:"compliance_reasons,omitempty"`
	} `json:"subscription_facet_attributes,omitempty" yaml:"subscription_facet_attributes,omitempty"`
	HostCollections []interface{} `json:"host_collections,omitempty" yaml:"host_collections,omitempty"`
	Interfaces      []struct {
		SubnetID        int    `json:"subnet_id,omitempty" yaml:"subnet_id,omitempty"`
		SubnetName      string `json:"subnet_name,omitempty" yaml:"subnet_name,omitempty"`
		Subnet6ID       int    `json:"subnet6_id,omitempty" yaml:"subnet6_id,omitempty"`
		Subnet6Name     string `json:"subnet6_name,omitempty" yaml:"subnet6_name,omitempty"`
		DomainID        int    `json:"domain_id,omitempty" yaml:"domain_id,omitempty"`
		DomainName      string `json:"domain_name,omitempty" yaml:"domain_name,omitempty"`
		CreatedAt       string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		UpdatedAt       string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
		Managed         bool   `json:"managed,omitempty" yaml:"managed,omitempty"`
		Identifier      string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		ID              int    `json:"id,omitempty" yaml:"id,omitempty"`
		Name            string `json:"name,omitempty" yaml:"name,omitempty"`
		IPv4            string `json:"ip,omitempty" yaml:"ip,omitempty"`
		IPv6            string `json:"ip6,omitempty" yaml:"ip6,omitempty"`
		MAC             string `json:"mac,omitempty" yaml:"mac,omitempty"`
		MTU             int    `json:"mtu,omitempty" yaml:"mtu,omitempty"`
		FQDN            string `json:"fqdn,omitempty" yaml:"fqdn,omitempty"`
		Primary         bool   `json:"primary,omitempty" yaml:"primary,omitempty"`
		Provision       bool   `json:"provision,omitempty" yaml:"provision,omitempty"`
		Type            string `json:"type,omitempty" yaml:"type,omitempty"`
		Execution       bool   `json:"execution,omitempty" yaml:"execution,omitempty"`
		Mode            string `json:"mode,omitempty" yaml:"mode,omitempty"`
		AttachedDevices string `json:"attached_devices,omitempty" yaml:"attached_devices,omitempty"`
		BondOptions     string `json:"bond_options,omitempty" yaml:"bond_options,omitempty"`
		Virtual         bool   `json:"virtual,omitempty" yaml:"virtual,omitempty"`
		Tag             string `json:"tag,omitempty" yaml:"tag,omitempty"`
		AttachedTo      string `json:"attached_to,omitempty" yaml:"attached_to,omitempty"`
	} `json:"interfaces,omitempty"`
	Facts       map[string]string `json:"facts,omitempty" yaml:"facts,omitempty"`
	Permissions struct {
		CockpitHosts                 bool `json:"cockpit_hosts,omitempty" yaml:"cockpit_hosts,omitempty"`
		ViewHosts                    bool `json:"view_hosts,omitempty"  yaml:"view_hosts,omitempty"`
		CreateHosts                  bool `json:"create_hosts,omitempty" yaml:"create_hosts,omitempty"`
		EditHosts                    bool `json:"edit_hosts,omitempty" yaml:"edit_hosts,omitempty"`
		DestroyHosts                 bool `json:"destroy_hosts,omitempty" yaml:"destroy_hosts,omitempty"`
		BuildHosts                   bool `json:"build_hosts,omitempty" yaml:"build_hosts,omitempty"`
		PowerHosts                   bool `json:"power_hosts,omitempty" yaml:"power_hosts,omitempty"`
		ConsoleHosts                 bool `json:"console_hosts,omitempty" yaml:"console_hosts,omitempty"`
		IpmiBootHosts                bool `json:"ipmi_boot_hosts,omitempty" yaml:"ipmi_boot_hosts,omitempty"`
		ViewDiscoveredHosts          bool `json:"view_discovered_hosts,omitempty" yaml:"view_discovered_hosts,omitempty"`
		SubmitDiscoveredHosts        bool `json:"submit_discovered_hosts,omitempty" yaml:"submit_discovered_hosts,omitempty"`
		AutoProvisionDiscoveredHosts bool `json:"auto_provision_discovered_hosts,omitempty" yaml:"auto_provision_discovered_hosts,omitempty"`
		ProvisionDiscoveredHosts     bool `json:"provision_discovered_hosts,omitempty" yaml:"provision_discovered_hosts,omitempty"`
		EditDiscoveredHosts          bool `json:"edit_discovered_hosts,omitempty" yaml:"edit_discovered_hosts,omitempty"`
		DestroyDiscoveredHosts       bool `json:"destroy_discovered_hosts,omitempty" yaml:"destroy_discovered_hosts,omitempty"`
		PlayRolesOnHost              bool `json:"play_roles_on_host,omitempty" yaml:"play_roles_on_host,omitempty"`
		ForgetStatusHosts            bool `json:"forget_status_hosts,omitempty" yaml:"forget_status_hosts,omitempty"`
	} `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	PuppetClasses            []interface{} `json:"puppetclasses,omitempty" yaml:"puppetclasses,omitempty"`
	ConfigGroups             []interface{} `json:"config_groups,omitempty" yaml:"config_groups,omitempty"`
	AllPuppetClasses         []interface{} `json:"all_puppetclasses,omitempty" yaml:"all_puppetclasses,omitempty"`
	ConfigurationStatus      int           `json:"configuration_status,omitempty" yaml:"configuration_status,omitempty"`
	ConfigurationStatusLabel string        `json:"configuration_status_label,omitempty" yaml:"configuration_status_label,omitempty"`
	BuildStatus              int           `json:"build_status,omitempty" yaml:"build_status,omitempty"`
	BuildStatusLabel         string        `json:"build_status_label,omitempty" yaml:"build_status_label,omitempty"` // test
}
