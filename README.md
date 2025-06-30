# ai-rules-link

<!-- Badges -->
[![License: PolyForm Noncommercial](https://img.shields.io/badge/license-PolyForm%20Noncommercial-blue.svg)](https://polyformproject.org/licenses/noncommercial/1.0.0/)
[![Go Report Card](https://goreportcard.com/badge/github.com/wchavarria03/ai-rules-link)](https://goreportcard.com/report/github.com/wchavarria03/ai-rules-link)
[![Go Reference](https://pkg.go.dev/badge/github.com/wchavarria03/ai-rules-link.svg)](https://pkg.go.dev/github.com/wchavarria03/ai-rules-link)
[![GitHub release](https://img.shields.io/github/v/release/wchavarria03/ai-rules-link.svg?style=flat)](https://github.com/wchavarria03/ai-rules-link/releases)
[![CI](https://github.com/wchavarria03/ai-rules-link/actions/workflows/ci.yml/badge.svg)](https://github.com/wchavarria03/ai-rules-link/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/wchavarria03/ai-rules-link/branch/main/graph/badge.svg)](https://codecov.io/gh/wchavarria03/ai-rules-link)

This tool standardizes AI-assisted development by generating and symlinking context-aware rules for the Cursor IDE. It allows developers to maintain a centralized repository of coding standards, best practices, and project-specific guidelines, then dynamically link them into their development environment. The tool supports both individual rule symlinking and consolidation into single files, making it easy to customize AI assistance based on technology stacks, project requirements, and personal coding preferences.

## Quick Start

### Install/Build

```sh
go build -o ai-rules-link .
```

### Usage

```bash
ai-rules-link rules --rule=go --rule=python
```
- Symlinks selected rules into your project's `.cursor/rules/` directory (e.g., `gorules.mdc`, `pythonrules.mdc`).
- Add `--global` to create rules in your home directory.
- Add `--consolidate` to merge selected rules into a single file (`consolidatedrules.mdc`).
- - Add `--force` to always overwrite destination files with embedded rules, even if they have been modified by the user.

For more details and advanced usage, see [docs/USAGE.md](docs/USAGE.md).

### Listing Symlinks

```bash
ai-rules-link status
```
- Lists all symlinks in `.cursor/rules/` and their targets.

## Development

- Build: `go build -o ai-rules-link .`
- Test: `make test`
- Docker: `docker build -t ai-rules-link .`

For detailed development, build, and contribution instructions, see [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).

## Documentation

- [Usage & CLI Options](docs/USAGE.md)
- [Canonical Rules Location](docs/RULES.md)
- [Contributing & Development](docs/DEVELOPMENT.md)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release notes and version history.

## License

This project is licensed under the PolyForm Noncommercial License 1.0.0. See [LICENSE](LICENSE) for details.

## How Rules Are Found

When you run the CLI, it looks for rule files in this order:
1. `${XDG_CONFIG_HOME}/ai-rules` (if set and exists)
2. `~/ai-rules` (in your home directory)
3. If neither exists, the CLI uses its own embedded rules (no setup needed)

You do **not** need to copy rule files manually—just use the binary! If you want to override or customize rules, create one of the above directories and add your own rule files.