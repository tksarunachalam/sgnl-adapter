# SGNL Ingestion Adapter Template

1. Clone this repository.
1. Update the names of `github.com/sgnl-ai/adapter-template/*` Golang packages in all files to match your new repository's name (e.g. `github.com/your-org/your-repo`):

   ```
   sed -e 's,^module github\.com/sgnl-ai/adapter-template,github.com/your-org/your-repo,' -i go.mod
   ```

   ```
   find pkg/ -type f -name '*.go' | xargs -n 1 sed -n -e 's,github\.com/sgnl-ai/adapter-template,github.com/your-org/your-repo,p' -i
   ```

1. Modify the adapter implementation in package `pkg/adapter` to query your datasource. All the code that must be modified is identified with `SCAFFOLDING` comments.
1. Delete the example datasource implementation in package `pkg/example_datasource`.
1. Build the Docker image with the `adapter` command.
   ```
   docker build -t adapter:latest .
   ```
1. **WARNING: Temporary workaround to allow `go` to download the adapter-framework module before it is made public.**
   **The image's history will contain your environment variables, incl. any GitHub credentials that it may contain.**
   ```
   docker build -t adapter:latest --build-arg GITHUB_USER="$GITHUB_USER" --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" .
   ```
1. Run the adapter server as a Docker container.
   ```
   docker run --rm -it adapter:latest
   ```
