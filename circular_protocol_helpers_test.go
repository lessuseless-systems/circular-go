package circularprotocol

import (
	"encoding/hex"
	"strings"
	"testing"
)

// TestHexFix tests the HexFix method
func TestHexFix(t *testing.T) {
	client := NewClient("", "")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase 0x prefix",
			input:    "0xabcdef123456",
			expected: "abcdef123456",
		},
		{
			name:     "uppercase 0X prefix",
			input:    "0Xabcdef123456",
			expected: "abcdef123456",
		},
		{
			name:     "mixed case 0x",
			input:    "0xABCDEF123456",
			expected: "ABCDEF123456",
		},
		{
			name:     "no prefix",
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
		{
			name:     "single character after 0x",
			input:    "0xa",
			expected: "a",
		},
		{
			name:     "address with 0x",
			input:    "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
			expected: "742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
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

// TestStringToHex tests the StringToHex method
func TestStringToHex(t *testing.T) {
	client := NewClient("", "")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: "68656c6c6f",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "68656c6c6f20776f726c64",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "numbers",
			input:    "12345",
			expected: "3132333435",
		},
		{
			name:     "special characters",
			input:    "!@#$%",
			expected: "21402324 25",
		},
		{
			name:     "JSON object",
			input:    `{"key":"value"}`,
			expected: hex.EncodeToString([]byte(`{"key":"value"}`)),
		},
		{
			name:     "function code",
			input:    "function test() { return 42; }",
			expected: hex.EncodeToString([]byte("function test() { return 42; }")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.StringToHex(tt.input)
			if result != tt.expected {
				t.Errorf("StringToHex(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHexToString tests the HexToString method
func TestHexToString(t *testing.T) {
	client := NewClient("", "")

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple hex",
			input:    "68656c6c6f",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "hex with spaces encoded",
			input:    "68656c6c6f20776f726c64",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix",
			input:    "0x68656c6c6f",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "numbers encoded",
			input:    "3132333435",
			expected: "12345",
			wantErr:  false,
		},
		{
			name:     "invalid hex (odd length)",
			input:    "abc",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid hex characters",
			input:    "zzzz",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.HexToString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("HexToString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestStringToHexRoundTrip tests encoding and decoding
func TestStringToHexRoundTrip(t *testing.T) {
	client := NewClient("", "")

	testStrings := []string{
		"hello world",
		"The quick brown fox",
		`{"Action":"CP_REGISTERWALLET","PublicKey":"abc123"}`,
		"function test() { return 42; }",
		"Special chars: !@#$%^&*()",
		"",
	}

	for _, original := range testStrings {
		t.Run(original, func(t *testing.T) {
			encoded := client.StringToHex(original)
			decoded, err := client.HexToString(encoded)
			if err != nil {
				t.Fatalf("Round trip failed: %v", err)
			}
			if decoded != original {
				t.Errorf("Round trip mismatch: got %q, want %q", decoded, original)
			}
		})
	}
}

// TestHashString tests the HashString method
func TestHashString(t *testing.T) {
	client := NewClient("", "")

	tests := []struct {
		name     string
		input    string
		expected int // expected length of hash
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: 64, // SHA256 produces 32 bytes = 64 hex characters
		},
		{
			name:     "empty string",
			input:    "",
			expected: 64,
		},
		{
			name:     "long string",
			input:    strings.Repeat("a", 1000),
			expected: 64,
		},
		{
			name:     "public key",
			input:    "04abcdef123456789abcdef123456789abcdef123456789abcdef123456789abc",
			expected: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.HashString(tt.input)
			if len(result) != tt.expected {
				t.Errorf("HashString(%q) length = %d, want %d", tt.input, len(result), tt.expected)
			}

			// Verify it's valid hex
			_, err := hex.DecodeString(result)
			if err != nil {
				t.Errorf("HashString(%q) produced invalid hex: %v", tt.input, err)
			}
		})
	}
}

// TestHashStringDeterministic tests that HashString is deterministic
func TestHashStringDeterministic(t *testing.T) {
	client := NewClient("", "")

	input := "test string for determinism"
	hash1 := client.HashString(input)
	hash2 := client.HashString(input)
	hash3 := client.HashString(input)

	if hash1 != hash2 || hash2 != hash3 {
		t.Errorf("HashString is not deterministic: %s, %s, %s", hash1, hash2, hash3)
	}
}

// TestHashStringUnique tests that different inputs produce different hashes
func TestHashStringUnique(t *testing.T) {
	client := NewClient("", "")

	inputs := []string{
		"hello",
		"hello ",
		"Hello",
		"world",
		"",
		" ",
	}

	hashes := make(map[string]string)
	for _, input := range inputs {
		hash := client.HashString(input)
		if existing, found := hashes[hash]; found {
			t.Errorf("Hash collision: %q and %q produced same hash %s", input, existing, hash)
		}
		hashes[hash] = input
	}
}

// TestGetFormattedTimestamp tests the GetFormattedTimestamp method
func TestGetFormattedTimestamp(t *testing.T) {
	client := NewClient("", "")

	timestamp := client.GetFormattedTimestamp()

	// Verify format: YYYY:MM:DD-HH:mm:ss
	parts := strings.Split(timestamp, "-")
	if len(parts) != 2 {
		t.Fatalf("Invalid timestamp format: %s (expected YYYY:MM:DD-HH:mm:ss)", timestamp)
	}

	datePart := parts[0]
	timePart := parts[1]

	// Check date part YYYY:MM:DD
	datePieces := strings.Split(datePart, ":")
	if len(datePieces) != 3 {
		t.Errorf("Invalid date format in timestamp: %s", datePart)
	}

	// Check time part HH:mm:ss
	timePieces := strings.Split(timePart, ":")
	if len(timePieces) != 3 {
		t.Errorf("Invalid time format in timestamp: %s", timePart)
	}

	// Verify year is reasonable (4 digits)
	if len(datePieces[0]) != 4 {
		t.Errorf("Year should be 4 digits: %s", datePieces[0])
	}

	// Verify month is 2 digits
	if len(datePieces[1]) != 2 {
		t.Errorf("Month should be 2 digits: %s", datePieces[1])
	}

	// Verify day is 2 digits
	if len(datePieces[2]) != 2 {
		t.Errorf("Day should be 2 digits: %s", datePieces[2])
	}

	// Test that it generates different timestamps over time
	timestamp1 := client.GetFormattedTimestamp()
	timestamp2 := client.GetFormattedTimestamp()

	// They might be the same if called in the same second, but format should be valid
	if timestamp1 != "" && timestamp2 != "" {
		t.Logf("Timestamp 1: %s", timestamp1)
		t.Logf("Timestamp 2: %s", timestamp2)
	}
}

// TestSignMessageAndVerify tests signing and verification
func TestSignMessageAndVerify(t *testing.T) {
	client := NewClient("", "")

	// Generate a private key (in real use, this would be a proper secp256k1 key)
	// For testing, we'll use a sample private key
	privateKey := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	message := "Hello, Circular Protocol!"

	// Test signing
	signature, err := client.SignMessage(message, privateKey)
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}

	if signature == "" {
		t.Error("Signature should not be empty")
	}

	// Verify signature is valid hex
	_, err = hex.DecodeString(signature)
	if err != nil {
		t.Errorf("Signature is not valid hex: %v", err)
	}

	// Test signing with 0x prefix on private key
	signature2, err := client.SignMessage(message, "0x"+privateKey)
	if err != nil {
		t.Fatalf("SignMessage with 0x prefix failed: %v", err)
	}

	// Should produce same signature regardless of 0x prefix
	if signature != signature2 {
		t.Errorf("Signatures differ with/without 0x prefix: %s vs %s", signature, signature2)
	}

	// Get public key for verification
	publicKey, err := client.GetPublicKey(privateKey)
	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}

	// Verify the signature
	valid := client.VerifySignature(publicKey, message, signature)
	if !valid {
		t.Error("Signature verification failed for valid signature")
	}

	// Test with wrong message
	invalidValid := client.VerifySignature(publicKey, "Wrong message", signature)
	if invalidValid {
		t.Error("Signature verification succeeded for wrong message")
	}
}

// TestGetPublicKey tests public key derivation
func TestGetPublicKey(t *testing.T) {
	client := NewClient("", "")

	privateKey := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	// Test without 0x prefix
	pubKey1, err := client.GetPublicKey(privateKey)
	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}

	if len(pubKey1) != 128 { // 64 bytes = 128 hex chars (uncompressed, without 0x04 prefix)
		t.Errorf("Public key length = %d, want 128", len(pubKey1))
	}

	// Test with 0x prefix
	pubKey2, err := client.GetPublicKey("0x" + privateKey)
	if err != nil {
		t.Fatalf("GetPublicKey with 0x prefix failed: %v", err)
	}

	// Should produce same public key
	if pubKey1 != pubKey2 {
		t.Errorf("Public keys differ with/without 0x prefix")
	}

	// Verify it's valid hex
	_, err = hex.DecodeString(pubKey1)
	if err != nil {
		t.Errorf("Public key is not valid hex: %v", err)
	}
}

// TestGetPublicKeyDeterministic tests that GetPublicKey is deterministic
func TestGetPublicKeyDeterministic(t *testing.T) {
	client := NewClient("", "")

	privateKey := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	pubKey1, _ := client.GetPublicKey(privateKey)
	pubKey2, _ := client.GetPublicKey(privateKey)
	pubKey3, _ := client.GetPublicKey(privateKey)

	if pubKey1 != pubKey2 || pubKey2 != pubKey3 {
		t.Error("GetPublicKey is not deterministic")
	}
}

// TestVersionConstant tests that the Version constant is correct
func TestVersionConstant(t *testing.T) {
	if Version != "1.0.9" {
		t.Errorf("Version constant = %q, want %q", Version, "1.0.9")
	}
}

// TestClientVersionField tests that client.version matches Version constant
func TestClientVersionField(t *testing.T) {
	client := NewClient("", "")

	if client.version != Version {
		t.Errorf("client.version = %q, want %q", client.version, Version)
	}

	if client.version != "1.0.9" {
		t.Errorf("client.version = %q, want %q", client.version, "1.0.9")
	}
}
