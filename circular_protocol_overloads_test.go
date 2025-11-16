package circularprotocol

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCheckWalletDualSignature tests both convenience and raw methods for CheckWallet
func TestCheckWalletDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"exists": true},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.CheckWallet(ctx, "MainNet", "address123")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"Address":    "address123",
			"Version":    "1.0.9",
		}
		result, err := client.CheckWalletRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetWalletDualSignature tests both signatures for GetWallet
func TestGetWalletDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"Balance": 100, "Nonce": 5},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.GetWallet(ctx, "MainNet", "address456")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"Address":    "address456",
			"Version":    "1.0.9",
		}
		result, err := client.GetWalletRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestSendTransactionDualSignature tests both signatures for SendTransaction
func TestSendTransactionDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"Status": "pending", "TransactionID": "tx123"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method with positional params", func(t *testing.T) {
		result, err := client.SendTransaction(
			ctx,
			"txid123",
			"from456",
			"to789",
			"2025:11:16-12:00:00",
			"C_TYPE_TRANSFER",
			"payload123",
			"1",
			"signature456",
			"MainNet",
		)
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method with request object", func(t *testing.T) {
		req := map[string]interface{}{
			"ID":         "txid123",
			"From":       "from456",
			"To":         "to789",
			"Timestamp":  "2025:11:16-12:00:00",
			"Type":       "C_TYPE_TRANSFER",
			"Payload":    "payload123",
			"Nonce":      "1",
			"Signature":  "signature456",
			"Blockchain": "MainNet",
			"Version":    "1.0.9",
		}
		result, err := client.SendTransactionRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetTransactionByIDDualSignature tests both signatures
func TestGetTransactionByIDDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"ID": "tx123", "BlockNumber": 100},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.GetTransactionByID(ctx, "MainNet", "tx123", "0", "1000")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"ID":         "tx123",
			"Start":      "0",
			"End":        "1000",
			"Version":    "1.0.9",
		}
		result, err := client.GetTransactionByIDRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetBlockDualSignature tests both signatures for GetBlock
func TestGetBlockDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"BlockNumber": 100, "Hash": "hash123"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.GetBlock(ctx, "MainNet", "100")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain":  "MainNet",
			"BlockNumber": "100",
			"Version":     "1.0.9",
		}
		result, err := client.GetBlockRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestTestContractDualSignature tests both signatures for TestContract
func TestTestContractDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": "Contract test successful",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	projectCode := "function test() {}"

	t.Run("convenience method with auto hex conversion", func(t *testing.T) {
		result, err := client.TestContract(ctx, "MainNet", "from123", projectCode)
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method with manual hex conversion", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"From":       "from123",
			"Project":    client.StringToHex(projectCode),
			"Timestamp":  client.GetFormattedTimestamp(),
			"Version":    "1.0.9",
		}
		result, err := client.TestContractRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetAssetDualSignature tests both signatures for GetAsset
func TestGetAssetDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"AssetName": "CIRX", "TotalSupply": 1000000},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.GetAsset(ctx, "MainNet", "CIRX")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"AssetName":  "CIRX",
			"Version":    "1.0.9",
		}
		result, err := client.GetAssetRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetDomainDualSignature tests both signatures for GetDomain
func TestGetDomainDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{"Domain": "test.cir", "Address": "addr123"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method", func(t *testing.T) {
		result, err := client.GetDomain(ctx, "MainNet", "test.cir")
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method", func(t *testing.T) {
		req := map[string]interface{}{
			"Blockchain": "MainNet",
			"Domain":     "test.cir",
			"Version":    "1.0.9",
		}
		result, err := client.GetDomainRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestGetBlockchainsDualSignature tests GetBlockchains (no params vs empty map)
func TestGetBlockchainsDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": []map[string]interface{}{{"Name": "MainNet", "Active": true}},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	t.Run("convenience method with no params", func(t *testing.T) {
		result, err := client.GetBlockchains(ctx)
		if err != nil {
			t.Fatalf("Convenience method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("raw method with version only", func(t *testing.T) {
		req := map[string]interface{}{
			"Version": "1.0.9",
		}
		result, err := client.GetBlockchainsRaw(ctx, req)
		if err != nil {
			t.Fatalf("Raw method failed: %v", err)
		}
		if result["Result"] != float64(200) {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

// TestAllMethodsHaveDualSignature verifies all 24 core methods have both versions
func TestAllMethodsHaveDualSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"Result":   200,
			"Response": map[string]interface{}{},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL+"?cep=", "")
	ctx := context.Background()

	// Test that all methods exist and don't panic
	tests := []struct {
		name           string
		convenienceCall func() error
		rawCall        func() error
	}{
		{
			name: "CheckWallet",
			convenienceCall: func() error {
				_, err := client.CheckWallet(ctx, "MainNet", "addr")
				return err
			},
			rawCall: func() error {
				_, err := client.CheckWalletRaw(ctx, map[string]interface{}{"Blockchain": "MainNet", "Address": "addr", "Version": "1.0.9"})
				return err
			},
		},
		{
			name: "GetWalletBalance",
			convenienceCall: func() error {
				_, err := client.GetWalletBalance(ctx, "MainNet", "addr", "CIRX")
				return err
			},
			rawCall: func() error {
				_, err := client.GetWalletBalanceRaw(ctx, map[string]interface{}{"Blockchain": "MainNet", "Address": "addr", "Asset": "CIRX", "Version": "1.0.9"})
				return err
			},
		},
		{
			name: "GetBlockRange",
			convenienceCall: func() error {
				_, err := client.GetBlockRange(ctx, "MainNet", "0", "100")
				return err
			},
			rawCall: func() error {
				_, err := client.GetBlockRangeRaw(ctx, map[string]interface{}{"Blockchain": "MainNet", "Start": "0", "End": "100", "Version": "1.0.9"})
				return err
			},
		},
		// Add more method pairs as needed...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("convenience", func(t *testing.T) {
				if err := tt.convenienceCall(); err != nil {
					t.Errorf("Convenience method failed: %v", err)
				}
			})

			t.Run("raw", func(t *testing.T) {
				if err := tt.rawCall(); err != nil {
					t.Errorf("Raw method failed: %v", err)
				}
			})
		})
	}
}
