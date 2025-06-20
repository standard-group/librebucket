<div style="background: #f0f0f0; padding: 20px; display: inline-block;">
   <img src="./docs/images/new-librebucket-logo.svg">
</div>

# Librebucket

[![Go Report Card](https://goreportcard.com/badge/github.com/standard-group/librebucket)](https://goreportcard.com/report/github.com/standard-group/librebucket) 
[![License: GPLv3 or later](https://img.shields.io/badge/License-GPLv3-yellow.svg)](https://opensource.org/license/gpl-3-0)
![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/standard-group/librebucket?utm_source=oss&utm_medium=github&utm_campaign=standard-group%2Flibrebucket&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)

A lightweight, self-hosted Git server with a clean web interface, built with Go.

Are you tired of big companies playing with your repositories, by removing them or trading with your data? Want to provide a small and free Git hosting for your community, friends or company? Then Librebucket (*/ˈliːbrə ˈbʌkɪt/*) for you, we promise that Librebucket will be **"independent Free/Libre Software forever"**!

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git 2.20 or later
- A operating system (Linux, macOS, Windows, even Raspberry/Orange Pi!)

### Installation

TODO: add instructions for installation

### Creating Your First Repository

1. Use the API to create a new repository:

   ```bash
   curl -X POST http://localhost:3000/api/v1/git/create \
     -H "Content-Type: application/json" \
     -d '{"username":"yourusername", "reponame":"yourrepo"}'
   ```

2. Clone your new repository:

   ```bash
   git clone http://localhost:3000/yourusername/yourrepo.git
   cd yourrepo
   ```

3. Make some changes and push:

   ```bash
   echo "# My Project" > README.md
   git add .
   git commit -m "Initial commit"
   git push -u origin main
   ```

### Web Interface

Access the web interface at `http://localhost:3000` to browse repositories, view code, and manage your projects.

## Configuration

TODO: do the configuration

## Development

### Building

```bash
go build -o librebucket cmd/librebucket/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Style

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go/).

## Security

- All repository access is private by default
- HTTPS is recommended for secure communication
- Authentication and authorization features coming soon

## Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) before submitting pull requests.

## Repo Activity

![Alt](https://repobeats.axiom.co/api/embed/8fb7eb654c3100485ab3adef9d5bd7e543568d94.svg "Repobeats analytics image")

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with the amazing Go programming language
- Inspired by [Forgejo](https://codeberg.org/forgejo/forgejo) and other self-hosted Git solutions
- Uses the [go-git](https://github.com/go-git/go-git) library for Git operations
