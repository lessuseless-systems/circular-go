package circularprotocol

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHexFixPreprocessing tests that HexFix correctly strips 0x prefix
func TestHexFixPreprocessing(t *testing.T) {
	client := NewClient("", "")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with 0x prefix lowercase",
			input:    "0xabcdef123456",
			expected: "abcdef123456",
		},
		{
			name:     "with 0X prefix uppercase",
			input:    "0Xabcdef123456",
			expected: "abcdef123456",
		},
		{
			name:     "without prefix",
			input:    "abcdef123456",
			expected: "abcdef123456",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only 0x",
			input:    "0x",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.HexFix(tt.input)
			if result != tt.expected {
				t.Errorf("HexFix(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckWalletAutoPreprocessing tests that CheckWallet auto-strips 0x and injects version
func TestCheckWalletAutoPreprocessing(t *testing.T) {
	// Create a test server that captures the request
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the request body
		if err := json.NewDecoder(r.Body).Decode(&capturedRequest); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Send a successful response
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"exists": true},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	// Call with 0x prefix - should be auto-stripped
	_, err := client.CheckWallet(ctx, "0xMainNet", "0x742d35cc6634c0532925a3b844bc9e7595f0beb")
	if err != nil {
		t.Fatalf("CheckWallet failed: %v", err)
	}

	// Verify preprocessing occurred
	if capturedRequest["Blockchain"] != "MainNet" {
		t.Errorf("Blockchain not preprocessed: got %v, want 'MainNet'", capturedRequest["Blockchain"])
	}
	if capturedRequest["Address"] != "742d35cc6634c0532925a3b844bc9e7595f0beb" {
		t.Errorf("Address not preprocessed: got %v, want '742d35cc6634c0532925a3b844bc9e7595f0beb'", capturedRequest["Address"])
	}
	if capturedRequest["Version"] != "1.0.9" {
		t.Errorf("Version not auto-injected: got %v, want '1.0.9'", capturedRequest["Version"])
	}
}

// TestGetWalletAutoPreprocessing tests GetWallet preprocessing
func TestGetWalletAutoPreprocessing(t *testing.T) {
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedRequest)
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"Balance": 100},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	_, err := client.GetWallet(ctx, "0xTestNet", "0xaddress123")
	if err != nil {
		t.Fatalf("GetWallet failed: %v", err)
	}

	if capturedRequest["Blockchain"] != "TestNet" {
		t.Errorf("Blockchain not preprocessed: got %v", capturedRequest["Blockchain"])
	}
	if capturedRequest["Address"] != "address123" {
		t.Errorf("Address not preprocessed: got %v", capturedRequest["Address"])
	}
	if capturedRequest["Version"] != "1.0.9" {
		t.Errorf("Version not auto-injected: got %v", capturedRequest["Version"])
	}
}

// TestStringToHexPreprocessing tests contract methods auto-convert strings to hex
func TestStringToHexPreprocessing(t *testing.T) {
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedRequest)
		response := map[string]interface{}{
			"Result":   200,
			"Response": "Contract tested",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	projectCode := "function test() { return 42; }"
	_, err := client.TestContract(ctx, "MainNet", "0xfrom123", projectCode)
	if err != nil {
		t.Fatalf("TestContract failed: %v", err)
	}

	// Verify string was converted to hex
	expectedHex := client.StringToHex(projectCode)
	if capturedRequest["Project"] != expectedHex {
		t.Errorf("Project not converted to hex: got %v, want %v", capturedRequest["Project"], expectedHex)
	}

	// Verify other preprocessing
	if capturedRequest["Blockchain"] != "MainNet" {
		t.Errorf("Blockchain not preprocessed: got %v", capturedRequest["Blockchain"])
	}
	if capturedRequest["From"] != "from123" {
		t.Errorf("From address not preprocessed: got %v", capturedRequest["From"])
	}
	if capturedRequest["Version"] != "1.0.9" {
		t.Errorf("Version not auto-injected: got %v", capturedRequest["Version"])
	}

	// Verify timestamp was auto-generated
	if _, ok := capturedRequest["Timestamp"].(string); !ok {
		t.Errorf("Timestamp not auto-generated")
	}
}

