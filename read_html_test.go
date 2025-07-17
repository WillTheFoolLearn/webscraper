package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadHTML(t *testing.T) {
	tests := []struct {
		name          string
		htmlBody      string
		rawBaseURL    string
		expected      []string
		errorContains string
	}{
		{
			name:       "absolute URL",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected: []string{"https://other.com/path/one"},
		},
		{
			name:       "relative URL",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:       "relative and absolute URL",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:       "paragraph containing a URL",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<p>This is a test that talks about the website http://www.google.com</p>
			</body>
		</html>			
		`,
			expected: nil,
		},
		{
			name:       "URLs and a paragraph containing a URL",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<p>This is a test that talks about the website http://www.google.com</p>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:       "No link",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<p>This is a test that talks about the website http://www.google.com</p>
			</body>
		</html>			
		`,
			expected: nil,
		},
		{
			name:       "relative URL",
			rawBaseURL: "://invalid",
			htmlBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected:      nil,
			errorContains: "unable to parse rawBaseURL",
		},
		{
			name:       "invalid href",
			rawBaseURL: "https://blog.boot.dev",
			htmlBody: `
		<html>
			<body>
				<a href="://invalid">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>			
		`,
			expected:      nil,
			errorContains: "unable to parse a.Val",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.htmlBody, tc.rawBaseURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error: %s, got none", i, tc.name, tc.errorContains)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
