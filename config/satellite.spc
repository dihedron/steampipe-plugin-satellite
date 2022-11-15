connection "satellite" {
    # the path to the plugin
    plugin    = "local/satellite"
    # Red Hat Satellite connection info
    endpoint_url = "https://satellite.example.com/api"
    username = "<username>"
    password = "<password>"
    organisation = "<organisation>"
    trace_level = "TRACE"
}
