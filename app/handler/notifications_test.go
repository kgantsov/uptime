package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kgantsov/uptime/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// GET /API/v1/notifications
// ---------------------------------------------------------------------------

func TestGetNotifications(t *testing.T) {
	tests := []struct {
		name              string
		seedNotifications []model.Notification
		expectedStatus    int
		expectedCount     int
	}{
		{
			name:              "empty list",
			seedNotifications: nil,
			expectedStatus:    http.StatusOK,
			expectedCount:     0,
		},
		{
			name: "single notification",
			seedNotifications: []model.Notification{
				{Name: "slack", CallbackType: "slack", Callback: "https://hooks.slack.com/xxx"},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "multiple notifications",
			seedNotifications: []model.Notification{
				{Name: "slack", CallbackType: "slack", Callback: "https://hooks.slack.com/aaa"},
				{Name: "telegram", CallbackType: "telegram", Callback: "https://api.telegram.org/bbb", CallbackChatID: "12345"},
				{Name: "pagerduty", CallbackType: "pagerduty", Callback: "https://events.pagerduty.com/ccc"},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			for i := range tc.seedNotifications {
				require.NoError(t, db.Create(&tc.seedNotifications[i]).Error)
			}

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodGet, "/API/v1/notifications", "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var notifications []model.Notification
			require.NoError(t, json.Unmarshal(readBody(t, resp), &notifications))
			assert.Len(t, notifications, tc.expectedCount)
		})
	}
}

// ---------------------------------------------------------------------------
// GET /API/v1/notifications/:notification_name
// ---------------------------------------------------------------------------

func TestGetNotification(t *testing.T) {
	tests := []struct {
		name             string
		notifNameParam   string
		seedNotification *model.Notification
		expectedStatus   int
		expectedCallback string
	}{
		{
			name:           "existing notification",
			notifNameParam: "slack",
			seedNotification: &model.Notification{
				Name:         "slack",
				CallbackType: "slack",
				Callback:     "https://hooks.slack.com/xxx",
			},
			expectedStatus:   http.StatusOK,
			expectedCallback: "https://hooks.slack.com/xxx",
		},
		{
			name:             "non-existing notification",
			notifNameParam:   "does-not-exist",
			seedNotification: nil,
			expectedStatus:   http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedNotification != nil {
				require.NoError(t, db.Create(tc.seedNotification).Error)
			}

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodGet, "/API/v1/notifications/"+tc.notifNameParam, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var notification model.Notification
				require.NoError(t, json.Unmarshal(readBody(t, resp), &notification))
				assert.Equal(t, tc.notifNameParam, notification.Name)
				assert.Equal(t, tc.expectedCallback, notification.Callback)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// POST /API/v1/notifications
// ---------------------------------------------------------------------------

func TestCreateNotification(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		expectedStatus   int
		expectedName     string
		expectedCallback string
	}{
		{
			name:             "valid slack notification",
			body:             `{"name":"slack","callback_type":"slack","callback":"https://hooks.slack.com/xxx"}`,
			expectedStatus:   http.StatusOK,
			expectedName:     "slack",
			expectedCallback: "https://hooks.slack.com/xxx",
		},
		{
			name:             "valid telegram notification with chat id",
			body:             `{"name":"telegram","callback_type":"telegram","callback":"https://api.telegram.org/bot123/sendMessage","callback_chat_id":"99887766"}`,
			expectedStatus:   http.StatusOK,
			expectedName:     "telegram",
			expectedCallback: "https://api.telegram.org/bot123/sendMessage",
		},
		{
			name:           "malformed JSON returns bad request",
			body:           `{"name": "slack", "callback_type":}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodPost, "/API/v1/notifications", tc.body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var notification model.Notification
				require.NoError(t, json.Unmarshal(readBody(t, resp), &notification))
				assert.Equal(t, tc.expectedName, notification.Name)
				assert.Equal(t, tc.expectedCallback, notification.Callback)

				// Verify it was persisted
				var dbNotif model.Notification
				require.NoError(t, db.Where("name = ?", tc.expectedName).First(&dbNotif).Error)
				assert.Equal(t, tc.expectedName, dbNotif.Name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// PATCH /API/v1/notifications/:notification_name
// ---------------------------------------------------------------------------

func TestUpdateNotification(t *testing.T) {
	tests := []struct {
		name             string
		notifNameParam   string
		body             string
		seedNotification *model.Notification
		expectedStatus   int
		expectedCallback string
		expectedChatID   string
	}{
		{
			name:           "update callback",
			notifNameParam: "slack",
			body:           `{"callback":"https://hooks.slack.com/NEW"}`,
			seedNotification: &model.Notification{
				Name:         "slack",
				CallbackType: "slack",
				Callback:     "https://hooks.slack.com/OLD",
			},
			expectedStatus:   http.StatusOK,
			expectedCallback: "https://hooks.slack.com/NEW",
		},
		{
			name:           "update callback_chat_id",
			notifNameParam: "telegram",
			body:           `{"callback_chat_id":"NEW_CHAT_ID"}`,
			seedNotification: &model.Notification{
				Name:           "telegram",
				CallbackType:   "telegram",
				Callback:       "https://api.telegram.org/xxx",
				CallbackChatID: "OLD_CHAT_ID",
			},
			expectedStatus: http.StatusOK,
			expectedChatID: "NEW_CHAT_ID",
		},
		{
			name:           "update callback_type",
			notifNameParam: "notif1",
			body:           `{"callback_type":"pagerduty"}`,
			seedNotification: &model.Notification{
				Name:         "notif1",
				CallbackType: "slack",
				Callback:     "https://example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:             "non-existing notification",
			notifNameParam:   "ghost",
			body:             `{"callback":"https://new.example.com"}`,
			seedNotification: nil,
			expectedStatus:   http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedNotification != nil {
				require.NoError(t, db.Create(tc.seedNotification).Error)
			}

			dispatcher := new(MockDispatcher)
			if tc.expectedStatus == http.StatusOK {
				dispatcher.On("Stop").Return()
				dispatcher.On("Start").Return()
			}

			app := newTestApp(t, db, dispatcher)
			resp := doRequest(t, app, http.MethodPatch, "/API/v1/notifications/"+tc.notifNameParam, tc.body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var notification model.Notification
				require.NoError(t, json.Unmarshal(readBody(t, resp), &notification))

				if tc.expectedCallback != "" {
					assert.Equal(t, tc.expectedCallback, notification.Callback)
				}
				if tc.expectedChatID != "" {
					assert.Equal(t, tc.expectedChatID, notification.CallbackChatID)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DELETE /API/v1/notifications/:notification_name
// ---------------------------------------------------------------------------

func TestDeleteNotification(t *testing.T) {
	tests := []struct {
		name             string
		notifNameParam   string
		seedNotification *model.Notification
		expectedStatus   int
	}{
		{
			name:           "existing notification",
			notifNameParam: "slack",
			seedNotification: &model.Notification{
				Name:         "slack",
				CallbackType: "slack",
				Callback:     "https://hooks.slack.com/xxx",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:             "non-existing notification – soft delete still succeeds",
			notifNameParam:   "ghost",
			seedNotification: nil,
			expectedStatus:   http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedNotification != nil {
				require.NoError(t, db.Create(tc.seedNotification).Error)
			}

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodDelete, "/API/v1/notifications/"+tc.notifNameParam, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.seedNotification != nil {
				// Verify the record is soft-deleted (not found without Unscoped)
				var count int64
				db.Model(&model.Notification{}).Where("name = ?", tc.notifNameParam).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
