package buildinfo

import "runtime"

var (
	Version    = "undefined"
	Revision   = "undefined"
	Branch     = "undefined"
	Prerelease = "undefined"
	IsSnapshot = "undefined"
	BuildUser  = "undefined"
	BuildDate  = "undefined"
)

// Info represents all available build information.
type Info struct {
	Version    string `json:"version"`
	Revision   string `json:"revision"`
	Branch     string `json:"branch"`
	Prerelease string `json:"prerelease"`
	IsSnapshot string `json:"is_snapshot"`
	BuildUser  string `json:"build_user"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	Compiler   string `json:"compiler"`
}

// NewInfo returns all available build information.
func NewInfo(
	version string,
	revision string,
	branch string,
	pre string,
	snap string,
	buildUser string,
	buildDate string,
) Info {
	return Info{
		Version:    version,
		Revision:   revision,
		Branch:     branch,
		Prerelease: pre,
		IsSnapshot: snap,
		BuildUser:  buildUser,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Compiler:   runtime.Compiler,
	}
}

func (i Info) Fields() []interface{} {
	return []interface{}{
		"version", i.Version,
		"revision", i.Revision,
		"branch", i.Branch,
		"build_user", i.BuildUser,
		"build_date", i.BuildDate,
		"go_version", i.GoVersion,
		"os", i.Os,
		"arch", i.Arch,
		"compiler", i.Compiler,
	}
}
