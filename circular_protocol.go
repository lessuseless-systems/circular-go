// Package circularprotocol provides a Go SDK for the Circular Protocol blockchain API.
// Generated from Nickel API specification
// Version: 1.0.9
//
// This SDK provides dual API patterns for maximum flexibility:
//
// 1. Convenience methods with positional parameters (auto-preprocessing):
//
//	client := circularprotocol.NewClient("https://nag.circularlabs.io/NAG.php?cep=", "")
//	// Auto-strips '0x' prefix, auto-injects version
//	result, err := client.CheckWallet(ctx, "MainNet", "0x742d35...")
//
// 2. Request object style (explicit control):
//
//	result, err := client.CheckWalletRaw(ctx, map[string]interface{}{
//		"Address":    "0x...",
//		"Blockchain": "714d2ac07a826b66ac56752eebd7c77b58d2ee842e523d913fd0ef06e6bdfcae",
//		"Version":    "1.0.9",
//	})
//
// Both patterns are fully supported and maintained.
package circularprotocol

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
)

const (
	// Version is the current SDK version
	Version = "1.0.9"
)

// Client represents a Circular Protocol API client
type Client struct {
	nagURL     string
	nagKey     string
	httpClient *http.Client
	headers    map[string]string
	version    string
}

// Config holds client configuration options
type Config struct {
	NAGURL     string
	NAGKey     string
	Timeout    time.Duration
	HTTPClient *http.Client
	Headers    map[string]string
}

// APIError represents an error returned by the Circular Protocol API
type APIError struct {
	Message    string
	StatusCode int
	Endpoint   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error on %s (status %d): %s", e.Endpoint, e.StatusCode, e.Message)
}

// NewClient creates a new Circular Protocol API client
func NewClient(nagURL, nagKey string) *Client {
	return NewClientWithConfig(Config{
		NAGURL:  nagURL,
		NAGKey:  nagKey,
		Timeout: 30 * time.Second,
	})
}

// NewClientWithConfig creates a new client with custom configuration
func NewClientWithConfig(cfg Config) *Client {
	nagURL := cfg.NAGURL
	if nagURL == "" {
		nagURL = "https://nag.circularlabs.io/NAG.php?cep="
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
		if cfg.Timeout == 0 {
			httpClient.Timeout = 30 * time.Second
		}
	}

	headers := cfg.Headers
	if headers == nil {
		headers = make(map[string]string)
	}

	return &Client{
		nagURL:     nagURL,
		nagKey:     cfg.NAGKey,
		httpClient: httpClient,
		headers:    headers,
		version:    Version,
	}
}

// SetNAGURL updates the NAG endpoint URL
func (c *Client) SetNAGURL(url string) {
	c.nagURL = url
}

// GetNAGURL returns the current NAG endpoint URL
func (c *Client) GetNAGURL() string {
	return c.nagURL
}

// SetNAGKey updates the NAG API key
func (c *Client) SetNAGKey(key string) {
	c.nagKey = key
}

// GetNAGKey returns the current NAG API key
func (c *Client) GetNAGKey() string {
	return c.nagKey
}

// SetHeader sets a custom HTTP header
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// makeRequest performs an HTTP request to the NAG API
func (c *Client) makeRequest(ctx context.Context, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	url := c.nagURL + "Circular_" + endpoint + "_"

	// Marshal request data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
			StatusCode: 0,
			Endpoint:   endpoint,
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to create request: %v", err),
			StatusCode: 0,
			Endpoint:   endpoint,
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	if c.nagKey != "" {
		req.Header.Set("X-NAG-Key", c.nagKey)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("request failed: %v", err),
			StatusCode: 0,
			Endpoint:   endpoint,
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to read response: %v", err),
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Message:    fmt.Sprintf("HTTP error: %s", resp.Status),
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to parse response: %v", err),
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}
	}

	// Check API-level errors
	resultCode, ok := result["Result"].(float64)
	if !ok || resultCode != 200 {
		errorMsg := "API request failed"
		if response, ok := result["Response"].(string); ok {
			errorMsg = response
		}
		return nil, &APIError{
			Message:    errorMsg,
			StatusCode: int(resultCode),
			Endpoint:   endpoint,
		}
	}

	// Return full response (with Result and Response fields)
	return result, nil
}

