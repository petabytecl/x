package cmdx

import "github.com/spf13/cobra"

// RegisterConfigFlag registers the --config / -c flag.
func RegisterConfigFlag(c *cobra.Command, defaultPath string) {
	c.PersistentFlags().StringP("config", "c", defaultPath, "Path to config file")
}

// RegisterDebugFlag registers the --debug / -d flag.
func RegisterDebugFlag(c *cobra.Command) {
	c.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
	_ = c.PersistentFlags().MarkHidden("debug")
}
