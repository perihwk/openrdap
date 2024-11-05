package bootstrap

import (
	"net/url"
	"testing"
)

// TestParseURLs tests the parseURLs function
func TestParseURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		hasError bool
	}{
		{
			name:     "Valid URLs",
			input:    []string{"https://example.com", "http://test.org"},
			expected: []string{"https://example.com", "http://test.org"},
			hasError: false,
		},
		{
			name:     "Invalid URLs",
			input:    []string{"https://example.com", "invalid.com"},
			expected: []string{"https://example.com", "invalid.com"},
			hasError: true,
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseURLs(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("Expected error: %v, got: %v", tt.hasError, err != nil)
			}
			if err == nil {
				urlStrings := urlsToStrings(result)
				for i, r := range urlStrings {
					if r != tt.expected[i] {
						t.Errorf("Expected %s, got %s", tt.expected[i], r)
					}
				}
			}
		})
	}
}

// TestUrlsToStrings tests the urlsToStrings function
func TestUrlsToStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []*url.URL
		expected []string
	}{
		{
			name: "Valid URLs",
			input: []*url.URL{
				{Scheme: "https", Host: "example.com"},
				{Scheme: "http", Host: "test.org"},
			},
			expected: []string{"https://example.com", "http://test.org"},
		},
		{
			name:     "Empty slice",
			input:    []*url.URL{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlsToStrings(tt.input)
			for i, r := range result {
				if r != tt.expected[i] {
					t.Errorf("Expected %s, got %s", tt.expected[i], r)
				}
			}
		})
	}
}

// TestParseASN tests the parseASN function
func TestParseASN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		hasError bool
	}{
		{
			name:     "Valid ASN without 'AS' prefix",
			input:    "12345",
			expected: 12345,
			hasError: false,
		},
		{
			name:     "Valid ASN with 'AS' prefix",
			input:    "AS12345",
			expected: 12345,
			hasError: false,
		},
		{
			name:     "Invalid ASN",
			input:    "ASABC",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseASN(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("Expected error: %v, got: %v", tt.hasError, err != nil)
			}
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
