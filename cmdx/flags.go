package cmdx

import "github.com/spf13/pflag"

// RegisterConfigFlag registers the --config / -c flag.
func RegisterConfigFlag(fs *pflag.FlagSet, defaultPath string) {
	fs.StringP("config", "c", defaultPath, "Path to config file")
}