// ============================================================================
// API Methods - Raw (Request Object Style)
// ============================================================================
// These methods accept a map[string]interface{} request object and provide
// explicit control over all parameters. No auto-preprocessing is applied.

// CheckWalletRaw checks if wallet exists (request object style)
// Checks whether a wallet address exists on the specified blockchain.
// Returns existence status and confirms the address format.
//
// Parameters (via request map):
//   - Blockchain: Blockchain identifier (no auto-preprocessing)
//   - Address: Wallet address (no auto-preprocessing)
//   - Version: SDK version (required)
//
// Example:
//
//	result, err := client.CheckWalletRaw(ctx, map[string]interface{}{
//		"Blockchain": "MainNet",
//		"Address":    "742d35...",
//		"Version":    "1.0.9",
//	})
func (c *Client) CheckWalletRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "CheckWallet", req)
}

// CheckWallet checks if wallet exists (convenience method with auto-preprocessing)
// Checks whether a wallet address exists on the specified blockchain.
//
// This method automatically:
//   - Strips '0x' prefix from blockchain and address parameters
//   - Injects Version field automatically
//
// Parameters:
//   - ctx: Context for request cancellation/timeout
//   - blockchain: Blockchain identifier (e.g., '0xMainNet' or 'MainNet')
//   - address: Wallet address (e.g., '0x742d35...' or '742d35...')
//
// Returns:
//   - map with "Result" (int) and "Response" fields
//   - error if request fails
//
// Example:
//
//	// Auto-strips '0x' prefix from both parameters
//	exists, err := client.CheckWallet(ctx, "0xMainNet", "0x742d35...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Wallet exists: %v\n", exists["Response"])
//
// See also:  GetWallet, GetWalletBalance
func (c *Client) CheckWallet(ctx context.Context, blockchain string, address string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Version":    c.version,
	}
	return c.CheckWalletRaw(ctx, req)
}

// GetWalletRaw gets wallet information (request object style)
// Retrieves complete wallet information including balance and nonce.
func (c *Client) GetWalletRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWallet", req)
}

// GetWallet gets complete wallet information (convenience method with auto-preprocessing)
// Retrieves comprehensive wallet information including balance, nonce, and state.
//
// This method automatically:
//   - Strips '0x' prefix from blockchain and address
//   - Injects Version field
//
// Parameters:
//   - ctx: Context for request cancellation/timeout
//   - blockchain: Blockchain identifier
//   - address: Wallet address
//
// Example:
//
//	wallet, err := client.GetWallet(ctx, "MainNet", "0x742d35...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Balance: %v, Nonce: %v\n", wallet["Response"].(map[string]interface{})["Balance"],
//		wallet["Response"].(map[string]interface{})["Nonce"])
//
// See also: CheckWallet, GetWalletBalance, GetWalletNonce
func (c *Client) GetWallet(ctx context.Context, blockchain string, address string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Version":    c.version,
	}
	return c.GetWalletRaw(ctx, req)
}

// GetLatestTransactionsRaw gets latest transactions (request object style)
func (c *Client) GetLatestTransactionsRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetLatestTransactions", req)
}

// GetLatestTransactions gets latest transactions for wallet (convenience method)
func (c *Client) GetLatestTransactions(ctx context.Context, blockchain string, address string) (map[string]interface{}, error) {
	return c.GetLatestTransactionsRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Version":    c.version,
	})
}

// GetWalletBalanceRaw gets wallet balance (request object style)
func (c *Client) GetWalletBalanceRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWalletBalance", req)
}

// GetWalletBalance gets wallet balance for specific asset (convenience method)
func (c *Client) GetWalletBalance(ctx context.Context, blockchain string, address string, asset string) (map[string]interface{}, error) {
	return c.GetWalletBalanceRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Asset":      asset,
		"Version":    c.version,
	})
}

// GetWalletNonceRaw gets wallet nonce (request object style)
func (c *Client) GetWalletNonceRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWalletNonce", req)
}

