package cmdx

import flag "github.com/spf13/pflag"

// RegisterConfigFlag registers the --config / -c flag.
func RegisterConfigFlag(fs *flag.FlagSet, defaultPath string) {
	fs.StringP("config", "c", defaultPath, "Path to config file")
}

// RegisterDebugFlag registers the --debug / -d flag.
func RegisterDebugFlag(fs *flag.FlagSet) {
	fs.BoolP("debug", "d", false, "Enable debug mode")
	_ = fs.MarkHidden("debug")
}
