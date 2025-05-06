package starred

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redjax/go-mygithub/internal/constants"
	"github.com/redjax/go-mygithub/internal/db"
	"github.com/redjax/go-mygithub/internal/domain/Github"
)

// Add 'starred' command handler
type GithubStarredCmd struct {
	// Github PAT
	AccessToken string `required:"" short:"t" long:"access-token" help:"Github Personal Acccess Token (PAT)"`
	// When present, save response to JSON file
	SaveJson bool `short:"s" long:"save-json" help:"Save response content to a file"`
	// JSON file to save response to
	OutputFile string `short:"o" long:"output" help:"Output file name" default:"starred_repos.json"`
	// When present, save response to DB
	SaveDB bool `long:"save-db" help:"Save response content to a database"`
	// Time between requests
	RequestSleep int `long:"request-sleep" help:"Time between requests (seconds)" default:"2"`

	// Get method for retrieving starred repos from Github API
	Get GetStarredCmd `cmd:"" help:"Get starred repositories"`
}

// Create class struct to pass into function when 'starred get' is used
type GetStarredCmd struct{}

// Handle 'starred get' command
func (u *GetStarredCmd) Run(cli *GithubStarredCmd) error {
	fmt.Println("Getting starred Github repositories.")

	// // Array to store all paginated starred repo responses
	// var allRepos []Github.Repository

	// Get API token from CLI
	token := cli.AccessToken
	// Get "Accept: " header value
	acceptHeaderVal := constants.GH_API_ACCCEPT_HEADER
	// Get Github starred repos API URL
	url := constants.GH_STARRED_URL

	allRepos, err := fetchAllStarredRepos(token, cli.RequestSleep, acceptHeaderVal, url)
	if err != nil {
		return fmt.Errorf("error fetching starred repositories: %w", err)
	}

	fmt.Printf("Fetched %d starred repositories.\n", len(allRepos))

	if cli.SaveDB {
		// Initialize your DB connection here (adjust to your actual DB init function)
		dbConn, err := db.InitDB()
		if err != nil {
			return fmt.Errorf("error initializing database: %w", err)
		}

		// Save repositories to the database
		err = db.SaveRepositories(dbConn, allRepos)
		if err != nil {
			return fmt.Errorf("error saving repositories to database: %w", err)
		}

		fmt.Println("Repositories saved to database successfully.")
	}

	// Save to JSON file if requested
	if cli.SaveJson {
		if cli.OutputFile == "" {
			return fmt.Errorf("output file path must be specified with -o/--output when using --save")
		}
		dir := filepath.Dir(cli.OutputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
		jsonBytes, err := json.MarshalIndent(allRepos, "", "  ")
		if err != nil {
			return fmt.Errorf("error formatting JSON: %w", err)
		}
		if err := os.WriteFile(cli.OutputFile, jsonBytes, 0644); err != nil {
			return fmt.Errorf("error saving response content to file: %s", err)
		}
		fmt.Printf("Starred repositories saved to: %s\n", cli.OutputFile)
	}

	return nil
}

// Helper: parseNextURL extracts the next page URL from the Link header.
func parseNextURL(linkHeader string) string {
	// Example Link header:
	// <https://api.github.com/user/starred?page=2>; rel="next", <https://api.github.com/user/starred?page=34>; rel="last"

	// Get URL parts
	parts := strings.Split(linkHeader, ",")

	// Iterate over parts to find rel="next"
	for _, p := range parts {
		if strings.Contains(p, `rel="next"`) {
			start := strings.Index(p, "<")
			end := strings.Index(p, ">")

			// Split the rel link into a URL
			if start != -1 && end != -1 && end > start {
				return p[start+1 : end]
			}
		}
	}

	// No next page
	return ""
}

// Helper: fetchAllStarredRepos fetches all starred repositories from the API
func fetchAllStarredRepos(token string, requestSleep int, acceptHeaderVal string, url string) ([]Github.Repository, error) {
	// Initialize HTTP client
	client := &http.Client{Timeout: 15 * time.Second}
	// Initialize array to hold all starred repositories
	var allRepos []Github.Repository
	// Initialize HTTP request page count
	page := 1

	for {
		fmt.Printf("Fetching page %d: %s\n", page, url)

		// Build HTTP request
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
		// Check for non-200 response
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		// Extract response body
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		// Unmarshal JSON
		var repos []Github.Repository
		if err := json.Unmarshal(bodyBytes, &repos); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		// Append repositories to allRepos array
		allRepos = append(allRepos, repos...)

		// Check for next page
		linkHeader := resp.Header.Get("Link")
		nextURL := parseNextURL(linkHeader)
		if nextURL == "" {
			break
		}

		// Increment page counter
		page++
		// Wait between requests
		time.Sleep(time.Duration(requestSleep) * time.Second)
		// Set URL for next loop
		url = nextURL
	}

	return allRepos, nil
}
