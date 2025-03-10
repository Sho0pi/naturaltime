# Contributing to NaturalTime 🕒

Thank you for your interest in contributing to NaturalTime! This document provides guidelines and instructions for contributing to the project.

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/your-username/naturaltime.git
   cd naturaltime
   ```
3. **Set up the development environment**
   ```bash
   # Install dependencies
   npm install

   # Build the JavaScript file
   make build
   ```

## Development Workflow

### JavaScript Changes

The core parsing functionality uses [chrono-node](https://github.com/wanasit/chrono). If you need to modify the JavaScript code:

1. Edit `naturaltime.js`
2. Build the JavaScript file:
   ```bash
   make build
   ```
3. Run tests to verify your changes:
   ```bash
   go test ./...
   ```

### Go Changes

1. Make your changes to the Go code
2. Format your code:
   ```bash
   go fmt ./...
   ```
3. Run the tests:
   ```bash
   go test ./...
   ```

## Pull Request Process

1. Create a branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Commit your changes with clear, descriptive commit messages
3. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
4. Create a pull request from your branch to the `main` branch of the original repository
5. In your pull request description, explain:
    - What changes you've made
    - Why you've made them
    - How they've been tested

## Code Style

Please follow these guidelines for code style:

### Go
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go fmt` and `go vet` before submitting code
- Add meaningful comments and documentation
- Write tests for new functionality

### JavaScript
- Follow modern JavaScript practices
- Keep the JavaScript code minimal and focused on the parsing functionality

## Testing

- Add tests for new functionality
- Ensure all existing tests pass
- Include examples of time expressions that your changes support

## Releasing

The project uses GitHub Actions for automated releases:

1. To create a new release, tag your commit with a version number and the `-pre` suffix:
   ```bash
   git tag v1.0.0-pre
   git push origin v1.0.0-pre
   ```
2. The GitHub Actions workflow will:
    - Build the JavaScript file
    - Create a new tag without the `-pre` suffix
    - Create a GitHub release

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).