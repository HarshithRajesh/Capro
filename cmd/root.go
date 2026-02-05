// Package cmd for CLI
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "capro",
	Short: "Caching Proxy",
	Long:  "Caching Proxy Server CLI",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while executing Capro '%s'\n", err)
		os.Exit(1)
	}
}
