# Go SDK - TODO

## ✅ Completed - Utility Methods Implementation

All utility/configuration methods have been successfully implemented to match Python and Dart SDKs!

### Configuration Methods ✅
- [x] `SetNAGURL(url string)` - Update NAG endpoint URL at runtime
- [x] `GetNAGURL() string` - Get current NAG endpoint URL
- [x] `SetNAGKey(key string)` - Update NAG API key at runtime
- [x] `GetNAGKey() string` - Get current NAG API key
- [x] `SetHeader(key string, value string)` - Set custom HTTP headers
- [x] `GetVersion() string` - Get SDK version

### Error Handling Methods ✅
- [x] `GetError() string` - Get last error message from SDK
- [x] `HandleError(result map[string]interface{}) bool` - Handle API error responses

### Lifecycle Methods ✅
- [x] `Dispose()` - Clean up resources (HTTP client, etc.)

### Additional Utility Methods ✅
- [x] `SetNode(address string)` - Set primary node address for querying blockchain
- [x] `GetNode() string` - Get current node address (bonus method)

## Implementation Summary

All 10+ methods have been added to the `Client` struct in `circular_protocol.go`:

1. **Configuration Methods**: Allow runtime configuration of NAG URL, API key, and custom headers
2. **Error Handling**: Track last error and provide error handling utilities
3. **Lifecycle Management**: Proper resource cleanup with `Dispose()`
4. **Node Management**: Set/get primary node for blockchain queries

The Go SDK now has full feature parity with Python and Dart SDKs while following Go naming conventions and best practices.
