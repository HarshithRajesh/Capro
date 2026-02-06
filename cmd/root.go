// Package cmd for CLI
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var clearCache bool

type Cache struct {
	Status int               `json:"status"`
	Body   []byte            `json:"body"`
	Header map[string]string `json:"headers"`
}

var cache = make(map[string]Cache)

const cacheFile = "cache.json"

func getCache() {
	file, err := os.ReadFile(cacheFile)
	if err != nil {
		fmt.Println("The cache file doesnt exist")
		return
	}

	err = json.Unmarshal(file, &cache)
	if err != nil {
		fmt.Print(err)
	}
}

func updateCache(key string, entry Cache) {
	cache[key] = entry
	data, _ := json.MarshalIndent(cache, "", " ")
	err := os.WriteFile(cacheFile, data, 0o644)
	if err != nil {
		fmt.Println("Failed to Write to the cacheFile")
	}
}

var rootCmd = &cobra.Command{
	Use:   "capro",
	Short: "Caching Proxy",
	Long:  "Caching Proxy Server CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if clearCache {
			emptyData := []byte("{}")
			err := os.WriteFile("cache.json", emptyData, 0o644)
			if err != nil {
				fmt.Println("Error Clearing Cache")
				os.Exit(1)
			}

			fmt.Println("Clearing Cache")
			return
		}

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Printf("Using default port 8080\n")
			port = 8080
		}
		origin, err := cmd.Flags().GetString("origin")
		if err != nil {
			fmt.Print("No origin was specified")
		}
		if origin == "" {
			fmt.Println("origin is required")
			return
		}
		getCache()

		fmt.Printf("Starting server on %d at %s\n", port, origin)
	},
}

func init() {
	rootCmd.Flags().BoolVar(&clearCache, "clear-cache", false, "Clear the local Cache")
	rootCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	rootCmd.Flags().String("origin", "https://api.github.com", "Origin Domain")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while executing Capro '%s'\n", err)
		os.Exit(1)
	}
}
