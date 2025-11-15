# Changelog

All notable changes to the Circular Protocol Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Modern context-aware SDK implementation (`circular_protocol.go`)
- Comprehensive test suite (unit, integration, and e2e tests)
- GitHub Actions CI/CD workflow for automated testing
- Complete project documentation (CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md, AGENTS.md)
- Context support for all API methods with `context.Context` parameter
- Proper Go error handling with `*APIError` type
- Client configuration options via `Config` struct
- `GetTransactionOutcome` method for polling transaction confirmation
- Cryptographic helper methods (SignMessage, VerifySignature, GetPublicKey, HashString)
- Encoding helper methods (StringToHex, HexToString, HexFix)
- `RegisterWallet` convenience method
- Professional README with comprehensive API reference

### Changed
- **Method naming aligned with API specification**:
  - Legacy `AddTransaction` → `SendTransaction` (matches API spec)
  - Legacy `GetTransactionbyID` → `GetTransactionByID` (proper Go naming)
  - Legacy `GetTransactionbyNode` → `GetTransactionByNode`
  - Legacy `GetTransactionbyAddress` → `GetTransactionByAddress`
  - Legacy `GetTransactionbyDate` → `GetTransactionByDate`
  - Legacy `GetLatestTransaction` → `GetLatestTransactions` (plural)
- Improved `.gitignore` with cleaner, more focused patterns
- Updated `utils.HashString` as primary function with `Sha256` as backward-compatible alias

### Fixed
- `RegisterWallet` now correctly uses client methods instead of undefined package functions
- Added missing imports in `circular_protocol.go`
- Fixed godoc comment formatting (multi-line comments now use `//` prefix)
- Fixed test file missing package qualifier

### Security
- Added comprehensive security documentation in SECURITY.md
- Documented best practices for API key and private key management
- Included secure coding guidelines for transaction handling

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
