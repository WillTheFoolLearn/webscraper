package main

import "testing"

func TestNormalizedURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "secure with slash",
			inputURL: "https://blog.boot.dev/cheese/",
			expected: "blog.boot.dev/cheese",
		},
		{
			name:     "not secure with slash",
			inputURL: "http://blog.boot.dev/robot/",
			expected: "blog.boot.dev/robot",
		},
		{
			name:     "not secure without slash",
			inputURL: "http://blog.boot.dev/frog",
			expected: "blog.boot.dev/frog",
		},
		{
			name:     "upper to lowercase",
			inputURL: "HTTP://BlOg.boot.Dev/cow",
			expected: "blog.boot.dev/cow",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