// TestCallContractAutoPreprocessing tests CallContract preprocessing
func TestCallContractAutoPreprocessing(t *testing.T) {
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedRequest)
		response := map[string]interface{}{
			"Result":   200,
			"Response": "Contract called",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	requestStr := "getBalance"
	_, err := client.CallContract(ctx, "0xMainNet", "0xcontract456", "0xfrom789", requestStr)
	if err != nil {
		t.Fatalf("CallContract failed: %v", err)
	}

	// Verify string was converted to hex
	expectedHex := client.StringToHex(requestStr)
	if capturedRequest["Request"] != expectedHex {
		t.Errorf("Request not converted to hex: got %v, want %v", capturedRequest["Request"], expectedHex)
	}

	// Verify hex preprocessing
	if capturedRequest["Blockchain"] != "MainNet" {
		t.Errorf("Blockchain not preprocessed: got %v", capturedRequest["Blockchain"])
	}
	if capturedRequest["Address"] != "contract456" {
		t.Errorf("Contract address not preprocessed: got %v", capturedRequest["Address"])
	}
	if capturedRequest["From"] != "from789" {
		t.Errorf("From address not preprocessed: got %v", capturedRequest["From"])
	}
}

// TestVersionAutoInjection tests that version is auto-injected in all convenience methods
func TestVersionAutoInjection(t *testing.T) {
	methods := []struct {
		name string
		call func(*Client, context.Context) error
	}{
		{
			name: "CheckWallet",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.CheckWallet(ctx, "MainNet", "address")
				return err
			},
		},
		{
			name: "GetWallet",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.GetWallet(ctx, "MainNet", "address")
				return err
			},
		},
		{
			name: "GetWalletNonce",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.GetWalletNonce(ctx, "MainNet", "address")
				return err
			},
		},
		{
			name: "GetBlockCount",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.GetBlockCount(ctx, "MainNet")
				return err
			},
		},
		{
			name: "GetAnalytics",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.GetAnalytics(ctx, "MainNet")
				return err
			},
		},
	}

	for _, tt := range methods {
		t.Run(tt.name, func(t *testing.T) {
			var capturedRequest map[string]interface{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewDecoder(r.Body).Decode(&capturedRequest)
				response := map[string]interface{}{
					"Result":   200,
					"Response": map[string]interface{}{},
				}
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			client := NewClient(server.URL+"?cep=", "")
			ctx := context.Background()

			err := tt.call(client, ctx)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}

			if capturedRequest["Version"] != "1.0.9" {
				t.Errorf("%s did not auto-inject version: got %v, want '1.0.9'", tt.name, capturedRequest["Version"])
			}
		})
	}
}

// TestNoPreprocessingInRawMethods tests that Raw methods don't preprocess
func TestNoPreprocessingInRawMethods(t *testing.T) {
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedRequest)
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"exists": true},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	// Call Raw method with 0x prefix - should NOT be stripped
	req := map[string]interface{}{
		"Blockchain": "0xMainNet",
		"Address":    "0x742d35cc",
		"Version":    "1.0.9",
	}

	_, err := client.CheckWalletRaw(ctx, req)
	if err != nil {
		t.Fatalf("CheckWalletRaw failed: %v", err)
	}

	// Verify NO preprocessing occurred
	if capturedRequest["Blockchain"] != "0xMainNet" {
		t.Errorf("Raw method preprocessed Blockchain: got %v, want '0xMainNet'", capturedRequest["Blockchain"])
	}
	if capturedRequest["Address"] != "0x742d35cc" {
		t.Errorf("Raw method preprocessed Address: got %v, want '0x742d35cc'", capturedRequest["Address"])
	}
}

// TestSendTransactionNoAutoPreprocessing verifies SendTransaction doesn't auto-preprocess transaction fields
func TestSendTransactionNoAutoPreprocessing(t *testing.T) {
	var capturedRequest map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedRequest)
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"Status": "pending"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	// SendTransaction should NOT preprocess ID, From, To, etc.
	// because these are pre-computed transaction fields
	_, err := client.SendTransaction(
		ctx,
		"txid123",            // ID
		"from456",            // From
		"to789",              // To
		"2025:11:16-12:00:00", // Timestamp
		"C_TYPE_REGISTERWALLET", // Type
		"payload",            // Payload
		"0",                  // Nonce
		"signature",          // Signature
		"MainNet",            // Blockchain
	)
	if err != nil {
		t.Fatalf("SendTransaction failed: %v", err)
	}

	// Verify fields are passed as-is (no preprocessing)
	if capturedRequest["ID"] != "txid123" {
		t.Errorf("SendTransaction preprocessed ID: got %v", capturedRequest["ID"])
	}
	if capturedRequest["From"] != "from456" {
		t.Errorf("SendTransaction preprocessed From: got %v", capturedRequest["From"])
	}
	if capturedRequest["To"] != "to789" {
		t.Errorf("SendTransaction preprocessed To: got %v", capturedRequest["To"])
	}
	if capturedRequest["Blockchain"] != "MainNet" {
		t.Errorf("SendTransaction preprocessed Blockchain: got %v", capturedRequest["Blockchain"])
	}

	// Verify version was injected
	if capturedRequest["Version"] != "1.0.9" {
		t.Errorf("Version not auto-injected: got %v", capturedRequest["Version"])
	}
}
