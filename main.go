package main

import (
	"github.com/dihedron/steampipe-plugin-satellite/satellite"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: satellite.Plugin})
}
