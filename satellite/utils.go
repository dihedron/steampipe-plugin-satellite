package satellite

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

var ErrNotImplemented = errors.New("not implemented")

// setLogLevel changes the current HCLog level; this seems necessary as the
// STEAMPIPE_LOG_LEVEL variable does not seem to be properly read by the plugins.
func setLogLevel(ctx context.Context, d *plugin.QueryData) {
	satelliteConfig := GetConfig(d.Connection)
	if satelliteConfig.TraceLevel != nil {
		level := *satelliteConfig.TraceLevel
		plugin.Logger(ctx).SetLevel(hclog.LevelFromString(level))
	}
}

// toPrettyJSON dumps the input object to JSON.
func toPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

// pointerTo returns a pointer to a given value.
// func pointerTo[T any](value T) *T {
// 	return &value
// }