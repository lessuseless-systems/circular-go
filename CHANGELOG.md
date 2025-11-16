# Changelog

All notable changes to the Circular Protocol Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.9] - 2025-11-16

Complete SDK alignment with circular-js-npm implementation. All changes are **100% backward compatible** - existing code continues to work unchanged.

### ✨ Added

**Dual Signature Pattern** - All 24 API methods now support both calling styles:
- **Convenience methods** (positional params) - Auto-preprocessing, version injection, simpler API
- **Raw methods** (request object) - Explicit control, no auto-preprocessing

**Auto-Preprocessing System**:
- `HexFix()` automatically strips '0x' prefix from blockchain, address, ID, nodeID parameters
- `StringToHex()` auto-converts project/request strings in contract methods
- Version automatically injected (`Version: "1.0.9"`)
- Timestamps auto-generated for contract methods

**New Convenience Methods** (24 total):
- Wallet: `CheckWallet`, `GetWallet`, `GetLatestTransactions`, `GetWalletBalance`, `GetWalletNonce`
- Transactions: `SendTransaction` (9 positional params), `GetPendingTransaction`, `GetTransactionByID`, `GetTransactionByNode`, `GetTransactionByAddress`, `GetTransactionByDate`
- Blocks: `GetBlock`, `GetBlockRange`, `GetBlockCount`, `GetAnalytics`
- Contracts: `TestContract`, `CallContract` (with auto-encoding)
- Assets: `GetAssetList`, `GetAsset`, `GetAssetSupply`, `GetVoucher`
- Network: `GetDomain`, `GetBlockchains`

**Exported Helper**:
- `HexFix()` method now exported for manual preprocessing

### 🔧 Changed

- **SDK Version**: `1.0.8` → `1.0.9`
- **All existing methods renamed to `*Raw`** suffix (backward compatibility maintained via convenience methods)
- **SendTransaction signature updated** to match JavaScript implementation (9 positional parameters)
- **RegisterWallet updated** to use new SendTransaction signature
- **Version management**: Now uses package constant `Version = "1.0.9"`
- **Client struct**: Added `version` field for automatic injection

### 🐛 Fixed

**Critical Endpoint Naming Bugs**:
- `GetDomain`: Fixed incorrect endpoint `"GetDomain"` → `"ResolveDomain"`
- `GetBlockCount`: Fixed incorrect endpoint `"GetBlockCount"` → `"GetBlockHeight"`

These were bugs that would have caused API calls to fail against the actual Circular Protocol API.

### 📚 Documentation

- Enhanced CONTRIBUTING.md with dual signature pattern guidelines
- Updated with auto-preprocessing rules and examples
- Added comprehensive method implementation examples

### 🎯 Migration Guide

**No migration required!** This release is 100% backward compatible.

**Optional**: Adopt new convenience methods for cleaner code:

```go
// Before (still works)
result, err := client.CheckWalletRaw(ctx, map[string]interface{}{
    "Blockchain": "MainNet",
    "Address":    "742d35...",
    "Version":    "1.0.9",
})

// After (recommended - more concise)
result, err := client.CheckWallet(ctx, "0xMainNet", "0x742d35...")
// Auto-strips '0x', auto-injects version
```

### 🔍 Technical Details

**Auto-Preprocessing Table**:

| Function | Applied To | Example |
|----------|-----------|---------|
| `HexFix()` | blockchain, address, ID, nodeID | `'0x123'` → `'123'` |
| `StringToHex()` | project, request (contracts) | `'hello'` → `'68656c6c6f'` |
| Version injection | All requests | Auto-adds `Version: "1.0.9"` |
| Timestamp generation | Contract methods | Auto-generates UTC timestamp |
| Code stripping | Voucher codes | Strips '0x' from codes |

**Methods with Dual Signatures**: 24 core API methods (all wallet, transaction, block, contract, asset, and network operations)

### ✅ Checklist

- [x] All 24 API methods have dual signatures
- [x] All methods auto-preprocess parameters
- [x] SendTransaction matches JavaScript signature (9 params)
- [x] Endpoint names corrected (getDomain, getBlockCount)
- [x] Version bumped to 1.0.9
- [x] RegisterWallet updated to new SendTransaction
- [x] HexFix exported as public method
- [x] CONTRIBUTING.md enhanced
- [x] **100% Backward compatibility maintained**

## [1.0.1] - 2024-11-XX

### Fixed
- GetTransactionOutcome polling implementation
- Transaction finality checking logic

## [1.0.0] - 2024-11-XX

### Added
- Initial release of Circular Protocol Go SDK
- Legacy API implementation in `circular_protocol_api` package
- Core wallet operations (CheckWallet, GetWallet, GetWalletBalance, GetWalletNonce)
- Transaction operations (SendTransaction, GetTransactionByID, GetPendingTransaction)
- Block operations (GetBlock, GetBlockRange, GetBlockCount, GetAnalytics)
- Smart contract operations (TestContract, CallContract)
- Asset operations (GetAsset, GetAssetList, GetAssetSupply, GetVoucher)
- Domain resolution (GetDomain)
- Network operations (GetBlockchains)
- Utility functions (HexFix, StringToHex, HexToString, SignMessage, VerifySignature)
- Wallet registration functionality
- Basic README documentation
- MIT License

### Changed
- Package structure organization with `circular_protocol_api` and `utils` packages

---

## Version History

- **1.0.1**: Bug fixes for transaction polling
- **1.0.0**: Initial public release
- **Unreleased**: Modern SDK with context support, comprehensive tests, and documentation

## Links

- [GitHub Repository](https://github.com/circular-protocol/circular-go)
- [Documentation](https://circular-protocol.gitbook.io/circular-sdk/api-docs/go)
- [Package on pkg.go.dev](https://pkg.go.dev/github.com/circular-protocol/circular-go)
- [Circular Protocol](https://circular-protocol.gitbook.io)

---

**Note**: This changelog follows the [Keep a Changelog](https://keepachangelog.com/) format and uses [Semantic Versioning](https://semver.org/).
