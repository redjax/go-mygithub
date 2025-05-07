# MyGithub

Reimplementing my [Python `mygithub` repository](https://github.com/redjax/mygithub) in Go.

## Details

Use Github Personal Access Tokens (PAT) to query Github's API. Currently only retrievees a user's starred repositories and optionally saves them to a JSON file or SQLite database.

## Usage

*Note: These examples may not be complete. Run `mygithub --help` to see all options & usage.*

```bash
$ mygithub --help

Usage:
  mygithub [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  starred     Operations on starred repositories

Flags:
  -t, --access-token string   GitHub Personal Access Token (PAT)
  -h, --help                  help for mygithub

Use "mygithub [command] --help" for more information about a command.
```

```bash
$ mygithub starred --help

Operations on starred repositories

Usage:
  mygithub starred [command]

Available Commands:
  get         Get starred repositories

Flags:
  -h, --help   help for starred

Global Flags:
  -t, --access-token string   GitHub Personal Access Token (PAT)

Use "mygithub starred [command] --help" for more information about a command.
```

```bash
$ mygithub starred get --help

Get starred repositories

Usage:
  mygithub starred get [flags]

Flags:
      --cache-dir string     Directory for HTTP cache storage (default ".httpcache")
      --cache-duration int   HTTP cache duration in minutes (0 to disable) (default 5)
  -h, --help                 help for get
  -o, --output string        Output file name (default "starred_repos.json")
      --request-sleep int    Time between requests (seconds)
      --save-db              Save response content to a database
      --save-json            Save response content to a file

Global Flags:
  -t, --access-token string   GitHub Personal Access Token (PAT)
```
## Links

- [Github docs: fine-grained Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)
- [Github API docs: get user's starred repositories](https://docs.github.com/en/rest/activity/starring?apiVersion=2022-11-28#list-repositories-starred-by-the-authenticated-user)
- [Pyinstrument docs](https://pyinstrument.readthedocs.io/en/latest/index.html)
  - [`pyinstrument` CLI args](https://pyinstrument.readthedocs.io/en/latest/guide.html)
  - [`pyinstrument` profile a specific chunk of code with a context manager](https://pyinstrument.readthedocs.io/en/latest/guide.html#profile-a-specific-chunk-of-code)
