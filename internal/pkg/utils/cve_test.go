package utils

import "testing"

func TestIsCVEFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid format",
			input:    "CVE-2024-12345",
			expected: true,
		},
		{
			name:     "valid format with longer number",
			input:    "CVE-2023-123456",
			expected: true,
		},
		{
			name:     "invalid format wrong prefix",
			input:    "CVM-2024-12345",
			expected: false,
		},
		{
			name:     "invalid format lowercase",
			input:    "cve-2024-12345",
			expected: false,
		},
		{
			name:     "invalid format wrong separator",
			input:    "CVE_2024_12345",
			expected: false,
		},
		{
			name:     "invalid format short number",
			input:    "CVE-2024-123",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCVEFormat(tt.input); got != tt.expected {
				t.Errorf("IsCVEFormat(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
