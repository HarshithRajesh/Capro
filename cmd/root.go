// Package cmd for CLI
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Cache struct {
	Status int               `json:"status"`
	Body   []byte            `json:"body"`
	Header map[string]string `json:"headers"`
}

var (
	err        error
	cache      = make(map[string]Cache)
	port       int
	origin     string
	clearCache bool
)

const cacheFile = "cache.json"

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if entry, found := cache[path]; found {
		fmt.Printf("Cache-Hit :%s\n", path)
		w.Header().Set("X-Cache", "HIT")
		for k, v := range entry.Header {
			w.Header().Set(k, v)
		}
		w.WriteHeader(entry.Status)

		if _, err = w.Write(entry.Body); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Cache-MISS :%s\n", path)
	newPath := strings.TrimSuffix(origin, "/") + path
	req, err := http.Get(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	respCode := req.StatusCode
	if respCode != 200 {
		log.Fatal("invalid url")
		os.Exit(1)
	}

	originalHeaders := req.Header
	headerMap := make(map[string]string)
	for k := range originalHeaders {
		headerMap[k] = req.Header.Get(k)
	}
	bodyString := string(body)

	entry := Cache{
		Status: respCode,
		Body:   []byte(bodyString),
		Header: headerMap,
	}

	updateCache(path, entry)

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-type", entry.Header["Content-type"])
	w.WriteHeader(entry.Status)
	w.Write([]byte(bodyString))
}

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

		port, err = cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Printf("Using default port 8080\n")
			port = 8080
		}
		origin, err = cmd.Flags().GetString("origin")
		if err != nil {
			fmt.Print("No origin was specified")
		}
		if origin == "" {
			fmt.Println("origin is required")
			return
		}
		getCache()
		http.HandleFunc("/", proxyHandler)
		fmt.Printf("Starting Proxy Server on port %d\n", port)
		fmt.Printf("Forwarding requests to: %s\n", origin)
		portadr := fmt.Sprintf(":%d", port)
		err = http.ListenAndServe(portadr, nil)
		if err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
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
