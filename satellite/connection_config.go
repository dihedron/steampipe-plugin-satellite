package satellite

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type satelliteConfig struct {
	Endpoint     *string `cty:"endpoint"`
	Username     *string `cty:"username"`
	Password     *string `cty:"password"`
	Organisation *string `cty:"organisation"`
	Location     *string `cty:"location"`
	TraceLevel   *string `cty:"trace_level"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"endpoint": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"organisation": {
		Type: schema.TypeString,
	},
	"location": {
		Type: schema.TypeString,
	},
	"trace_level": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &satelliteConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) satelliteConfig {
	if connection == nil || connection.Config == nil {
		return satelliteConfig{}
	}
	config, _ := connection.Config.(satelliteConfig)
	return config
}
