# Contributing to Circular Protocol Go SDK

Thank you for your interest in contributing to the Circular Protocol Go SDK! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Git
- Basic understanding of the Circular Protocol blockchain

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/circular-go.git
   cd circular-go
   ```

3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/circular-protocol/circular-go.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Verify everything works:
   ```bash
   go test ./...
   ```

## Development Workflow

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** following the coding standards

3. **Write tests** for your changes

4. **Run tests** to ensure everything passes:
   ```bash
   go test ./... -v
   go test -race ./...
   go test -cover ./...
   ```

5. **Update documentation** as needed

6. **Commit your changes** using conventional commits

7. **Push to your fork**:
   ```bash
   git push origin feat/your-feature-name
   ```

8. **Open a Pull Request** against the `main` branch

## Coding Standards

### Go Style Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Use `golint` and `go vet` to catch common issues
- Follow the existing code style in the repository

### Code Organization

- Keep functions focused and small
- Use descriptive variable and function names
- Add comments for exported functions (godoc format)
- Group related functionality together

### Dual Signature Pattern

All API methods must support two calling patterns:

1. **Raw method** (request object style - explicit control):
```go
// CheckWalletRaw checks if wallet exists (request object style)
func (c *Client) CheckWalletRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
    return c.makeRequest(ctx, "CheckWallet", req)
}
```

2. **Convenience method** (positional parameters with auto-preprocessing):
```go
// CheckWallet checks if wallet exists (convenience method)
// Auto-strips '0x' prefix, auto-injects version
func (c *Client) CheckWallet(ctx context.Context, blockchain string, address string) (map[string]interface{}, error) {
    return c.CheckWalletRaw(ctx, map[string]interface{}{
        "Blockchain": c.HexFix(blockchain),
        "Address":    c.HexFix(address),
        "Version":    c.version,
    })
}
```

**Auto-preprocessing rules:**
- Use `c.HexFix()` for blockchain, address, ID, nodeID parameters
- Use `c.StringToHex()` for project, request parameters in contract methods
- Always inject `Version: c.version`
- Auto-generate timestamps for contract methods with `c.GetFormattedTimestamp()`

### Example

```go
// SendTransactionRaw submits transaction (request object style)
func (c *Client) SendTransactionRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
    return c.makeRequest(ctx, "AddTransaction", req)
}

// SendTransaction submits transaction with positional parameters (convenience method)
// Matches JavaScript implementation signature, auto-injects version
func (c *Client) SendTransaction(
    ctx context.Context,
    id string,
    from string,
    to string,
    timestamp string,
    txType string,
    payload string,
    nonce string,
    signature string,
    blockchain string,
) (map[string]interface{}, error) {
    return c.SendTransactionRaw(ctx, map[string]interface{}{
        "ID":         id,
        "From":       from,
        "To":         to,
        "Timestamp":  timestamp,
        "Type":       txType,
        "Payload":    payload,
        "Nonce":      nonce,
        "Signature":  signature,
        "Blockchain": blockchain,
        "Version":    c.version,
    })
}
```

## Testing

### Test Requirements

- All new features must include tests
- Aim for high test coverage (>80%)
- Include unit tests, integration tests, and e2e tests where appropriate

### Test Types

1. **Unit Tests** - Test individual functions
   ```bash
   go test ./... -short
   ```

2. **Integration Tests** - Test API interactions with mock servers
   ```bash
   go test ./... -run Integration
   ```

3. **E2E Tests** - Test against real endpoints (use with caution)
   ```bash
   go test ./... -run E2E
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestCheckWallet
```

## Submitting Changes

### Pull Request Process

1. **Ensure all tests pass** before submitting

2. **Update the README.md** if you're adding new features

3. **Update CHANGELOG.md** with your changes following [Keep a Changelog](https://keepachangelog.com/) format

4. **Ensure 100% backward compatibility** unless this is a major version bump

5. **Write a comprehensive PR description** including:
   - Overview of changes
   - Why the changes were made
   - Examples of usage (before/after if applicable)
   - Migration guide (if breaking changes)
   - Checklist of completed items

### PR Title Format

Use conventional commits format:

```
<type>: <description> (vX.Y.Z -> vA.B.C)
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `perf:` - Performance improvements

**Examples:**
- `feat: add GetTransactionOutcome polling method (v1.0.8 -> v1.0.9)`
- `fix: correct method naming to match API specification (v1.0.9 -> v1.0.10)`
- `docs: add comprehensive API reference to README`

### PR Checklist

Your PR should include:

- [ ] All tests pass
- [ ] New tests added for new features
- [ ] Documentation updated (README, godoc comments)
- [ ] CHANGELOG.md updated
- [ ] Code follows Go best practices
- [ ] No breaking changes (or clearly documented if unavoidable)
- [ ] Commit messages follow conventional commits
- [ ] Examples added for new features

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Examples

```
feat(client): add context support to all API methods

Add context.Context parameter to all client methods for
better cancellation and timeout control.

BREAKING CHANGE: All client methods now require a context parameter
```

```
fix(wallet): correct RegisterWallet to use client methods

RegisterWallet was calling undefined package-level functions.
Updated to use client methods (c.HashString, c.StringToHex, etc.)

Fixes #123
```

## Documentation

### README Updates

When adding new features, update the README.md with:
- Feature description
- Usage examples
- API reference additions

### Code Comments

- All exported functions must have godoc comments
- Complex logic should include inline comments
- Use complete sentences for comments
- Start comments with the function/type name

### Example Documentation

```go
// Client represents a Circular Protocol API client.
// It manages connections to the Circular blockchain network
// and provides methods for interacting with wallets, transactions,
// blocks, and smart contracts.
type Client struct {
    nagURL     string
    nagKey     string
    httpClient *http.Client
    headers    map[string]string
}
```

## Questions or Need Help?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Join the community discussions
- Refer to the [official documentation](https://circular-protocol.gitbook.io/circular-sdk/api-docs/go)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to the Circular Protocol Go SDK! 🎉
