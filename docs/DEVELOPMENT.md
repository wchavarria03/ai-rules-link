# Contributing & Development

## Local Build & Run

### Build and Run Locally (No Docker)

1. Build the binary:
   ```sh
   go build -o ai-rules-link .
   ```
2. Run the CLI in your target project:
   ```sh
   ./ai-rules-link rules --rule=go --rule=python --rule=personalcommits
   ```
   - This will operate on the current working directory (your project).

### Build and Run with Docker

1. Build the Docker image:
   ```sh
   docker build -t ai-rules-link .
   ```
2. Run the CLI in your target project using Docker (replace the path with your project):
   ```sh
   docker run --rm -it \
     -v /path/to/your/project:/work \
     -v $HOME/.sync-rules/rules:/root/.sync-rules/rules \
     -w /work \
     ai-rules-link rules --rule=go --rule=python --rule=personalcommits
   ```
   - This mounts your project and rules directory into the container and runs the CLI inside the container.

## Testing

- Run all tests:
  ```sh
  make test
  ```

## Rules Lookup Order

The CLI looks for rules in this order:
1. `${XDG_CONFIG_HOME}/ai-rules` (if set and exists)
2. `~/ai-rules` (in your home directory)
3. If neither exists, the CLI uses its own embedded rules (no setup needed)

This means you do not need to copy rule files for development or testing unless you want to test overrides. Rule files should be named without underscores, e.g., `gorules.mdc`, `pythonrules.mdc`, `personalcommitsrules.mdc`.

## Contributing

- Please follow the project's code style and best practices.
- See [ARCHITECTURE.md](../ARCHITECTURE.md) for design and structure guidelines.
- Pull requests and issues are welcome! 