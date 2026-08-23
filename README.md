# GitHub Container Registry Credential Helper

A tool to manage auth with personal access tokens (PATs) for individuals
to authenticate to GitHub Container Registry and push images. Probably
not meant for production workloads, as Actions will generate a short-lived
token for each job.

Based on similar tools like
[docker/docker-credential-helpers](https://github.com/docker/docker-credential-helpers),
[awslabs/amazon-ecr-credential-helper](https://github.com/awslabs/amazon-ecr-credential-helper), and
[GoogleCloudPlatform/docker-credential-gcr](https://github.com/GoogleCloudPlatform/docker-credential-gcr).

# How it Works

Uses your GitHub username/password to create a personal access token. Uses the [`gh`](https://github.com/cli/cli)
utilities to ensure the PAT has `write:packages` access. If it does not, it will
run the web-auth flow from the CLI and overwrite the existing token.

Once that is finished, we just rely on the `gh` configuration files as the source
of truth for a username and token, and provide that to the docker daemon
when requested.

# Configuration

Add this as a `credsHelper` in your Docker CLI Config:

```json
// $HOME/.docker/config.json
{
  "credHelpers": {
    "ghcr.io": "ghcr-login"
  }
}
```

## Testing

Build like: `go build -o docker-credential-ghcr-login main.go`.

```bash
$ echo "https://ghcr.io" | ./docker-credential-ghcr-login get
{"ServerURL":"https://ghcr.io","Username":"xxx","Secret":"gho_xxx"}
$ ./docker-credential-ghcr-login list
{"github.com":"xxx"}
```
