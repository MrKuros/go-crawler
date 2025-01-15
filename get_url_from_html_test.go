package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
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
			name:     "mixed relative and nested URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<a href="/home">
			<span>Example</span>
		</a>
		<a href="/about">
			<span>About Us</span>
		</a>
		<a href="https://external.com/contact">
			<span>Contact</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://example.com/home", "https://example.com/about", "https://external.com/contact"},
		},
		{
			name:     "missing http scheme and query params",
			inputURL: "https://mywebsite.org",
			inputBody: `
<html>
	<body>
		<a href="/profile">
			<span>Profile</span>
		</a>
		<a href="/search?q=golang">
			<span>Search Golang</span>
		</a>
		<a href="https://anotherdomain.com/page?item=42">
			<span>Another Page</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://mywebsite.org/profile", "https://mywebsite.org/search?q=golang", "https://anotherdomain.com/page?item=42"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getURLsFromHTML(test.inputBody, test.inputURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
