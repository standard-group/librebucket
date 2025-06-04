# Contributing to LibreBucket

Thank you for your interest in contributing to LibreBucket! We welcome contributions from the community to help improve this project.

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork** the repository on [0xacab.org](https://0xacab.org/volumetech/librebucket)
2. **Clone** your fork locally:

   ```bash
   git clone https://0xacab.org/techplayz32/librebucket.git
   cd librebucket
   ```

3. **Create a branch** for your changes:

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes** and commit them with a descriptive message
5. **Push** your changes to your fork
6. Open a **merge request** to the `no-masters` branch

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git 2.20 or later

### Building

```bash
go build -o librebucket cmd/librebucket/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Style

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go/) and for markdown, [markdownlint's rules](https://github.com/DavidAnson/markdownlint/blob/v0.38.0/doc/Rules.md). When are you making markdown files, please make sure that you using the [markdownlint](https://github.com/DavidAnson/markdownlint) and following the rules of it.

Before submitting your changes, please run:

```bash
gofmt -s -w .
```

## Pull Request Guidelines

- **Keep PRs focused** - Each PR should address a single issue or add a single feature
- **Write tests** - New features should include tests
- **Update documentation** - Update the README or other relevant documentation
- **Squash commits** - Keep your commit history clean by squashing related commits
- **Write good commit messages** - Follow the [Conventional Commits](https://www.conventionalcommits.org/), or preferably at least [How to Write a Git Commit Message](https://cbea.ms/git-commit/) specification

## Reporting Issues

When reporting issues, please include:

- A clear description of the issue
- Steps to reproduce the issue
- Expected vs actual behavior
- Any relevant error messages or logs
- Your environment (OS, Go version, etc.)

## Feature Requests

We welcome feature requests! Please open an issue to discuss your idea before implementing it.

## License

By contributing to LibreBucket, you agree that your contributions will be licensed under the [MIT License](LICENSE).
