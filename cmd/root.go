package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httpreq",
	Short: "Make HTTP requests to JSON APIs using YAML and JSON",
	Long: `Make HTTP requests to APIs that expect JSON inputs using a YAML
configuration. The payload should be written as a multiline string.`,
}

// Execute ...
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
