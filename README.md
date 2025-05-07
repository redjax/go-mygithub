# MyGithub

Reimplementing my [Python `mygithub` repository](https://github.com/redjax/mygithub) in Go.

## Details

Use Github Personal Access Tokens (PAT) to query Github's API. Currently only retrievees a user's starred repositories and optionally saves them to a JSON file or SQLite database.

## Links

- [Github docs: fine-grained Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)
- [Github API docs: get user's starred repositories](https://docs.github.com/en/rest/activity/starring?apiVersion=2022-11-28#list-repositories-starred-by-the-authenticated-user)
- [Pyinstrument docs](https://pyinstrument.readthedocs.io/en/latest/index.html)
  - [`pyinstrument` CLI args](https://pyinstrument.readthedocs.io/en/latest/guide.html)
  - [`pyinstrument` profile a specific chunk of code with a context manager](https://pyinstrument.readthedocs.io/en/latest/guide.html#profile-a-specific-chunk-of-code)
