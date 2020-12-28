package cmdx

import "github.com/spf13/pflag"

// RegisterConfigFlag registers the --config / -c flag.
func RegisterConfigFlag(fs *pflag.FlagSet, defaultPath string) {
	fs.StringP("config", "c", defaultPath, "Path to config file")
}

// RegisterDebugFlag registers the --debug / -d flag.
func RegisterDebugFlag(fs *pflag.FlagSet) {
	fs.Bool("debug", false, "Enable debug mode")
}