// GetWalletNonce gets wallet nonce (convenience method)
func (c *Client) GetWalletNonce(ctx context.Context, blockchain string, address string) (map[string]interface{}, error) {
	return c.GetWalletNonceRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Version":    c.version,
	})
}

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

// GetPendingTransactionRaw gets pending transaction (request object style)
func (c *Client) GetPendingTransactionRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetPendingTransaction", req)
}

// GetPendingTransaction gets pending transaction by ID (convenience method)
func (c *Client) GetPendingTransaction(ctx context.Context, blockchain string, txID string) (map[string]interface{}, error) {
	return c.GetPendingTransactionRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"ID":         c.HexFix(txID),
		"Version":    c.version,
	})
}

// GetTransactionByIDRaw finds transaction by ID (request object style)
func (c *Client) GetTransactionByIDRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyID", req)
}

// GetTransactionByID finds transaction by ID (convenience method)
func (c *Client) GetTransactionByID(ctx context.Context, blockchain string, txID string, start string, end string) (map[string]interface{}, error) {
	return c.GetTransactionByIDRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"ID":         c.HexFix(txID),
		"Start":      start,
		"End":        end,
		"Version":    c.version,
	})
}

// GetTransactionByNodeRaw finds transactions by node (request object style)
func (c *Client) GetTransactionByNodeRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyNode", req)
}

// GetTransactionByNode finds transactions by node ID (convenience method)
func (c *Client) GetTransactionByNode(ctx context.Context, blockchain string, nodeID string, start string, end string) (map[string]interface{}, error) {
	return c.GetTransactionByNodeRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"NodeID":     c.HexFix(nodeID),
		"Start":      start,
		"End":        end,
		"Version":    c.version,
	})
}

// GetTransactionByAddressRaw finds transactions by address (request object style)
func (c *Client) GetTransactionByAddressRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyAddress", req)
}

// GetTransactionByAddress finds transactions by address (convenience method)
func (c *Client) GetTransactionByAddress(ctx context.Context, blockchain string, address string, start string, end string) (map[string]interface{}, error) {
	return c.GetTransactionByAddressRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"Start":      start,
		"End":        end,
		"Version":    c.version,
	})
}

// GetTransactionByDateRaw finds transactions by date (request object style)
func (c *Client) GetTransactionByDateRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyDate", req)
}

// GetTransactionByDate finds transactions by date range (convenience method)
func (c *Client) GetTransactionByDate(ctx context.Context, blockchain string, address string, startDate string, endDate string) (map[string]interface{}, error) {
	return c.GetTransactionByDateRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"StartDate":  startDate,
		"EndDate":    endDate,
		"Version":    c.version,
	})
}

// GetBlockRaw gets specific block (request object style)
func (c *Client) GetBlockRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlock", req)
}

// GetBlock gets specific block by number (convenience method)
func (c *Client) GetBlock(ctx context.Context, blockchain string, blockNumber string) (map[string]interface{}, error) {
	return c.GetBlockRaw(ctx, map[string]interface{}{
		"Blockchain":  c.HexFix(blockchain),
		"BlockNumber": blockNumber,
		"Version":     c.version,
	})
}

// GetBlockRangeRaw gets range of blocks (request object style)
func (c *Client) GetBlockRangeRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockRange", req)
}

// GetBlockRange gets range of blocks (convenience method)
func (c *Client) GetBlockRange(ctx context.Context, blockchain string, start string, end string) (map[string]interface{}, error) {
	return c.GetBlockRangeRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Start":      start,
		"End":        end,
		"Version":    c.version,
	})
}

// GetBlockCountRaw gets blockchain height (request object style)
func (c *Client) GetBlockCountRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockHeight", req)
}

// GetBlockCount gets blockchain height (convenience method)
func (c *Client) GetBlockCount(ctx context.Context, blockchain string) (map[string]interface{}, error) {
	return c.GetBlockCountRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Version":    c.version,
	})
}

// GetAnalyticsRaw gets blockchain analytics (request object style)
func (c *Client) GetAnalyticsRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAnalytics", req)
}

