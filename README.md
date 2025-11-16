# Circular Protocol - Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/circular-protocol/circular-go.svg)](https://pkg.go.dev/github.com/circular-protocol/circular-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/circular-protocol/circular-go)](https://goreportcard.com/report/github.com/circular-protocol/circular-go)

The **Circular Protocol Go SDK** is the official Go library for seamless integration with the Circular blockchain ecosystem. This open-source SDK provides a comprehensive suite of tools for efficient and secure interaction with blockchain networks, managing wallets, assets, smart contracts, and more.

## 🔥 Key Features

- **Blockchain Interaction**: Connect and interact with Circular's blockchain networks
- **Smart Contracts**: Deploy, test, and interact with smart contracts
- **Wallet Management**: Create, retrieve, and manage blockchain wallets with balance tracking
- **Asset Management**: Issue and manage assets, handle transfers, and retrieve supply information
- **Domain Management**: Resolve blockchain domain names to wallet addresses
- **Transaction Management**: Send transactions, track status, and search the blockchain
- **Analytics**: Access blockchain performance data and insights
- **Cryptographic Helpers**: Built-in utilities for key generation, signing, and hashing
- **Go 1.19+ Support**: Compatible with modern Go versions
- **Context Support**: Full context.Context integration for cancellation and timeouts

```bash
go get github.com/circular-protocol/circular-go
```

## 🚀 Quick Start

The SDK supports two API styles for maximum flexibility:

### Option 1: Convenience Methods (Recommended)

Simpler API with auto-preprocessing and version injection:

```go
package main

import (
    "context"
    "fmt"
    "log"

    circularprotocol "github.com/circular-protocol/circular-go"
)

func main() {
    // Initialize the client
    client := circularprotocol.NewClient(
        "https://nag.circularlabs.io/NAG.php?cep=",
        "", // API key (optional)
    )

    ctx := context.Background()

    // Check if a wallet exists - auto-strips '0x', auto-injects version
    result, err := client.CheckWallet(ctx, "MainNet", "0xd55872dbe508fd27445889b9d81bbc9411bb0f1353153a249f2fb34ef2690310")
    if err != nil {
        log.Fatalf("API Error: %v", err)
    }

    fmt.Printf("Wallet exists: %v\n", result["Response"])
}
```

### Option 2: Raw Methods (Explicit Control)

Full control over request parameters, no auto-preprocessing:

```go
package main

import (
    "context"
    "fmt"
    "log"

    circularprotocol "github.com/circular-protocol/circular-go"
)

func main() {
    client := circularprotocol.NewClient(
        "https://nag.circularlabs.io/NAG.php?cep=",
        "", // API key (optional)
    )

    ctx := context.Background()

    // Check wallet with explicit request object
    params := map[string]interface{}{
        "Address":    "d55872dbe508fd27445889b9d81bbc9411bb0f1353153a249f2fb34ef2690310",
        "Blockchain": "MainNet",
        "Version":    "1.0.9",
    }

    result, err := client.CheckWalletRaw(ctx, params)
    if err != nil {
        log.Fatalf("API Error: %v", err)
    }

    fmt.Printf("Wallet exists: %v\n", result["Response"])
}
```

### Key Differences

| Feature | Convenience Methods | Raw Methods |
|---------|---------------------|-------------|
| **Method Names** | `CheckWallet`, `GetWallet`, etc. | `CheckWalletRaw`, `GetWalletRaw`, etc. |
| **Parameters** | Positional (blockchain, address, ...) | Request object (map[string]interface{}) |
| **Auto-preprocessing** | ✅ Strips '0x', converts to hex | ❌ Manual |
| **Version injection** | ✅ Automatic | ❌ Manual |
| **Best for** | Rapid development, cleaner code | Advanced use cases, explicit control |

## 📜 API Reference

The Circular Protocol Go SDK provides **37 methods** across multiple categories for comprehensive blockchain interaction.

### Wallet Operations (6 methods)

- **`CheckWallet`** - Verify wallet existence on the blockchain
- **`GetWallet`** - Retrieve complete wallet details and metadata
- **`GetLatestTransactions`** - Get recent wallet activity and transaction history
- **`GetWalletBalance`** - Query current wallet balance across assets
- **`GetWalletNonce`** - Get transaction nonce for the wallet
- **`RegisterWallet`** - Register new wallet on the blockchain

