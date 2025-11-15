package circularprotocol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/circular-protocol/circular-go/circularprotocol"
)

var (
	client     *circularprotocol.Client
	mockServer *httptest.Server
)

// TestMain sets up the mock server before running tests
func TestMain(m *testing.M) {
	// Create mock server
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock server responds with success for all requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"Result": 200,
			"Response": {
				"data": "mock_success",
				"array": ["item1", "item2"],
				"count": 10,
				"exists": true
			}
		}`))
	}))
	defer mockServer.Close()

	// Create client pointing to mock server
	client = circularprotocol.NewClient(mockServer.URL+"/", "mock-key")

	// Run tests
	m.Run()
}
func Test_check_wallet(t *testing.T) {
	// Should successfully check if wallet exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.CheckWallet(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["exists"] != true {
	t.Errorf("Expected Response.exists to be true, got %%v", result["Response"].(map[string]interface{})["exists"])
}

	t.Logf("  ✅ Should successfully check if wallet exists")
}
func Test_get_latest_transactions(t *testing.T) {
	// Should fetch recent transactions for wallet
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Limit"] = 10
request["Version"] = "1.0.8"

	result, err := client.GetLatestTransactions(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["transactions"].([]interface{}); !ok {
	t.Errorf("Expected Response.transactions to be an array")
}

	t.Logf("  ✅ Should fetch recent transactions for wallet")
}
func Test_get_wallet(t *testing.T) {
	// Should retrieve wallet details
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetWallet(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["address"] == nil {
	t.Errorf("Expected Response.address to be defined")
}

	t.Logf("  ✅ Should retrieve wallet details")
}
func Test_get_wallet_balance(t *testing.T) {
	// Should get wallet balance for specific asset
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Asset"] = "0xC123"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetWalletBalance(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["balance"] == nil {
	t.Errorf("Expected Response.balance to be defined")
}

	t.Logf("  ✅ Should get wallet balance for specific asset")
}
func Test_get_wallet_nonce(t *testing.T) {
	// Should get current wallet nonce
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetWalletNonce(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if int(result["Response"].(map[string]interface{})["nonce"].(float64)) < 0 {
	t.Errorf("Expected Response.nonce to be >= 0, got %%v", result["Response"].(map[string]interface{})["nonce"])
}

	t.Logf("  ✅ Should get current wallet nonce")
}
func Test_add_transaction(t *testing.T) {
	// Should submit a transaction to blockchain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["From"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["ID"] = "0xaabbccdd11223344"
request["Nonce"] = 1
request["Payload"] = "0x1234"
request["Signature"] = "0xsignature"
request["Timestamp"] = "1234567890"
request["To"] = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
request["Type"] = "transfer"
request["Version"] = "1.0.8"

	result, err := client.SendTransaction(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["transaction_id"] == nil {
	t.Errorf("Expected Response.transaction_id to be defined")
}

	t.Logf("  ✅ Should submit a transaction to blockchain")
}
func Test_get_pending_transaction(t *testing.T) {
	// Should get pending transactions
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetPendingTransaction(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["transactions"].([]interface{}); !ok {
	t.Errorf("Expected Response.transactions to be an array")
}

	t.Logf("  ✅ Should get pending transactions")
}
func Test_get_transaction_by_address(t *testing.T) {
	// Should get transactions by address
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetTransactionByAddress(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["transactions"].([]interface{}); !ok {
	t.Errorf("Expected Response.transactions to be an array")
}

	t.Logf("  ✅ Should get transactions by address")
}
func Test_get_transaction_by_date(t *testing.T) {
	// Should get transactions by date range
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["EndDate"] = "2024-12-31"
request["StartDate"] = "2024-01-01"
request["Version"] = "1.0.8"

	result, err := client.GetTransactionByDate(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["transactions"].([]interface{}); !ok {
	t.Errorf("Expected Response.transactions to be an array")
}

	t.Logf("  ✅ Should get transactions by date range")
}
func Test_get_transaction_by_id(t *testing.T) {
	// Should get transaction by ID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["TransactionID"] = "0xaabbccdd11223344"
request["Version"] = "1.0.8"

	result, err := client.GetTransactionByID(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["transaction"] == nil {
	t.Errorf("Expected Response.transaction to be defined")
}

	t.Logf("  ✅ Should get transaction by ID")
}
func Test_get_transaction_by_node(t *testing.T) {
	// Should get transactions by node
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Node"] = "0xnode123"
request["Version"] = "1.0.8"

	result, err := client.GetTransactionByNode(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["transactions"].([]interface{}); !ok {
	t.Errorf("Expected Response.transactions to be an array")
}

	t.Logf("  ✅ Should get transactions by node")
}
func Test_get_asset(t *testing.T) {
	// Should get asset details
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Asset"] = "0xC123"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetAsset(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["asset"] == nil {
	t.Errorf("Expected Response.asset to be defined")
}

	t.Logf("  ✅ Should get asset details")
}
func Test_get_asset_list(t *testing.T) {
	// Should get list of all assets
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetAssetList(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["assets"].([]interface{}); !ok {
	t.Errorf("Expected Response.assets to be an array")
}

	t.Logf("  ✅ Should get list of all assets")
}
func Test_get_asset_supply(t *testing.T) {
	// Should get asset supply information
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Asset"] = "0xC123"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetAssetSupply(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["total_supply"] == nil {
	t.Errorf("Expected Response.total_supply to be defined")
}

	t.Logf("  ✅ Should get asset supply information")
}
func Test_get_voucher(t *testing.T) {
	// Should get voucher details
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"
request["VoucherID"] = "0xvoucher123"

	result, err := client.GetVoucher(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["voucher"] == nil {
	t.Errorf("Expected Response.voucher to be defined")
}

	t.Logf("  ✅ Should get voucher details")
}
func Test_get_blockchains(t *testing.T) {
	// Should list supported blockchains
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Version"] = "1.0.8"

	result, err := client.GetBlockchains(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["blockchains"].([]interface{}); !ok {
	t.Errorf("Expected Response.blockchains to be an array")
}
arr, ok := result["Response"].(map[string]interface{})["blockchains"].([]interface{})
if !ok {
	t.Errorf("Expected Response.blockchains to be an array")
} else {
	found := false
	for _, item := range arr {
		if item == "MainNet" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected array to contain 'MainNet'")
	}
}

	t.Logf("  ✅ Should list supported blockchains")
}
func Test_get_analytics(t *testing.T) {
	// Should get blockchain analytics
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetAnalytics(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["analytics"] == nil {
	t.Errorf("Expected Response.analytics to be defined")
}

	t.Logf("  ✅ Should get blockchain analytics")
}
func Test_get_block(t *testing.T) {
	// Should get block by number
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Block"] = 12345
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetBlock(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["block"] == nil {
	t.Errorf("Expected Response.block to be defined")
}

	t.Logf("  ✅ Should get block by number")
}
func Test_get_block_count(t *testing.T) {
	// Should get current blockchain height
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	result, err := client.GetBlockCount(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if int(result["Response"].(map[string]interface{})["count"].(float64)) <= 0 {
	t.Errorf("Expected Response.count to be greater than 0, got %%v", result["Response"].(map[string]interface{})["count"])
}

	t.Logf("  ✅ Should get current blockchain height")
}
func Test_get_block_range(t *testing.T) {
	// Should get range of blocks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["EndBlock"] = 10010
request["StartBlock"] = 10000
request["Version"] = "1.0.8"

	result, err := client.GetBlockRange(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if _, ok := result["Response"].(map[string]interface{})["blocks"].([]interface{}); !ok {
	t.Errorf("Expected Response.blocks to be an array")
}

	t.Logf("  ✅ Should get range of blocks")
}
func Test_get_domain(t *testing.T) {
	// Should resolve domain to address
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["Domain"] = "myname.circular"
request["Version"] = "1.0.8"

	result, err := client.GetDomain(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["address"] == nil {
	t.Errorf("Expected Response.address to be defined")
}

	t.Logf("  ✅ Should resolve domain to address")
}
func Test_call_contract(t *testing.T) {
	// Should call contract method
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["ContractAddress"] = "0xcontract123"
request["Method"] = "balanceOf"
request["Parameters"] = "[
  "0xwallet123"
]"
request["Version"] = "1.0.8"

	result, err := client.CallContract(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["result"] == nil {
	t.Errorf("Expected Response.result to be defined")
}

	t.Logf("  ✅ Should call contract method")
}
func Test_test_contract(t *testing.T) {
	// Should test contract execution (dry run)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Blockchain"] = ""MainNet""
request["ContractAddress"] = "0xcontract123"
request["Method"] = "transfer"
request["Parameters"] = "[
  "0xrecipient",
  "1000"
]"
request["Version"] = "1.0.8"

	result, err := client.TestContract(ctx, request)
	if err != nil {
		t.Fatalf("API call failed: %%v", err)
	}

if result["Result"] != float64(200) {
	t.Errorf("Expected Result to be 200, got %%v", result["Result"])
}
if result["Response"].(map[string]interface{})["result"] == nil {
	t.Errorf("Expected Response.result to be defined")
}

	t.Logf("  ✅ Should test contract execution (dry run)")
}
func Test_connection_error(t *testing.T) {
	// Should handle network connection errors gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	invalidClient := circularprotocol.NewClient("http://localhost:9999", "")

	request := make(map[string]interface{})
request["Address"] = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	_, err := invalidClient.CheckWallet(ctx, request)
	if err == nil {
		t.Fatalf("Expected connection error, but call succeeded")
	}

	t.Logf("  ✅ Should handle network connection errors gracefully")
}
func Test_invalid_address(t *testing.T) {
	// Should handle invalid address gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := make(map[string]interface{})
request["Address"] = "invalid"
request["Blockchain"] = ""MainNet""
request["Version"] = "1.0.8"

	_, err := client.CheckWallet(ctx, request)
	if err == nil {
		t.Errorf("Expected API error, but call succeeded")
	}

	t.Logf("  ✅ Should handle invalid address gracefully")
}