// GetAnalytics gets blockchain analytics (convenience method)
func (c *Client) GetAnalytics(ctx context.Context, blockchain string) (map[string]interface{}, error) {
	return c.GetAnalyticsRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Version":    c.version,
	})
}

// TestContractRaw tests smart contract (request object style)
func (c *Client) TestContractRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "TestContract", req)
}

// TestContract tests smart contract execution (convenience method)
// Auto-strips '0x' prefix, auto-converts project to hex, auto-generates timestamp
func (c *Client) TestContract(ctx context.Context, blockchain string, from string, project string) (map[string]interface{}, error) {
	return c.TestContractRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"From":       c.HexFix(from),
		"Project":    c.StringToHex(project),
		"Timestamp":  c.GetFormattedTimestamp(),
		"Version":    c.version,
	})
}

// CallContractRaw calls smart contract (request object style)
func (c *Client) CallContractRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "CallContract", req)
}

// CallContract calls smart contract function (convenience method)
// Auto-strips '0x' prefix, auto-converts request to hex, auto-generates timestamp
func (c *Client) CallContract(ctx context.Context, blockchain string, address string, from string, request string) (map[string]interface{}, error) {
	return c.CallContractRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Address":    c.HexFix(address),
		"From":       c.HexFix(from),
		"Request":    c.StringToHex(request),
		"Timestamp":  c.GetFormattedTimestamp(),
		"Version":    c.version,
	})
}

// GetAssetListRaw lists all assets (request object style)
func (c *Client) GetAssetListRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAssetList", req)
}

// GetAssetList lists all assets on blockchain (convenience method)
func (c *Client) GetAssetList(ctx context.Context, blockchain string) (map[string]interface{}, error) {
	return c.GetAssetListRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Version":    c.version,
	})
}

// GetAssetRaw gets specific asset (request object style)
func (c *Client) GetAssetRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAsset", req)
}

// GetAsset gets specific asset information (convenience method)
func (c *Client) GetAsset(ctx context.Context, blockchain string, assetName string) (map[string]interface{}, error) {
	return c.GetAssetRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"AssetName":  assetName,
		"Version":    c.version,
	})
}

// GetAssetSupplyRaw gets asset supply (request object style)
func (c *Client) GetAssetSupplyRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAssetSupply", req)
}

// GetAssetSupply gets asset supply information (convenience method)
func (c *Client) GetAssetSupply(ctx context.Context, blockchain string, assetName string) (map[string]interface{}, error) {
	return c.GetAssetSupplyRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"AssetName":  assetName,
		"Version":    c.version,
	})
}

// GetVoucherRaw retrieves voucher (request object style)
func (c *Client) GetVoucherRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetVoucher", req)
}

// GetVoucher retrieves voucher information (convenience method)
// Auto-strips '0x' prefix from code
func (c *Client) GetVoucher(ctx context.Context, blockchain string, code string) (map[string]interface{}, error) {
	return c.GetVoucherRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Code":       c.HexFix(code),
		"Version":    c.version,
	})
}

// GetDomainRaw resolves domain (request object style)
func (c *Client) GetDomainRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "ResolveDomain", req)
}

// GetDomain resolves domain to wallet address (convenience method)
func (c *Client) GetDomain(ctx context.Context, blockchain string, domain string) (map[string]interface{}, error) {
	return c.GetDomainRaw(ctx, map[string]interface{}{
		"Blockchain": c.HexFix(blockchain),
		"Domain":     domain,
		"Version":    c.version,
	})
}

// GetBlockchainsRaw lists available blockchains (request object style)
func (c *Client) GetBlockchainsRaw(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockchains", req)
}

// GetBlockchains lists available blockchains (convenience method)
func (c *Client) GetBlockchains(ctx context.Context) (map[string]interface{}, error) {
	return c.GetBlockchainsRaw(ctx, map[string]interface{}{
		"Version": c.version,
	})
}

// ============================================================================
// Convenience Methods
// ============================================================================
// These methods wrap underlying API calls to simplify common workflows

