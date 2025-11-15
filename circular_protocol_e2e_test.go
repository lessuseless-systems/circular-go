package circularprotocol_test

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/circular-protocol/circular-go"
)

// Circular Protocol Go SDK E2E Tests
// Generated from Nickel E2E test specifications
//
// These tests run against REAL NAG endpoints.
// They only execute when required environment variables are present.
//
// Required ENV vars:
// - CIRCULAR_TEST_ADDRESS: Test wallet address (must exist on blockchain)
//
// Optional ENV vars:
// - CIRCULAR_NAG_URL: NAG endpoint URL (default: https://nag.circularlabs.io/NAG.php?cep=)
// - CIRCULAR_TEST_BLOCKCHAIN: Blockchain network (default: 0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2)
// - CIRCULAR_API_KEY: Optional API key
// - CIRCULAR_E2E_TIMEOUT: Request timeout in seconds (default: 30)
//
// Run with:
//   CIRCULAR_TEST_ADDRESS=0x... go test -v -tags=e2e
//
// Or skip if ENV vars not present:
//   go test -v -tags=e2e  # Will skip all tests

var requiredEnvVars = []string{"CIRCULAR_TEST_ADDRESS"}

func TestMain(m *testing.M) {
	// Check for required environment variables
	missingVars := []string{}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		println("⏭️  Skipping E2E tests - missing required environment variables:")
		for _, v := range missingVars {
			println("   -", v)
		}
		println("\nTo run E2E tests, set:")
		println("  CIRCULAR_TEST_ADDRESS=0x... go test -v -tags=e2e")
		os.Exit(0)
	}

	nagURL := getEnvOrDefault("CIRCULAR_NAG_URL", "https://nag.circularlabs.io/NAG.php?cep=")
	testAddress := os.Getenv("CIRCULAR_TEST_ADDRESS")
	blockchain := getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2")

	println("\n🌐 Running E2E tests against:", nagURL)
	println("📍 Test address:", testAddress)
	println("⛓️  Blockchain:", blockchain)
	println()

	os.Exit(m.Run())
}

func setupClient(t *testing.T) *circularprotocol.Client {
	nagURL := getEnvOrDefault("CIRCULAR_NAG_URL", "https://nag.circularlabs.io/NAG.php?cep=")
	apiKey := os.Getenv("CIRCULAR_API_KEY")
	return circularprotocol.NewClient(nagURL, apiKey)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Wallet API E2E Tests
func TestCheck_wallet(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Address": "${CIRCULAR_TEST_ADDRESS}",
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.CheckWallet(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Check if test wallet exists on blockchain")
}
func TestGet_latest_transactions(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Address": "${CIRCULAR_TEST_ADDRESS}",
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetLatestTransactions(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Get latest transactions for wallet")
}
func TestGet_wallet(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Address": "${CIRCULAR_TEST_ADDRESS}",
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetWallet(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Retrieve wallet details from blockchain")
}
func TestGet_wallet_balance(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Address": "${CIRCULAR_TEST_ADDRESS}",
  "Asset": "CIRX",
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetWalletBalance(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Get wallet balance from blockchain")
}
func TestGet_wallet_nonce(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Address": "${CIRCULAR_TEST_ADDRESS}",
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetWalletNonce(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Get wallet nonce from blockchain")
}

// Network API E2E Tests
func TestGet_blockchains(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetBlockchains(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to equal 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["Blockchains"].([]interface{}); !ok {
	t.Errorf("Expected Response.Blockchains to be an array")
}

	t.Logf("✅ E2E: Retrieve list of available blockchains")
}

// Block API E2E Tests
func TestGet_analytics(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetAnalytics(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Get blockchain analytics and statistics")
}
func TestGet_block(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "BlockNumber": 1,
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetBlock(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Retrieve specific block by number")
}
func TestGet_block_count(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetBlockCount(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Get current block count from blockchain")
}
func TestGet_block_range(t *testing.T) {
	client := setupClient(t)

	requestJSON := `{
  "Blockchain": "${CIRCULAR_TEST_BLOCKCHAIN}",
  "EndBlock": 10,
  "StartBlock": 1,
  "Version": "1.0.8"
}`
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_ADDRESS}", os.Getenv("CIRCULAR_TEST_ADDRESS"))
	requestJSON = strings.ReplaceAll(requestJSON, "$${CIRCULAR_TEST_BLOCKCHAIN}", getEnvOrDefault("CIRCULAR_TEST_BLOCKCHAIN", "0x8a20baa40c45dc5055aeb26197c203e576ef389d9acb171bd62da11dc5ad72b2"))

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		t.Fatalf("Failed to unmarshal request: %%v", err)
	}

	result, err := client.GetBlockRange(context.Background(), request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] == nil {
	t.Errorf("Expected Result to be defined")
}

	t.Logf("✅ E2E: Retrieve range of blocks")
}