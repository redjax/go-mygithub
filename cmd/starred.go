package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// "time"

	"github.com/redjax/go-mygithub/internal/cache"
	"github.com/redjax/go-mygithub/internal/constants"
	"github.com/redjax/go-mygithub/internal/db"
	"github.com/redjax/go-mygithub/internal/domain/Github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	saveJson      bool
	outputFile    string
	saveDB        bool
	requestSleep  int
	cacheDir      string
	cacheDuration int
)

var starredCmd = &cobra.Command{
	Use:   "starred",
	Short: "Operations on starred repositories",
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get starred repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load token from flag, env, or config
		token := loadGitHubToken()
		if token == "" {
			return fmt.Errorf("GitHub access token not provided (use --access-token, GITHUB_TOKEN env, or config file)")
		}

		acceptHeaderVal := constants.GH_API_ACCCEPT_HEADER
		url := constants.GH_STARRED_URL

		allRepos, err := fetchAllStarredRepos(token, requestSleep, acceptHeaderVal, url)
		if err != nil {
			return fmt.Errorf("error fetching starred repositories: %w", err)
		}

		fmt.Printf("Fetched %d starred repositories.\n", len(allRepos))

		if saveDB {
			dbConn, err := db.InitDB()
			if err != nil {
				return fmt.Errorf("error initializing database: %w", err)
			}
			err = db.SaveRepositories(dbConn, allRepos)
			if err != nil {
				return fmt.Errorf("error saving repositories to database: %w", err)
			}
			fmt.Println("Repositories saved to database successfully.")
		}

		if saveJson {
			if outputFile == "" {
				return fmt.Errorf("output file path must be specified with --output when using --save-json")
			}
			dir := filepath.Dir(outputFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			jsonBytes, err := json.MarshalIndent(allRepos, "", "  ")
			if err != nil {
				return fmt.Errorf("error formatting JSON: %w", err)
			}
			if err := os.WriteFile(outputFile, jsonBytes, 0644); err != nil {
				return fmt.Errorf("error saving response content to file: %s", err)
			}
			fmt.Printf("Starred repositories saved to: %s\n", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(starredCmd)
	starredCmd.AddCommand(getCmd)

	// These are regular (not persistent) flags for starred get
	getCmd.Flags().BoolVar(&saveJson, "save-json", false, "Save response content to a file")
	getCmd.Flags().StringVarP(&outputFile, "output", "o", "starred_repos.json", "Output file name")
	getCmd.Flags().BoolVar(&saveDB, "save-db", false, "Save response content to a database")
	getCmd.Flags().IntVar(&requestSleep, "request-sleep", 2, "Time between requests (seconds)")

	// Bind flags to viper
	viper.BindPFlag("save_json", getCmd.Flags().Lookup("save-json"))
	viper.BindPFlag("output_file", getCmd.Flags().Lookup("output"))
	viper.BindPFlag("save_db", getCmd.Flags().Lookup("save-db"))
	viper.BindPFlag("request_sleep", getCmd.Flags().Lookup("request-sleep"))

	// Cache flags
	getCmd.Flags().StringVar(&cacheDir, "cache-dir", ".httpcache", "Directory for HTTP cache storage")
	getCmd.Flags().IntVar(&cacheDuration, "cache-duration", 5, "HTTP cache duration in minutes (0 to disable)")
	viper.BindPFlag("cache_dir", getCmd.Flags().Lookup("cache-dir"))
	viper.BindPFlag("cache_duration", getCmd.Flags().Lookup("cache-duration"))
	viper.SetDefault("cache_dir", ".httpcache")
	viper.SetDefault("cache_duration", 5)
}

func loadGitHubToken() string {
	// Priority: flag > env > config
	token := viper.GetString("access_token")
	return token
}

func fetchAllStarredRepos(token string, requestSleep int, acceptHeaderVal string, url string) ([]Github.Repository, error) {
	client := cache.NewCachingClient(
		viper.GetString("cache_dir"),
		viper.GetInt("cache_duration"),
	)

	var allRepos []Github.Repository
	page := 1

	for {
		fmt.Printf("Fetching page %d: %s\n", page, url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		req.Header.Set("Accept", acceptHeaderVal)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		if resp.StatusCode == 403 {
			reset := resp.Header.Get("X-RateLimit-Reset")
			resp.Body.Close()
			return nil, fmt.Errorf("rate limited by GitHub API, reset at %s", reset)
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var repos []Github.Repository
		if err := json.Unmarshal(bodyBytes, &repos); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		allRepos = append(allRepos, repos...)

		fmt.Printf("  Got %d repos (total so far: %d)\n", len(repos), len(allRepos))

		linkHeader := resp.Header.Get("Link")
		nextURL := parseNextURL(linkHeader)
		if nextURL == "" {
			break
		}
		page++
		// time.Sleep(time.Duration(requestSleep) * time.Second)
		url = nextURL
	}
	return allRepos, nil
}

func parseNextURL(linkHeader string) string {
	parts := strings.Split(linkHeader, ",")
	for _, p := range parts {
		if strings.Contains(p, `rel="next"`) {
			start := strings.Index(p, "<")
			end := strings.Index(p, ">")
			if start != -1 && end != -1 && end > start {
				return p[start+1 : end]
			}
		}
	}
	return ""
}
