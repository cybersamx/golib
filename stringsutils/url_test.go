package stringsutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/stringsutils"
)

func TestMaskPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{
			input: "postgres://pguser:password@localhost:5433/db_test?sslmode=disable",
			want:  "postgres://pguser:***@localhost:5433/db_test?sslmode=disable",
		},
		{
			input: "mysql://mysqluser:password@localhost:3307/db_test?ssl-mode=disabled",
			want:  "mysql://mysqluser:***@localhost:3307/db_test?ssl-mode=disabled",
		},
		{
			input: "http://user:password@localhost/path",
			want:  "http://user:***@localhost/path",
		},
		{
			input: "http://localhost/path",
			want:  "http://localhost/path",
		},
	}

	for _, test := range tests {
		masked := MaskPassword(test.input)

		assert.Equal(t, test.want, masked)
	}
}

func TestMapParams(t *testing.T) {
	tests := []struct {
		path   string
		params map[string]string
		want   string
	}{
		{
			path: "",
			want: "",
		},
		{
			path: ":",
			want: ":",
		},
		{
			path: "/:",
			want: "/:",
		},
		{
			path:   ":name",
			params: map[string]string{"name": "david"},
			want:   "david",
		},
		{
			path:   ":name/",
			params: map[string]string{"name": "david"},
			want:   "david/",
		},
		{
			path:   "/:name",
			params: map[string]string{"name": "david"},
			want:   "/david",
		},
		{
			path:   "/:name:",
			params: map[string]string{"name:": "david"},
			want:   "/david",
		},
		{
			path:   "/::name:",
			params: map[string]string{":name:": "david"},
			want:   "/david",
		},
		{
			path:   "/users/:name",
			params: map[string]string{"name": "david"},
			want:   "/users/david",
		},
		{
			path:   "/users/:name/",
			params: map[string]string{"name": "david"},
			want:   "/users/david/",
		},
		{
			path:   "/areas/:zip/people/:name",
			params: map[string]string{"zip": "90405", "name": "nancy"},
			want:   "/areas/90405/people/nancy",
		},
		{
			path:   "/areas/:zip/people/:lastname/:firstname",
			params: map[string]string{"zip": "90405", "lastname": "mclean", "firstname": "jon"},
			want:   "/areas/90405/people/mclean/jon",
		},
	}

	for _, test := range tests {
		sub := MapParams(test.path, test.params)
		assert.Equal(t, test.want, sub)
	}
}

func TestTrimFileExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		filePath string
		want     string
	}{
		{"image.png", "image"},
		{"image.x", "image"},
		{"image.x.y.z", "image.x.y"},
		{"image.", "image"},
		{"/absolute/path/image.png", "/absolute/path/image"},
		{"relative/path/image.png", "relative/path/image"},
		{"", ""},
	}

	for _, test := range tests {
		trimmed := TrimFileExtension(test.filePath)
		assert.Equal(t, test.want, trimmed)
	}
}

func TestGetPathFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		url  string
		want string
	}{
		{"/", "/"},
		{"https://example.com", "/"},
		{"https://example.com/", "/"},
		{"/products/books/isbn-123456", "/products/books/isbn-123456"},
		{"https://example.com/products/books/isbn-123456", "/products/books/isbn-123456"},
		{"https://example.com/products/xyz?ref=123&x=fb", "/products/xyz"},
		{"https://example.com:443/products/xyz?ref=123&x=fb", "/products/xyz"},
	}

	for _, test := range tests {
		result := GetPathFromURL(test.url)
		assert.Equal(t, test.want, result)
	}
}
