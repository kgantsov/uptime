package handler

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// GET /heartbeats/latencies
// ---------------------------------------------------------------------------

func TestGetHeartbeatsLatencies(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		seedHeartbeats []model.Heartbeat
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "no heartbeats",
			queryParams:    nil,
			seedHeartbeats: nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:        "all heartbeats – no filter",
			queryParams: nil,
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 120},
				{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: 80},
				{ServiceID: 1, Status: "failed", StatusCode: 500, ResponseTime: 300},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:        "filter by service_id",
			queryParams: map[string]string{"service_id": "1"},
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 120},
				{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: 80},
				{ServiceID: 1, Status: "failed", StatusCode: 500, ResponseTime: 300},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:        "custom size",
			queryParams: map[string]string{"size": "2"},
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 100},
				{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: 110},
				{ServiceID: 3, Status: "up", StatusCode: 200, ResponseTime: 120},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "invalid size returns bad request",
			queryParams:    map[string]string{"size": "not-a-number"},
			seedHeartbeats: nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			for i := range tc.seedHeartbeats {
				require.NoError(t, db.Create(&tc.seedHeartbeats[i]).Error)
			}

			h := newTestHandler(db, nil)
			c, rec := echoCtxJSON(http.MethodGet, "/API/v1/heartbeats/latencies", "")

			// Attach query params via the underlying request's URL
			q := c.Request().URL.Query()
			for k, v := range tc.queryParams {
				q.Set(k, v)
			}
			c.Request().URL.RawQuery = q.Encode()

			err := h.GetHeartbeatsLatencies(c)
			if tc.expectedStatus == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)

				var heartbeats []model.Heartbeat
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &heartbeats))
				assert.Len(t, heartbeats, tc.expectedCount)
			} else {
				var he *echo.HTTPError
				if assert.ErrorAs(t, err, &he) {
					assert.Equal(t, tc.expectedStatus, he.Code)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GET /heartbeats/latencies/last
// ---------------------------------------------------------------------------

func TestGetHeartbeatsLastLatencies(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		seedHeartbeats []model.Heartbeat
		expectedStatus int
		// At most (size * number-of-distinct-services) rows can be returned.
		maxExpectedCount int
	}{
		{
			name:             "no heartbeats",
			queryParams:      nil,
			seedHeartbeats:   nil,
			expectedStatus:   http.StatusOK,
			maxExpectedCount: 0,
		},
		{
			name:        "default size of 3 – two services with 5 heartbeats each",
			queryParams: nil,
			seedHeartbeats: func() []model.Heartbeat {
				var hbs []model.Heartbeat
				now := time.Now()
				for i := 0; i < 5; i++ {
					hbs = append(hbs,
						model.Heartbeat{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: int64(i * 10), CreatedAt: now.Add(time.Duration(i) * time.Second)},
						model.Heartbeat{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: int64(i * 20), CreatedAt: now.Add(time.Duration(i) * time.Second)},
					)
				}
				return hbs
			}(),
			expectedStatus:   http.StatusOK,
			maxExpectedCount: 6, // 3 per service × 2 services
		},
		{
			name:        "size=1 – one heartbeat per service",
			queryParams: map[string]string{"size": "1"},
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 100},
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 200},
				{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: 150},
			},
			expectedStatus:   http.StatusOK,
			maxExpectedCount: 2, // 1 per service × 2 services
		},
		{
			name:           "invalid size returns bad request",
			queryParams:    map[string]string{"size": "abc"},
			seedHeartbeats: nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			for i := range tc.seedHeartbeats {
				require.NoError(t, db.Create(&tc.seedHeartbeats[i]).Error)
			}

			h := newTestHandler(db, nil)
			c, rec := echoCtxJSON(http.MethodGet, "/API/v1/heartbeats/latencies/last", "")

			q := c.Request().URL.Query()
			for k, v := range tc.queryParams {
				q.Set(k, v)
			}
			c.Request().URL.RawQuery = q.Encode()

			err := h.GetHeartbeatsLastLatencies(c)
			if tc.expectedStatus == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)

				var heartbeats []model.Heartbeat
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &heartbeats))
				assert.LessOrEqual(t, len(heartbeats), tc.maxExpectedCount)
			} else {
				var he *echo.HTTPError
				if assert.ErrorAs(t, err, &he) {
					assert.Equal(t, tc.expectedStatus, he.Code)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GET /heartbeats/stats/:days
// ---------------------------------------------------------------------------

func TestGetHeartbeatStats(t *testing.T) {
	tests := []struct {
		name           string
		daysParam      string
		seedHeartbeats []model.Heartbeat
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "no heartbeats in range",
			daysParam:      "7",
			seedHeartbeats: nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:      "heartbeats within last 7 days",
			daysParam: "7",
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 100, CreatedAt: time.Now().Add(-24 * time.Hour)},
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 110, CreatedAt: time.Now().Add(-48 * time.Hour)},
				{ServiceID: 1, Status: "failed", StatusCode: 500, ResponseTime: 300, CreatedAt: time.Now().Add(-12 * time.Hour)},
				{ServiceID: 2, Status: "up", StatusCode: 200, ResponseTime: 90, CreatedAt: time.Now().Add(-6 * time.Hour)},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3, // (service1, up), (service1, failed), (service2, up)
		},
		{
			name:      "heartbeats outside the window are excluded",
			daysParam: "1",
			seedHeartbeats: []model.Heartbeat{
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 100, CreatedAt: time.Now().Add(-2 * time.Hour)},
				{ServiceID: 1, Status: "up", StatusCode: 200, ResponseTime: 110, CreatedAt: time.Now().Add(-240 * time.Hour)}, // 10 days ago
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1, // only the heartbeat within the last day
		},
		{
			name:           "invalid days param",
			daysParam:      "not-a-number",
			seedHeartbeats: nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			for i := range tc.seedHeartbeats {
				require.NoError(t, db.Create(&tc.seedHeartbeats[i]).Error)
			}

			h := newTestHandler(db, nil)
			c, rec := echoCtxJSON(http.MethodGet, "/API/v1/heartbeats/stats/"+tc.daysParam, "")
			c.SetParamNames("days")
			c.SetParamValues(tc.daysParam)

			err := h.GetHeartbeatStats(c)
			if tc.expectedStatus == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)

				var stats []model.HeartbeatStatsPoint
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &stats))
				assert.Len(t, stats, tc.expectedCount)
			} else {
				var he *echo.HTTPError
				if assert.ErrorAs(t, err, &he) {
					assert.Equal(t, tc.expectedStatus, he.Code)
				}
			}
		})
	}
}
