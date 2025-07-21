# Contributing to Cursor Admin API Exporter

We welcome contributions to the Cursor Admin API Exporter! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker (for container builds)
- Helm 3.x (for Kubernetes deployment)
- Make (for build automation)
- Git

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/matanbaruch/cursor-admin-api-exporter.git
   cd cursor-admin-api-exporter
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up your environment:
   ```bash
   cp env.example .env
   # Edit .env with your Cursor API token
   ```

5. Run tests to ensure everything works:
   ```bash
   make test
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards
3. Add or update tests as needed
4. Run the test suite:
   ```bash
   make test
   ```

5. Run linters:
   ```bash
   make lint
   ```

6. Format your code:
   ```bash
   make fmt
   ```

### Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation changes
- `style:` formatting changes
- `refactor:` code refactoring
- `test:` adding or modifying tests
- `chore:` maintenance tasks

Examples:
```
feat: add support for custom API endpoints
fix: resolve memory leak in metrics collection
docs: update installation instructions
```

### Pull Request Process

1. Push your changes to your fork
2. Create a pull request against the main branch
3. Fill out the pull request template
4. Ensure all CI checks pass
5. Request review from maintainers

## Code Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Testing

- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Mock external dependencies

Example test structure:
```go
func TestNewCursorClient(t *testing.T) {
    tests := []struct {
        name     string
        baseURL  string
        token    string
        expected *CursorClient
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Documentation

- Update README.md for user-facing changes
- Add or update inline code comments
- Update API documentation
- Include examples for new features

## Types of Contributions

### Bug Reports

When reporting bugs, please include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Feature Requests

For new features, please provide:
- Clear description of the feature
- Use case and motivation
- Proposed implementation approach
- Potential impact on existing functionality

### Code Contributions

We welcome contributions including:
- Bug fixes
- New features
- Performance improvements
- Documentation updates
- Test coverage improvements

## Development Guidelines

### Adding New Metrics

When adding new metrics:

1. Define the metric in the appropriate exporter file
2. Add the metric to the `Describe` method
3. Implement collection logic in the `Collect` method
4. Add comprehensive tests
5. Update documentation

### API Client Changes

When modifying the API client:

1. Update the client interface if needed
2. Add appropriate error handling
3. Include unit tests with mocked responses
4. Update integration tests if applicable

### Configuration Changes

When adding new configuration options:

1. Add environment variable support
2. Update the Helm chart values
3. Add validation logic
4. Update documentation
5. Maintain backward compatibility

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test types
make test-unit
make test-integration
make test-performance

# Generate coverage report
make coverage
```

### Test Categories

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test API interactions (requires API token)
- **Performance Tests**: Test performance characteristics
- **End-to-End Tests**: Test complete workflows

### Coverage Requirements

- New code should have at least 80% test coverage
- Critical paths should have 100% coverage
- Integration tests should cover happy paths
- Error handling should be thoroughly tested

## Release Process

Releases are automated through GitHub Actions:

1. Changes are merged to main branch
2. Version is automatically bumped
3. Docker images are built and pushed
4. Helm charts are packaged and published
5. GitHub release is created with binaries
6. Documentation is updated

## Getting Help

- Join our [GitHub Discussions](https://github.com/matanbaruch/cursor-admin-api-exporter/discussions)
- Check existing [Issues](https://github.com/matanbaruch/cursor-admin-api-exporter/issues)
- Review the [Documentation](docs/)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Documentation acknowledgments

Thank you for contributing to Cursor Admin API Exporter!