// Register wallet on blockchain (Convenience Method)
// Registers a wallet on the specified blockchain by creating and sending
a C_TYPE_REGISTERWALLET transaction. This convenience method handles all
transaction construction internally:

- Derives From/To addresses from public key (sha256)
- Builds Payload: hex(JSON.stringify({Action: "CP_REGISTERWALLET", PublicKey: publicKey}))
- Calculates transaction ID: sha256(blockchain + from + to + payload + nonce + timestamp)
- Sets Nonce to "0" and Signature to "" (empty for registration)
- Calls sendTransaction with constructed parameters

Without registration, the wallet will not be reachable on the blockchain.
The same wallet can be registered on multiple blockchains.
//
// This is a convenience method that wraps sendTransaction().
// It handles transaction construction internally.
func (c *Client) RegisterWallet(ctx context.Context, blockchain string, publicKey string) (map[string]interface{}, error) {
	// Derive addresses from public key
	from := c.HashString(publicKey)
	to := from
	nonce := "0"
	txType := "C_TYPE_REGISTERWALLET"

	// Build payload
	payloadObj := map[string]string{
		"Action":    "CP_REGISTERWALLET",
		"PublicKey": publicKey,
	}

	payloadJSON, err := json.Marshal(payloadObj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	payload := c.StringToHex(string(payloadJSON))
	timestamp := c.GetFormattedTimestamp()

	// Calculate transaction ID
	id := c.HashString(blockchain + from + to + payload + nonce + timestamp)
	signature := ""

	// Call SendTransaction with positional parameters
	return c.SendTransaction(ctx, id, from, to, timestamp, txType, payload, nonce, signature, blockchain)
}

// ============================================================================
// Helper Methods - Cryptography
// ============================================================================
// SignMessage signs a message using secp256k1 with DER encoding
// message: Message to sign (will be SHA256 hashed)
// privateKey: Private key in hex format (with or without '0x' prefix)
// Returns: DER-encoded signature as hex string
func (c *Client) SignMessage(message string, privateKey string) (string, error) {
	// Remove 0x prefix if present
	cleanKey := hexFix(privateKey)

	// Parse private key
	privKeyBytes, err := hex.DecodeString(cleanKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// Hash the message
	msgHash := sha256.Sum256([]byte(message))

	// Sign with compact signature (Bitcoin-style)
	signature, err := ecdsa.SignASN1(rand.Reader, privKey.ToECDSA(), msgHash[:])
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a DER-encoded signature
// publicKey: Public key in hex format (uncompressed, 64 bytes without 0x04 prefix)
// message: Original message that was signed
// signatureHex: DER-encoded signature in hex format
// Returns: true if signature is valid, false otherwise
func (c *Client) VerifySignature(publicKey string, message string, signatureHex string) bool {
	// Parse public key
	cleanPubKey := hexFix(publicKey)
	pubKeyBytes, err := hex.DecodeString(cleanPubKey)
	if err != nil {
		return false
	}

	// Add uncompressed point prefix if needed (0x04)
	if len(pubKeyBytes) == 64 {
		pubKeyBytes = append([]byte{0x04}, pubKeyBytes...)
	}

	// Parse public key
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false
	}

	// Hash the message
	msgHash := sha256.Sum256([]byte(message))

	// Parse signature
	signatureBytes, err := hex.DecodeString(hexFix(signatureHex))
	if err != nil {
		return false
	}

	// Verify DER-encoded signature
	return ecdsa.VerifyASN1(pubKey.ToECDSA(), msgHash[:], signatureBytes)
}

// GetPublicKey derives public key from private key
// privateKey: Private key in hex format (with or without '0x' prefix)
// Returns: Public key in uncompressed hex format (64 bytes, without 0x04 prefix)
func (c *Client) GetPublicKey(privateKey string) (string, error) {
	// Remove 0x prefix if present
	cleanKey := hexFix(privateKey)

	// Parse private key
	privKeyBytes, err := hex.DecodeString(cleanKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	if privKey == nil {
		return "", fmt.Errorf("failed to derive public key")
	}

	// Get uncompressed public key bytes (65 bytes with 0x04 prefix)
	pubKeyBytes := pubKey.SerializeUncompressed()

	// Remove 0x04 prefix to return 64 bytes
	return hex.EncodeToString(pubKeyBytes[1:]), nil
}

// HashString computes SHA256 hash of a string
// str: String to hash
// Returns: SHA256 hash as hex string
func (c *Client) HashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

// GetFormattedTimestamp returns current timestamp in Circular Protocol format
// Format: YYYY:MM:DD-hh:mm:ss (UTC)
// Returns: Formatted timestamp string
func (c *Client) GetFormattedTimestamp() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%d:%02d:%02d-%02d:%02d:%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
}

// ============================================================================
// Helper Methods - Encoding
// ============================================================================

// HexFix normalizes hex strings (removes 0x prefix if present)
// This method is used internally for auto-preprocessing but is also exported
// for users who want to manually preprocess hex strings.
//
// Parameters:
//   - hexString: Hex string to normalize (with or without '0x' prefix)
//
// Returns:
//   - Normalized hex string without '0x' prefix
//
// Example:
//
//	normalized := client.HexFix("0xabcdef")  // returns "abcdef"
//	normalized := client.HexFix("abcdef")    // returns "abcdef"
func (c *Client) HexFix(hexString string) string {
	if len(hexString) >= 2 && (hexString[:2] == "0x" || hexString[:2] == "0X") {
		return hexString[2:]
	}
	return hexString
}

// hexFix is an internal convenience wrapper
func hexFix(hexString string) string {
	if len(hexString) >= 2 && (hexString[:2] == "0x" || hexString[:2] == "0X") {
		return hexString[2:]
	}
	return hexString
}

// StringToHex converts string to hex encoding
func (c *Client) StringToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// HexToString converts hex encoding to string
func (c *Client) HexToString(hexStr string) (string, error) {
	normalized := hexFix(hexStr)
	bytes, err := hex.DecodeString(normalized)
	if err != nil {
		return "", fmt.Errorf("invalid hex string: %w", err)
	}
	return string(bytes), nil
}


// ============================================================================
// Helper Methods - Advanced
// ============================================================================
// GetTransactionOutcome polls for transaction confirmation
// blockchain: Blockchain network (e.g., 'MainNet', 'testnet')
// txID: Transaction ID to monitor
// start: Start block number for search
// end: End block number for search
// timeoutSec: Maximum time to wait in seconds (default: 120)
// intervalSec: Polling interval in seconds (default: 5)
// Returns: Transaction response when confirmed
func (c *Client) GetTransactionOutcome(
	ctx context.Context,
	blockchain string,
	txID string,
	start string,
	end string,
	timeoutSec int,
	intervalSec int,
) (map[string]interface{}, error) {
	if timeoutSec <= 0 {
		timeoutSec = 120
	}
	if intervalSec <= 0 {
		intervalSec = 5
	}

	startTime := time.Now()
	timeout := time.Duration(timeoutSec) * time.Second
	interval := time.Duration(intervalSec) * time.Second

	for {
		// Check if timeout exceeded
		elapsed := time.Since(startTime)
		if elapsed >= timeout {
			err := fmt.Errorf("transaction %s timed out after %d seconds", txID, timeoutSec)
			return nil, err
		}

		// Check transaction status
		request := map[string]interface{}{
			"Blockchain": blockchain,
			"ID":         txID,
			"Start":      start,
			"End":        end,
			"Version":    "2.0.0-alpha.1",
		}

		tx, err := c.GetTransactionByID(ctx, request)
		if err != nil {
			// If error is not just "pending", return error
			if !strings.Contains(strings.ToLower(err.Error()), "pending") {
				return nil, err
			}

			// Otherwise, wait and retry
			time.Sleep(interval)
			continue
		}

		// Check if transaction is confirmed (has BlockNumber)
		if response, ok := tx["Response"].(map[string]interface{}); ok {
			if blockNum, ok := response["BlockNumber"].(float64); ok && blockNum > 0 {
				// Transaction confirmed
				return tx, nil
			}
		}

		// Still pending, wait before next check
		time.Sleep(interval)
	}
}