### Transaction Operations (7 methods)

- **`SendTransaction`** - Submit new transaction to the blockchain
- **`GetPendingTransaction`** - Check transaction status in the mempool
- **`GetTransactionByID`** - Query transaction by unique identifier
- **`GetTransactionByNode`** - Query transactions by validator node
- **`GetTransactionByAddress`** - Query all transactions for a wallet address
- **`GetTransactionByDate`** - Query transactions within a date range
- **`GetTransactionOutcome`** - Poll for transaction confirmation with automatic retries

### Blockchain Operations (5 methods)

- **`GetBlock`** - Retrieve block data by block number or hash
- **`GetBlockRange`** - Query multiple blocks within a range
- **`GetBlockCount`** - Get current blockchain height (latest block number)
- **`GetAnalytics`** - Retrieve blockchain performance metrics and analytics
- **`GetBlockchains`** - List all supported blockchain networks

### Contract Operations (2 methods)

- **`TestContract`** - Validate smart contract logic before deployment
- **`CallContract`** - Execute smart contract function call

### Asset Operations (4 methods)

- **`GetAssetList`** - List all available assets on the blockchain
- **`GetAsset`** - Get detailed asset information and metadata
- **`GetAssetSupply`** - Query total and circulating supply for an asset
- **`GetVoucher`** - Retrieve voucher data and redemption details

### Domain Operations (1 method)

- **`GetDomain`** - Query blockchain domain registry (resolve domain to address)

---

### Cryptographic Utilities (4 methods)

- **`SignMessage`** - Generate ECDSA secp256k1 signatures (DER format)
- **`VerifySignature`** - Verify message signatures against public keys
- **`GetPublicKey`** - Derive public key from private key (128 hex characters, uncompressed, no 0x04 prefix)
- **`HashString`** - Generate SHA-256 hash of string input

**Implementation Details:**
- **TypeScript/JavaScript**: `crypto-browserify` (browser-compatible)
- **Python**: `ecdsa` + `hashlib` (standard library)
- **Java**: Bouncy Castle library for secp256k1
- **PHP**: `phpseclib3` elliptic curve cryptography
- **Go**: `btcsuite/btcd/btcec/v2` secp256k1
- **Dart**: `pointycastle` package

---

### Encoding/Format Utilities (4 methods)

- **`HexFix`** - Normalize hex strings (remove `0x` prefix if present)
- **`StringToHex`** - Convert UTF-8 string to hexadecimal encoding
- **`HexToString`** - Convert hexadecimal string to UTF-8
- **`GetFormattedTimestamp`** - Get current UTC timestamp in Circular Protocol format (`YYYY:MM:DD-HH:mm:ss`)

---

### NAG Configuration (4 methods)

- **`GetNagUrl`** - Retrieve current NAG endpoint URL
- **`SetNagUrl`** - Configure NAG endpoint URL
- **`GetNagKey`** - Retrieve current API key
- **`SetNagKey`** - Configure API key for authenticated requests

---

## 📊 Total Methods: 37

Breaking down by category:
- **6** Wallet Operations
- **7** Transaction Operations
- **5** Blockchain Operations
- **2** Smart Contract Operations
- **4** Asset Management
- **1** Domain Management
- **4** Cryptographic Utilities
- **4** Encoding/Format Utilities
- **4** NAG Configuration

> **Note**: For detailed parameter types, response structures, and advanced usage examples, refer to the **[Go SDK Documentation](https://circular-protocol.gitbook.io/circular-sdk/api-docs/go)**.

## 🤝 Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](https://github.com/circular-protocol/circular-canonical/blob/main/CONTRIBUTING.md) file in the canonical repository for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 Resources

- **[Go SDK Documentation](https://circular-protocol.gitbook.io/circular-sdk/api-docs/go)** - Complete API reference
- **[Circular Protocol Docs](https://circular-protocol.gitbook.io)** - Protocol documentation
- **[Circular Canonical](https://github.com/circular-protocol/circular-canonical)** - Single source of truth
- **[Package on pkg.go.dev](https://pkg.go.dev/github.com/circular-protocol/circular-go)** - Official Go package

## ℹ️ About

**Version**: 1.0.8
**License**: MIT
**Generated**: Auto-generated from [Circular Canonical](https://github.com/circular-protocol/circular-canonical) specification

---

© 2025 Circular Global Ledgers, Inc. - Open source for private and commercial use