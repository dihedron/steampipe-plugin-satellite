.PHONY: plugin
plugin:
	@go build

.PHONY: clean
clean:
	@rm -rf steampipe-plugin-satellite

.PHONY: install
install: plugin
	@mkdir -p ~/.steampipe/plugins/local/satellite
	@cp steampipe-plugin-satellite ~/.steampipe/plugins/local/satellite/satellite.plugin
#	@cp config/satellite.spc ~/.steampipe/config/satellite.spc