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

func loadCache() {
	file, err := os.ReadFile(cacheFile)
	if err == nil {
		json.Unmarshal(file, &cache)
	}
}

func updateCache(path string, entry Cache) {
	cache[path] = entry
	data, _ := json.MarshalIndent(cache, "", "  ")
	os.WriteFile(cacheFile, data, 0o644)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	newPath := strings.TrimSuffix(origin, "/") + path
	if entry, found := cache[newPath]; found {
		fmt.Printf("Cache-Hit :%s\n", path)
		w.Header().Set("X-Cache", "HIT")
		for k, v := range entry.Header {
			w.Header().Set(k, v)
		}
		w.WriteHeader(entry.Status)

		if _, err = w.Write(entry.Body); err != nil {
			http.Error(w, "Message", http.StatusBadGateway)
			return
		}
		return
	}

	fmt.Printf("Cache-MISS :%s\n", path)
	resp, err := http.Get(newPath)
	if err != nil {
		http.Error(w, "Message", http.StatusBadGateway)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Message", http.StatusBadGateway)
	}
	respCode := resp.StatusCode
	if respCode != 200 {
		log.Fatal("invalid url")
		return
	}

	originalHeaders := resp.Header
	headerMap := make(map[string]string)
	for k := range originalHeaders {
		headerMap[k] = resp.Header.Get(k)
	}
	bodyString := string(body)

	entry := Cache{
		Status: respCode,
		Body:   []byte(bodyString),
		Header: headerMap,
	}

	updateCache(newPath, entry)

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", entry.Header["Content-Type"])
	w.WriteHeader(entry.Status)
	w.Write([]byte(bodyString))
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
				return
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
			fmt.Println("origin is respuired")
			return
		}
		loadCache()
		http.HandleFunc("/", proxyHandler)
		fmt.Printf("Starting Proxy Server on port %d\n", port)
		fmt.Printf("Forwarding respuests to: %s\n", origin)
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
		return
	}
}

//Add Mutex for the concurrent respuests
//
