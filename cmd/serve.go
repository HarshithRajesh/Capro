package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		origin, err := cmd.Flags().GetString("origin")
		if err != nil {
			return err
		}

		fmt.Printf("Serving on :%d\n", port)
		fmt.Printf("Origin on :%s\n", origin)
		return nil
	},
}

func init() {
	serveCmd.Flags().Int("port", 8080, "port to listen on")
	serveCmd.Flags().String("origin", "https://api.github.com", "Origin Domain")
	rootCmd.AddCommand(serveCmd)
}
