# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

The Circular Protocol team takes security bugs seriously. We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.

### How to Report a Security Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

**Email**: [dannydenovi29@gmail.com](mailto:dannydenovi29@gmail.com)

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

This information will help us triage your report more quickly.

### What to Expect

After you have submitted a vulnerability report, you can expect:

1. **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours.

2. **Investigation**: We will investigate the issue and determine its impact and severity.

3. **Updates**: We will keep you informed about our progress in addressing the vulnerability.

4. **Resolution**: Once the vulnerability has been fixed, we will:
   - Release a security patch
   - Publicly disclose the vulnerability (with credit to you, if desired)
   - Update this security policy if necessary

### Disclosure Policy

- We ask that you give us a reasonable amount of time to fix the vulnerability before any public disclosure.
- We will credit you in our release notes (unless you prefer to remain anonymous).
- We follow a coordinated disclosure model and will work with you to understand and address the issue.

## Security Best Practices

When using the Circular Protocol Go SDK, please follow these security best practices:

### 1. API Keys and Secrets

- **Never commit API keys or private keys** to version control
- Use environment variables or secure secret management systems
- Rotate API keys regularly
- Use different keys for development, staging, and production

```go
// Good: Load from environment
api := circular.NewClient(
    os.Getenv("CIRCULAR_NAG_URL"),
    os.Getenv("CIRCULAR_API_KEY"),
)

// Bad: Hardcoded credentials
api := circular.NewClient(
    "https://nag.circularlabs.io/NAG.php?cep=",
    "my-secret-key-123", // Don't do this!
)
```

### 2. Private Key Management

- **Never expose private keys** in logs, error messages, or client-side code
- Store private keys securely (hardware wallets, secure enclaves, or encrypted storage)
- Never transmit private keys over insecure channels
- Use the SDK's built-in cryptographic functions for signing

```go
// Good: Load from secure storage
privateKey := secureStorage.GetPrivateKey()
signature, err := client.SignMessage(message, privateKey)

// Bad: Hardcoded private key
signature, err := client.SignMessage(message, "0x1234567890abcdef...") // Don't do this!
```

### 3. Input Validation

- Always validate and sanitize user inputs before using them in transactions
- Verify addresses are properly formatted
- Check transaction amounts are within expected ranges
- Validate smart contract inputs

```go
// Good: Validate inputs
if !isValidAddress(address) {
    return fmt.Errorf("invalid address format")
}
if amount <= 0 {
    return fmt.Errorf("amount must be positive")
}

result, err := client.SendTransaction(ctx, txParams)
```

### 4. Error Handling

- Don't expose sensitive information in error messages
- Log errors securely without including private data
- Handle errors gracefully to avoid information leakage

```go
// Good: Generic error message
if err != nil {
    log.Printf("Transaction failed: %v", err)
    return fmt.Errorf("transaction failed")
}

// Bad: Exposing sensitive data
if err != nil {
    return fmt.Errorf("failed to sign with key %s: %v", privateKey, err) // Don't do this!
}
```

### 5. Network Security

- Always use HTTPS for API communications
- Verify SSL/TLS certificates
- Use timeouts and context cancellation to prevent hanging connections
- Consider rate limiting for production applications

```go
// Good: Use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.SendTransaction(ctx, txParams)
```

### 6. Dependency Management

- Keep the SDK and its dependencies up to date
- Regularly run `go get -u` to update dependencies
- Monitor for security advisories
- Use `go mod verify` to ensure dependency integrity

```bash
# Update dependencies
go get -u github.com/circular-protocol/circular-go

# Verify module integrity
go mod verify

# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

### 7. Transaction Security

- Always verify transaction details before signing
- Use nonces correctly to prevent replay attacks
- Implement transaction amount limits
- Double-check recipient addresses

```go
// Good: Verify before sending
func sendWithConfirmation(client *circular.Client, tx map[string]interface{}) error {
    // Verify transaction details
    if !confirmTransaction(tx) {
        return fmt.Errorf("transaction not confirmed")
    }

    ctx := context.Background()
    result, err := client.SendTransaction(ctx, tx)
    if err != nil {
        return err
    }

    // Wait for confirmation
    return waitForConfirmation(client, result)
}
```

## Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed. Updates will be announced through:

- GitHub Security Advisories
- Release notes
- Project README

To stay informed about security updates:

- Watch the repository for releases
- Subscribe to the project's GitHub notifications
- Follow the project's official channels

## Comments on this Policy

If you have suggestions on how this process could be improved, please submit a pull request or open an issue.

---

**Last Updated**: November 2024
