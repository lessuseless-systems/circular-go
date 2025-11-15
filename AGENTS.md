# Circular Protocol Go SDK - Architecture Guide

This document provides a comprehensive overview of the Circular Protocol Go SDK architecture, design decisions, and implementation details.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Package Structure](#package-structure)
- [Design Patterns](#design-patterns)
- [API Client](#api-client)
- [Error Handling](#error-handling)
- [Cryptography](#cryptography)
- [Testing Strategy](#testing-strategy)
- [Migration Guide](#migration-guide)

## Overview

The Circular Protocol Go SDK provides two implementations for interacting with the Circular blockchain:

1. **Modern SDK** (`circular_protocol.go`) - Context-aware, idiomatic Go implementation
2. **Legacy API** (`circular_protocol_api/`) - Original package-level implementation

Both implementations are maintained for backward compatibility, but new projects should use the modern SDK.

### Key Features

- **Context Support**: Full `context.Context` integration for cancellation and timeouts
- **Type Safety**: Proper Go types and error handling
- **Client Pattern**: Object-oriented design with `Client` struct
- **Comprehensive Testing**: Unit, integration, and e2e test coverage
- **Cryptographic Functions**: Built-in secp256k1 signing and verification
- **Flexible Configuration**: Customizable HTTP client, headers, and timeouts

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  (Your Go Application using Circular Protocol SDK)          │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ├──────────────────┬───────────────────┐
                     │                  │                   │
         ┌───────────▼──────────┐  ┌────▼────┐  ┌──────────▼─────────┐
         │   Modern SDK         │  │  Utils  │  │   Legacy API       │
         │ (circular_protocol)  │  │         │  │ (circular_protocol │
         │                      │  │         │  │  _api)             │
         └───────────┬──────────┘  └────┬────┘  └──────────┬─────────┘
                     │                  │                   │
                     │                  │                   │
         ┌───────────▼──────────────────▼───────────────────▼─────────┐
         │                  HTTP Client Layer                          │
         │         (Network requests to Circular NAG API)             │
         └────────────────────────────┬───────────────────────────────┘
                                      │
                          ┌───────────▼──────────┐
                          │  Circular Blockchain  │
                          │   NAG API Endpoint    │
                          └──────────────────────┘
```

### Request Flow

```
1. Application calls SDK method
   ↓
2. SDK validates input and constructs request
   ↓
3. Request is serialized to JSON
   ↓
4. HTTP POST request sent to NAG endpoint
   ↓
5. Response received and parsed
   ↓
6. SDK validates response structure
   ↓
7. Result returned to application
```

## Package Structure

```
circular-go/
├── circular_protocol.go          # Modern SDK implementation
├── circular_protocol_test.go     # Unit tests
├── circular_protocol_integration_test.go  # Integration tests
├── circular_protocol_e2e_test.go # End-to-end tests
│
├── circular_protocol_api/        # Legacy implementation
│   ├── init.go
│   └── circular_protocol_api.go
│
├── utils/                        # Shared utilities
│   ├── init.go
│   └── utils.go                  # Crypto & helper functions
│
├── .github/
│   └── workflows/
│       └── test.yml              # CI/CD pipeline
│
├── README.md                     # User documentation
├── CONTRIBUTING.md               # Contribution guidelines
├── CODE_OF_CONDUCT.md            # Community guidelines
├── SECURITY.md                   # Security policy
├── CHANGELOG.md                  # Version history
├── AGENTS.md                     # This file
├── LICENSE                       # MIT License
├── go.mod                        # Go modules
└── go.sum                        # Dependency checksums
```

## Design Patterns

### 1. Client Pattern

The modern SDK uses the client pattern for better encapsulation and state management:

```go
type Client struct {
    nagURL     string
    nagKey     string
    httpClient *http.Client
    headers    map[string]string
    lastError  string
}

// Constructor with default configuration
func NewClient(nagURL, nagKey string) *Client

// Constructor with custom configuration
func NewClientWithConfig(cfg Config) *Client
```

**Benefits**:
- Encapsulates configuration and state
- Allows multiple client instances
- Supports custom HTTP clients and headers
- Testable with mock HTTP clients

### 2. Functional Options Pattern

Configuration uses functional options for flexibility:

```go
type Config struct {
    NAGURL     string
    NAGKey     string
    Timeout    time.Duration
    HTTPClient *http.Client
    Headers    map[string]string
}

client := NewClientWithConfig(Config{
    NAGURL:  "https://custom-nag.example.com",
    NAGKey:  apiKey,
    Timeout: 60 * time.Second,
    Headers: map[string]string{
        "X-Custom-Header": "value",
    },
})
```

### 3. Context Pattern

All API methods accept `context.Context` for cancellation and timeouts:

```go
func (c *Client) CheckWallet(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error)
```

**Benefits**:
- Request cancellation
- Timeout control
- Request-scoped values
- Better resource management

### 4. Builder Pattern

Transaction construction uses a builder-like pattern:

```go
request := map[string]interface{}{
    "ID":         id,
    "From":       from,
    "To":         to,
    "Timestamp":  timestamp,
    "Type":       txType,
    "Payload":    payload,
    "Nonce":      nonce,
    "Signature":  signature,
    "Blockchain": blockchain,
    "Version":    "1.0.8",
}
```

## API Client

### Client Lifecycle

```go
// 1. Create client
client := circular.NewClient(nagURL, apiKey)

// 2. Use client for multiple requests
ctx := context.Background()

// 3. Make API calls
wallet, err := client.GetWallet(ctx, req)
balance, err := client.GetWalletBalance(ctx, req)
tx, err := client.SendTransaction(ctx, req)

// 4. No explicit cleanup needed (HTTP client reused)
```

### Request Handling

The `makeRequest` method centralizes all HTTP communication:

```go
func (c *Client) makeRequest(ctx context.Context, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
    // 1. Construct URL
    url := c.nagURL + "Circular_" + endpoint + "_"

    // 2. Marshal JSON
    jsonData, err := json.Marshal(data)

    // 3. Create HTTP request with context
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))

    // 4. Set headers
    req.Header.Set("Content-Type", "application/json")
    if c.nagKey != "" {
        req.Header.Set("X-NAG-Key", c.nagKey)
    }

    // 5. Send request
    resp, err := c.httpClient.Do(req)

    // 6. Parse and validate response
    // ...
}
```

### API Method Pattern

All API methods follow a consistent pattern:

```go
// MethodName Brief description
// Longer description with usage details
func (c *Client) MethodName(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
    return c.makeRequest(ctx, "MethodName", req)
}
```

## Error Handling

### Error Types

```go
// APIError represents an error from the Circular Protocol API
type APIError struct {
    Message    string
    StatusCode int
    Endpoint   string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error on %s (status %d): %s",
        e.Endpoint, e.StatusCode, e.Message)
}
```

### Error Handling Patterns

```go
// Check for specific error types
result, err := client.CheckWallet(ctx, req)
if err != nil {
    if apiErr, ok := err.(*circular.APIError); ok {
        // Handle API-specific error
        log.Printf("API error on %s: %s", apiErr.Endpoint, apiErr.Message)
    } else {
        // Handle other errors (network, timeout, etc.)
        log.Printf("Request failed: %v", err)
    }
    return err
}
```

### Response Validation

```go
// API responses are validated for:
// 1. HTTP status code (must be 200)
// 2. Result field (must be 200)
// 3. Response field presence

if resp.StatusCode != http.StatusOK {
    return nil, &APIError{
        Message:    fmt.Sprintf("HTTP error: %s", resp.Status),
        StatusCode: resp.StatusCode,
        Endpoint:   endpoint,
    }
}

resultCode, ok := result["Result"].(float64)
if !ok || resultCode != 200 {
    return nil, &APIError{
        Message:    errorMsg,
        StatusCode: int(resultCode),
        Endpoint:   endpoint,
    }
}
```

## Cryptography

### Secp256k1 Implementation

The SDK uses `btcsuite/btcd/btcec/v2` for secp256k1 cryptography:

```go
import "github.com/btcsuite/btcd/btcec/v2"
```

### Key Operations

#### 1. Derive Public Key

```go
func (c *Client) GetPublicKey(privateKey string) (string, error) {
    // Remove 0x prefix
    cleanKey := hexFix(privateKey)

    // Parse private key bytes
    privKeyBytes, err := hex.DecodeString(cleanKey)

    // Generate secp256k1 private key
    privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

    // Return uncompressed public key (64 bytes, no 0x04 prefix)
    return hex.EncodeToString(pubKeyBytes[1:]), nil
}
```

#### 2. Sign Message

```go
func (c *Client) SignMessage(message string, privateKey string) (string, error) {
    // Parse private key
    privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

    // Hash message with SHA-256
    msgHash := sha256.Sum256([]byte(message))

    // Sign with ECDSA (DER format)
    signature, err := ecdsa.SignASN1(rand.Reader, privKey.ToECDSA(), msgHash[:])

    return hex.EncodeToString(signature), nil
}
```

#### 3. Verify Signature

```go
func (c *Client) VerifySignature(publicKey string, message string, signatureHex string) bool {
    // Parse public key (add 0x04 prefix if needed)
    if len(pubKeyBytes) == 64 {
        pubKeyBytes = append([]byte{0x04}, pubKeyBytes...)
    }

    // Parse signature
    pubKey, err := btcec.ParsePubKey(pubKeyBytes)

    // Hash message
    msgHash := sha256.Sum256([]byte(message))

    // Verify DER-encoded signature
    return ecdsa.VerifyASN1(pubKey.ToECDSA(), msgHash[:], signatureBytes)
}
```

### Address Generation

Addresses are generated by hashing the public key:

```go
func (c *Client) HashString(str string) string {
    hash := sha256.Sum256([]byte(str))
    return hex.EncodeToString(hash[:])
}

// Usage:
address := client.HashString(publicKey)
```

## Testing Strategy

### Test Levels

#### 1. Unit Tests (`circular_protocol_test.go`)

Tests individual functions in isolation:

```go
func TestNewClient(t *testing.T) {
    client := circularprotocol.NewClient("http://test.com", "test-key")
    // Assert client configuration
}
```

#### 2. Integration Tests (`circular_protocol_integration_test.go`)

Tests API interactions with mock servers:

```go
mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"Result": 200, "Response": {...}}`))
}))

client := circularprotocol.NewClient(mockServer.URL, "mock-key")
result, err := client.CheckWallet(ctx, request)
```

#### 3. E2E Tests (`circular_protocol_e2e_test.go`)

Tests against real endpoints (optional, requires configuration):

```go
func TestE2E_CheckWallet(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping e2e test in short mode")
    }

    client := circularprotocol.NewClient(realNAGURL, realAPIKey)
    // Test against real blockchain
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run only unit tests (fast)
go test ./... -short

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestCheckWallet

# Run with race detection
go test -race ./...
```

### Mock Strategy

- **HTTP Mocking**: Use `httptest.Server` for HTTP mocking
- **Response Mocking**: Return predefined JSON responses
- **Error Testing**: Test error handling paths
- **Context Testing**: Test context cancellation and timeouts

## Migration Guide

### From Legacy API to Modern SDK

#### Before (Legacy API)

```go
import "github.com/circular-protocol/circular-go/circular_protocol_api"

result := circular_protocol_api.CheckWallet(blockchain, address)
if result["Result"].(float64) != 200 {
    // Handle error
}
```

#### After (Modern SDK)

```go
import circular "github.com/circular-protocol/circular-go"

client := circular.NewClient(nagURL, apiKey)
ctx := context.Background()

req := map[string]interface{}{
    "Address":    address,
    "Blockchain": blockchain,
    "Version":    "1.0.8",
}

result, err := client.CheckWallet(ctx, req)
if err != nil {
    // Handle error
}
```

### Key Differences

| Aspect | Legacy API | Modern SDK |
|--------|-----------|------------|
| **Style** | Package-level functions | Client methods |
| **Context** | Not supported | Full support |
| **Errors** | Return in response | Go error type |
| **Configuration** | Global variables | Client instance |
| **Testing** | Difficult to mock | Easy to mock |
| **Timeout** | No built-in support | Via context |

### Backward Compatibility

Both implementations are maintained:
- **Legacy API**: `circular_protocol_api` package still works
- **Modern SDK**: New `circularprotocol` package
- **Utils**: Shared between both (with aliases for compatibility)

Example of backward-compatible changes:

```go
// Old function still works (alias)
func Sha256(data string) string {
    return HashString(data)
}

// New function (primary)
func HashString(data string) string {
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

## Best Practices

### 1. Always Use Context

```go
// Good: Pass context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
result, err := client.CheckWallet(ctx, req)

// Bad: No context (compilation error in modern SDK)
result := circular_protocol_api.CheckWallet(blockchain, address)
```

### 2. Reuse Client Instances

```go
// Good: Create once, use many times
client := circular.NewClient(nagURL, apiKey)
for _, addr := range addresses {
    result, err := client.GetWallet(ctx, makeRequest(addr))
}

// Bad: Create new client for each request
for _, addr := range addresses {
    client := circular.NewClient(nagURL, apiKey)
    result, err := client.GetWallet(ctx, makeRequest(addr))
}
```

### 3. Handle Errors Properly

```go
// Good: Check and handle errors
result, err := client.SendTransaction(ctx, tx)
if err != nil {
    if apiErr, ok := err.(*circular.APIError); ok {
        log.Printf("API error: %s", apiErr.Message)
    }
    return fmt.Errorf("transaction failed: %w", err)
}

// Bad: Ignore errors
result, _ := client.SendTransaction(ctx, tx)
```

### 4. Secure Key Management

```go
// Good: Load from secure storage
privateKey := os.Getenv("PRIVATE_KEY")
signature, err := client.SignMessage(message, privateKey)

// Bad: Hardcode keys
signature, err := client.SignMessage(message, "0x1234...")
```

## Performance Considerations

### HTTP Client Reuse

The SDK reuses HTTP connections:

```go
// HTTP client is reused across requests
client := circular.NewClient(nagURL, apiKey)

// Each request reuses the connection pool
for i := 0; i < 1000; i++ {
    client.CheckWallet(ctx, req) // Efficient connection reuse
}
```

### Connection Pooling

Customize the HTTP client for specific needs:

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := circular.NewClientWithConfig(circular.Config{
    HTTPClient: httpClient,
})
```

### Context Timeouts

Use appropriate timeouts:

```go
// Short timeout for quick operations
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
wallet, err := client.CheckWallet(ctx, req)

// Longer timeout for transactions
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
tx, err := client.SendTransaction(ctx, txReq)
```

## Future Enhancements

Planned improvements for future versions:

1. **Typed Requests/Responses**: Replace `map[string]interface{}` with typed structs
2. **Retry Logic**: Automatic retry with exponential backoff
3. **Rate Limiting**: Built-in rate limiting for API calls
4. **Batch Operations**: Support for batch API requests
5. **WebSocket Support**: Real-time blockchain updates
6. **Event Streaming**: Subscribe to blockchain events
7. **Smart Contract SDK**: Higher-level abstractions for contracts
8. **Wallet Management**: Built-in wallet creation and management

## Resources

- [Go Documentation](https://pkg.go.dev/github.com/circular-protocol/circular-go)
- [Circular Protocol Docs](https://circular-protocol.gitbook.io)
- [GitHub Repository](https://github.com/circular-protocol/circular-go)
- [Issue Tracker](https://github.com/circular-protocol/circular-go/issues)

---

**Last Updated**: November 2024
**Version**: 1.0.9 (Unreleased)
