// Package circularprotocol provides a Go SDK for the Circular Protocol blockchain API.
// Generated from Nickel API specification
// Version: 1.0.8
//
// Example usage:
//
//	client := circularprotocol.NewClient("https://nag.circularlabs.io/NAG.php?cep=", "")
//	result, err := client.CheckWallet(context.Background(), map[string]interface{}{
//		"Address":    "0x...",
//		"Blockchain": "714d2ac07a826b66ac56752eebd7c77b58d2ee842e523d913fd0ef06e6bdfcae",
//		"Version":    "1.0.8",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
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

// Client represents a Circular Protocol API client
type Client struct {
	nagURL     string
	nagKey     string
	httpClient *http.Client
	headers    map[string]string
	lastError string
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
// API Methods
// ============================================================================
// CheckWallet Check if wallet exists
// Checks whether a wallet address exists on the specified blockchain.
Returns existence status and confirms the address format.
func (c *Client) CheckWallet(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "CheckWallet", req)
}
// GetWallet Get wallet information
// Retrieves complete wallet information including balance and nonce.
Returns all wallet properties including current state on the blockchain.
func (c *Client) GetWallet(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWallet", req)
}
// GetLatestTransactions Get latest transactions for wallet
// Retrieves the latest transactions for a wallet address.
Returns an array of transaction objects with details.
func (c *Client) GetLatestTransactions(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetLatestTransactions", req)
}
// GetWalletBalance Get wallet balance for specific asset
// Retrieves the balance of a specified asset in a wallet.
Returns the balance amount for the requested asset.
func (c *Client) GetWalletBalance(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWalletBalance", req)
}
// GetWalletNonce Get wallet nonce
// Retrieves the nonce (transaction counter) of a wallet.
The nonce is used for transaction ordering and must increment with each transaction.
func (c *Client) GetWalletNonce(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetWalletNonce", req)
}
// SendTransaction Submit transaction to blockchain
// Submits a transaction to the blockchain. Requires a complete signed transaction
// including ID, addresses, payload, nonce, and signature.
func (c *Client) SendTransaction(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "AddTransaction", req)
}
// GetPendingTransaction Get pending transaction by ID
// Searches for a transaction by ID among pending transactions.
Returns the transaction if it exists and is still pending.
func (c *Client) GetPendingTransaction(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetPendingTransaction", req)
}
// GetTransactionByID Find transaction by ID
// Finds a transaction by ID within a specified block range.
// Searches through blocks to locate the transaction.
func (c *Client) GetTransactionByID(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyID", req)
}
// GetTransactionByNode Find transactions by node ID
// Finds transactions by node ID within a specified block range.
// Returns all transactions associated with the node.
func (c *Client) GetTransactionByNode(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyNode", req)
}
// GetTransactionByAddress Find transactions by address
// Finds transactions by wallet address within a specified block range.
// Returns transactions where the address is sender or recipient.
func (c *Client) GetTransactionByAddress(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyAddress", req)
}
// GetTransactionByDate Find transactions by date range
// Finds transactions by wallet address within a specified date range.
// Returns all transactions for the address between the dates.
func (c *Client) GetTransactionByDate(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetTransactionbyDate", req)
}
// GetBlock Get specific block
// Retrieves a desired block by block number.
Returns complete block information including transactions and hash.
func (c *Client) GetBlock(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlock", req)
}
// GetBlockRange Get range of blocks
// Retrieves all blocks in a specified range.
If End = 0, then Start is the number of blocks from the last one minted going backward.
func (c *Client) GetBlockRange(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockRange", req)
}
// GetBlockCount Get blockchain height
// Retrieves the blockchain block height (total number of blocks).
Also known as getBlockHeight in some documentation.
func (c *Client) GetBlockCount(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockCount", req)
}
// GetAnalytics Get blockchain analytics
// Retrieves blockchain analytics and statistics.
Returns comprehensive information about the blockchain state.
func (c *Client) GetAnalytics(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAnalytics", req)
}
// TestContract Test smart contract execution
// Tests smart contract execution locally without sending a transaction.
Useful for testing contract logic before deploying or executing.
func (c *Client) TestContract(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "TestContract", req)
}
// CallContract Call smart contract function
// Calls a smart contract function on the blockchain.
Executes the specified function with provided parameters.
func (c *Client) CallContract(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "CallContract", req)
}
// GetAssetList List all assets on blockchain
// Retrieves the list of all assets minted on a specific blockchain.
Returns an array of asset information.
func (c *Client) GetAssetList(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAssetList", req)
}
// GetAsset Get specific asset information
// Retrieves an asset descriptor with complete asset information.
Returns detailed information about the specified asset.
func (c *Client) GetAsset(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAsset", req)
}
// GetAssetSupply Get asset supply information
// Retrieves the total, circulating, and residual supply of a specified asset.
Returns comprehensive supply metrics.
func (c *Client) GetAssetSupply(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetAssetSupply", req)
}
// GetVoucher Retrieve voucher information
// Retrieves an existing voucher by code.
Code is automatically stripped of 0x prefix if present.
func (c *Client) GetVoucher(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetVoucher", req)
}
// GetDomain Resolve domain to wallet address
// Resolves a domain name to a wallet address.
A single wallet can have multiple domain associations.
Also known as resolveDomain.
func (c *Client) GetDomain(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetDomain", req)
}
// GetBlockchains List available blockchains
// Retrieves the list of blockchains available in the network.
Returns information about all active and inactive blockchains.
func (c *Client) GetBlockchains(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	return c.makeRequest(ctx, "GetBlockchains", req)
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

	// Build request
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

	// Call SendTransaction
	return c.SendTransaction(ctx, request)
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
// hexFix normalizes hex strings (removes 0x prefix if present)
// This is a package-level helper function
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

// padNumber pads number with leading zero if single digit
// This is a package-level helper function
func padNumber(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	}
	return fmt.Sprintf("%d", num)
}

// GetFormattedTimestamp returns current timestamp in Circular Protocol format
// Format: YYYY:MM:DD-hh:mm:ss (UTC)
func (c *Client) GetFormattedTimestamp() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%d:%s:%s-%s:%s:%s",
		now.Year(),
		padNumber(int(now.Month())),
		padNumber(now.Day()),
		padNumber(now.Hour()),
		padNumber(now.Minute()),
		padNumber(now.Second()))
}

// ============================================================================
// Helper Methods - Advanced
// ============================================================================
// GetError returns the last error message
func (c *Client) GetError() string {
	return c.lastError
}

// handleError stores error message for later retrieval
func (c *Client) handleError(err error) {
	if err != nil {
		c.lastError = err.Error()
	} else {
		c.lastError = "Unknown error"
	}
}

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
			c.handleError(err)
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
				c.handleError(err)
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