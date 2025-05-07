package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redjax/go-mygithub/internal/cache"
	"github.com/redjax/go-mygithub/internal/constants"
	"github.com/redjax/go-mygithub/internal/db"
	"github.com/redjax/go-mygithub/internal/domain/Github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cobra flags
var (
	saveJson      bool
	outputFile    string
	saveDB        bool
	requestSleep  int
	cacheDir      string
	cacheDuration int
)

// Init "starred" subcommand
var starredCmd = &cobra.Command{
	Use:   "starred",
	Short: "Operations on starred repositories",
}

// Init "starred get" subcommand
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get starred repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load token from flag, env, or config
		token := loadGithubToken()
		if token == "" {
			return fmt.Errorf("GitHub access token not provided (use --access-token, GITHUB_TOKEN env, or config file)")
		}

		// Set "Accept:" header
		acceptHeaderVal := constants.GH_API_ACCCEPT_HEADER
		// Set Github API stars URL
		url := constants.GH_STARRED_URL

		// Make HTTP requests to fetch user's starred repositories
		allRepos, err := fetchAllStarredRepos(token, requestSleep, acceptHeaderVal, url)
		if err != nil {
			return fmt.Errorf("error fetching starred repositories: %w", err)
		}

		if len(allRepos) == 0 {
			return fmt.Errorf("no starred repositories returned for this PAT")
		}

		fmt.Printf("Fetched %d starred repositories.\n", len(allRepos))

		if saveDB {
			// Save fetched repositories to database

			// Initialize database
			dbConn, err := db.InitDB()
			if err != nil {
				return fmt.Errorf("error initializing database: %w", err)
			}

			// Save retrieved repositories
			err = db.SaveRepositories(dbConn, allRepos)
			if err != nil {
				return fmt.Errorf("error saving repositories to database: %w", err)
			}
			fmt.Println("Repositories saved to database successfully.")
		}

		if saveJson {
			// Save repositories to JSON

			// Validate outputFile
			if outputFile == "" {
				return fmt.Errorf("output file path must be specified with --output when using --save-json")
			}

			// Ensure file's parent dir exists
			dir := filepath.Dir(outputFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}

			// Marshal response content into JSON
			jsonBytes, err := json.MarshalIndent(allRepos, "", "  ")
			if err != nil {
				return fmt.Errorf("error formatting JSON: %w", err)
			}

			// Write to file
			if err := os.WriteFile(outputFile, jsonBytes, 0644); err != nil {
				return fmt.Errorf("error saving response content to file: %s", err)
			}
			fmt.Printf("Starred repositories saved to: %s\n", outputFile)
		}

		return nil
	},
}

// "starred" CLI entrypoint
func init() {

	// Add starred subcommand to root CLI
	rootCmd.AddCommand(starredCmd)
	// Add "get" subcommand to starred subcommand
	starredCmd.AddCommand(getCmd)

	// Save to JSON file
	getCmd.Flags().BoolVar(&saveJson, "save-json", false, "Save response content to a file")
	// File to save output to
	getCmd.Flags().StringVarP(&outputFile, "output", "o", "starred_repos.json", "Output file name")
	// Save to database
	getCmd.Flags().BoolVar(&saveDB, "save-db", false, "Save response content to a database")
	// Time between requests
	getCmd.Flags().IntVar(&requestSleep, "request-sleep", 0, "Time between requests (seconds)")

	// HTTP cache control flags
	getCmd.Flags().StringVar(&cacheDir, "cache-dir", ".httpcache", "Directory for HTTP cache storage")
	getCmd.Flags().IntVar(&cacheDuration, "cache-duration", 5, "HTTP cache duration in minutes (0 to disable)")

	// Bind flags to viper
	viper.BindPFlag("save_json", getCmd.Flags().Lookup("save-json"))
	viper.BindPFlag("output_file", getCmd.Flags().Lookup("output"))
	viper.BindPFlag("save_db", getCmd.Flags().Lookup("save-db"))
	viper.BindPFlag("request_sleep", getCmd.Flags().Lookup("request-sleep"))
	viper.BindPFlag("cache_dir", getCmd.Flags().Lookup("cache-dir"))
	viper.BindPFlag("cache_duration", getCmd.Flags().Lookup("cache-duration"))
	viper.SetDefault("cache_dir", ".httpcache")
	viper.SetDefault("cache_duration", 5)
}

// Load Github PAT from viper
func loadGithubToken() string {
	// Priority: flag > env > config
	token := viper.GetString("access_token")

	return token
}

// Fetch all starred repositories from Github
func fetchAllStarredRepos(token string, requestSleep int, acceptHeaderVal string, url string) ([]Github.Repository, error) {
	// Get HTTP cache client
	client := cache.NewCachingClient(
		viper.GetString("cache_dir"),
		viper.GetInt("cache_duration"),
	)

	// Initialize array for storing repository schemas
	var allRepos []Github.Repository
	// Initialize page count
	page := 1

	for {
		fmt.Printf("Fetching page %d: %s\n", page, url)

		// Build request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		// Set request headers
		req.Header.Set("Accept", acceptHeaderVal)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Make HTTP request
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		// Check for rate limiting
		if resp.StatusCode == 403 {
			reset := resp.Header.Get("X-RateLimit-Reset")
			resp.Body.Close()
			return nil, fmt.Errorf("rate limited by GitHub API, reset at %s", reset)
		}
		// Check for unexpected status
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		// Read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		// Unmarshal this page of repositories into Repository schemas
		var repos []Github.Repository
		if err := json.Unmarshal(bodyBytes, &repos); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		// Append loop's Repository schemas to allRepos array
		allRepos = append(allRepos, repos...)

		fmt.Printf("  Got %d repos (total so far: %d)\n", len(repos), len(allRepos))

		// Check for next page
		linkHeader := resp.Header.Get("Link")
		nextURL := parseNextURL(linkHeader)
		if nextURL == "" {
			// No more pages
			break
		}

		// Increment page count
		page++

		// Wait before next request
		time.Sleep(time.Duration(requestSleep) * time.Second)
		// Set URL for next loop
		url = nextURL
	}

	return allRepos, nil
}

// Extraxt next URL from response
func parseNextURL(linkHeader string) string {
	// Parse Link header
	parts := strings.Split(linkHeader, ",")

	for _, p := range parts {
		// Find rel="next"
		if strings.Contains(p, `rel="next"`) {
			// Extract URL
			start := strings.Index(p, "<")
			end := strings.Index(p, ">")
			if start != -1 && end != -1 && end > start {
				return p[start+1 : end]
			}
		}
	}

	return ""
}
