package satellite

import (
	"errors"
	"strings"
)

// ParseNVRA returns name, version, release, arch from an RPM NVREA e.g.:
// - foo-1.0-1.i386.rpm returns foo, 1.0, 1, i386
// - bar-9-123a.ia64.rpm returns bar, 9, 123a, 1, ia64
// This is the Golang porting of the official implementation in Python's in rpmUtils.miscutils.
func ParseNVRA(nvra string) (name string, ver string, rel string, arch string, err error) {
	// tuned-profiles-cpu-partitioning-2.18.0-1.2.20220511git9fa66f19.el8fdp.noarch
	//                                                                      ^----->|
	archIndex := strings.LastIndex(nvra, ".")
	if archIndex == -1 {
		return "", "", "", "", errors.New("invalid format: no arch info")
	}
	arch = nvra[archIndex+1:]

	// tuned-profiles-cpu-partitioning-2.18.0-1.2.20220511git9fa66f19.el8fdp.noarch
	//                                       ^----------------------------->|
	relIndex := strings.LastIndex(nvra[:archIndex], "-")
	rel = nvra[relIndex+1 : archIndex]

	// tuned-profiles-cpu-partitioning-2.18.0-1.2.20220511git9fa66f19.el8fdp.noarch
	//                                ^----->|
	verIndex := strings.LastIndex(nvra[:relIndex], "-")
	ver = nvra[verIndex+1 : relIndex]

	// tuned-profiles-cpu-partitioning-2.18.0-1.2.20220511git9fa66f19.el8fdp.noarch
	// ------------------------------>|
	name = nvra[:verIndex]
	return
}
