package monitor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPCHecker(t *testing.T) {
	tests := []struct {
		status             string
		expectedStatusCode int
		statusCode         int
		timeout            int
		sleep              int
	}{
		{expectedStatusCode: 200, statusCode: 200, status: StatusUp, timeout: 1, sleep: 0},
		{expectedStatusCode: 200, statusCode: 0, status: StatusTimeout, timeout: 1, sleep: 2},
		{expectedStatusCode: 200, statusCode: 500, status: StatusFailed, timeout: 1, sleep: 0},
	}

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, "/API/v1/healthz", req.URL.String())

			if tc.sleep > 0 {
				time.Sleep(time.Duration(tc.sleep) * time.Second)
			}

			rw.WriteHeader(tc.statusCode)
			rw.Write([]byte(`OK`))
		}))

		checker := NewHTTPCHecker(
			"test checker", fmt.Sprintf("%s/API/v1/healthz", server.URL), 1, tc.expectedStatusCode,
		)
		statusCode, status := checker.Check()
		assert.Equal(t, tc.statusCode, statusCode)
		assert.Equal(t, tc.status, status)

		server.Close()
	}
}
