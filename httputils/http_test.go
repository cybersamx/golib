package httputils_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cybersamx/golib/timeutils"
	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/httputils"
)

func TestIsStatusCode2xx(t *testing.T) {
	t.Parallel()

	tcases := []struct {
		status int
		want   bool
	}{
		{status: 199, want: false},
		{status: 200, want: true},
		{status: 299, want: true},
		{status: 300, want: false},
	}

	for _, tc := range tcases {
		assert.Equal(t, tc.want, IsStatusCode2xx(tc.status))
	}
}

func TestWriteNoCacheHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	WriteNoCacheHeaders(w)

	result := timeutils.CompareTimeAndTimeString(
		time.Unix(0, 0),
		w.Header().Get("Expires"),
		time.Local,
	)

	assert.True(t, result)
	assert.Equal(t, "no-cache, private, max-age=0", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
}
