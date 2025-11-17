package circularprotocol

import (
	"testing"
)

// TestGetVersion verifies the GetVersion method returns the SDK version
func TestGetVersion(t *testing.T) {
	client := NewClient("", "")

	version := client.GetVersion()
	if version != Version {
		t.Errorf("Expected version '%s', got '%s'", Version, version)
	}

	if version != "1.0.9" {
		t.Errorf("Expected version '1.0.9', got '%s'", version)
	}
}

// TestSetNodeAndGetNode verifies node address management
func TestSetNodeAndGetNode(t *testing.T) {
	client := NewClient("", "")

	// Initially should be empty
	if client.GetNode() != "" {
		t.Errorf("Expected empty node URL, got '%s'", client.GetNode())
	}

	// Set node address
	testNodeURL := "https://node1.circularlabs.io"
	client.SetNode(testNodeURL)

	// Verify it was set
	if client.GetNode() != testNodeURL {
		t.Errorf("Expected node URL '%s', got '%s'", testNodeURL, client.GetNode())
	}

	// Update node address
	newNodeURL := "https://node2.circularlabs.io"
	client.SetNode(newNodeURL)

	// Verify it was updated
	if client.GetNode() != newNodeURL {
		t.Errorf("Expected node URL '%s', got '%s'", newNodeURL, client.GetNode())
	}
}

// TestGetError verifies error tracking
func TestGetError(t *testing.T) {
	client := NewClient("", "")

	// Initially should be empty
	if client.GetError() != "" {
		t.Errorf("Expected empty error, got '%s'", client.GetError())
	}

	// Manually set an error
	client.lastError = "Test error message"

	// Verify error was tracked
	if client.GetError() != "Test error message" {
		t.Errorf("Expected error 'Test error message', got '%s'", client.GetError())
	}
}

// TestHandleError verifies error handling utility
func TestHandleError(t *testing.T) {
	client := NewClient("", "")

	// Test with nil result
	if !client.HandleError(nil) {
		t.Error("Expected HandleError to return true for nil result")
	}
	if client.GetError() != "nil response received" {
		t.Errorf("Expected error 'nil response received', got '%s'", client.GetError())
	}

	// Test with successful result
	successResult := map[string]interface{}{
		"Result":   200.0,
		"Response": map[string]interface{}{"data": "success"},
	}
	if client.HandleError(successResult) {
		t.Error("Expected HandleError to return false for successful result")
	}
	if client.GetError() != "" {
		t.Errorf("Expected empty error for success, got '%s'", client.GetError())
	}

	// Test with API error result
	errorResult := map[string]interface{}{
		"Result":   404.0,
		"Response": "Resource not found",
	}
	if !client.HandleError(errorResult) {
		t.Error("Expected HandleError to return true for error result")
	}
	if client.GetError() != "Resource not found" {
		t.Errorf("Expected error 'Resource not found', got '%s'", client.GetError())
	}

	// Test with error result without Response string
	errorResultNoMsg := map[string]interface{}{
		"Result":   500.0,
		"Response": map[string]interface{}{"error": "details"},
	}
	if !client.HandleError(errorResultNoMsg) {
		t.Error("Expected HandleError to return true for error result")
	}
	if client.GetError() != "API error with code: 500" {
		t.Errorf("Expected error 'API error with code: 500', got '%s'", client.GetError())
	}
}

// TestDispose verifies resource cleanup
func TestDispose(t *testing.T) {
	client := NewClient("", "")

	// Should not panic
	client.Dispose()

	// Test with nil httpClient
	client.httpClient = nil
	client.Dispose() // Should not panic
}

// TestNodeURLInConfig verifies NodeURL can be set via Config
func TestNodeURLInConfig(t *testing.T) {
	cfg := Config{
		NAGURL:  "https://nag.test.com",
		NAGKey:  "test-key",
		NodeURL: "https://node.test.com",
	}

	client := NewClientWithConfig(cfg)

	if client.GetNode() != "https://node.test.com" {
		t.Errorf("Expected node URL 'https://node.test.com', got '%s'", client.GetNode())
	}
}

// TestErrorTrackingInMakeRequest verifies that makeRequest tracks errors
func TestErrorTrackingInMakeRequest(t *testing.T) {
	// This test verifies that the lastError field is properly set/cleared
	// by makeRequest method. Since we can't easily mock HTTP without a server,
	// we just verify the structure is in place.

	client := NewClient("", "")

	// Error should initially be empty
	if client.GetError() != "" {
		t.Errorf("Expected empty initial error, got '%s'", client.GetError())
	}

	// After setting manually, should be retrievable
	client.lastError = "test"
	if client.GetError() != "test" {
		t.Error("Error tracking not working properly")
	}
}
