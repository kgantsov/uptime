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
// GET /API/v1/services
// ---------------------------------------------------------------------------

func TestGetServices(t *testing.T) {
	tests := []struct {
		name           string
		seedServices   []model.Service
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "empty list",
			seedServices:   nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "single service",
			seedServices: []model.Service{
				{Name: "alpha", URL: "https://alpha.example.com", Enabled: true, CheckInterval: 60},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "multiple services",
			seedServices: []model.Service{
				{Name: "alpha", URL: "https://alpha.example.com", Enabled: true, CheckInterval: 60},
				{Name: "beta", URL: "https://beta.example.com", Enabled: false, CheckInterval: 30},
				{Name: "gamma", URL: "https://gamma.example.com", Enabled: true, CheckInterval: 120},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			for i := range tc.seedServices {
				require.NoError(t, db.Create(&tc.seedServices[i]).Error)
			}

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodGet, "/API/v1/services", "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var services []model.Service
			require.NoError(t, json.Unmarshal(readBody(t, resp), &services))
			assert.Len(t, services, tc.expectedCount)
		})
	}
}

// ---------------------------------------------------------------------------
// GET /API/v1/services/:service_id
// ---------------------------------------------------------------------------

func TestGetService(t *testing.T) {
	tests := []struct {
		name           string
		serviceIDParam string
		seedService    *model.Service
		expectedStatus int
	}{
		{
			name:           "existing service",
			serviceIDParam: "1",
			seedService:    &model.Service{Name: "alpha", URL: "https://alpha.example.com", Enabled: true},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing service",
			serviceIDParam: "999",
			seedService:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid id – not a number",
			serviceIDParam: "abc",
			seedService:    nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedService != nil {
				require.NoError(t, db.Create(tc.seedService).Error)
			}

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodGet, "/API/v1/services/"+tc.serviceIDParam, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var svc model.Service
				require.NoError(t, json.Unmarshal(readBody(t, resp), &svc))
				assert.Equal(t, tc.seedService.Name, svc.Name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// POST /API/v1/services
// ---------------------------------------------------------------------------

func TestCreateService(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedName   string
	}{
		{
			name:           "valid service",
			body:           `{"name":"svc1","url":"https://svc1.example.com","enabled":true,"check_interval":60,"timeout":5,"retries":3,"accepted_status_code":200}`,
			expectedStatus: http.StatusOK,
			expectedName:   "svc1",
		},
		{
			name:           "minimal valid service",
			body:           `{"name":"minimal","url":"https://minimal.example.com"}`,
			expectedStatus: http.StatusOK,
			expectedName:   "minimal",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			dispatcher := new(MockDispatcher)
			dispatcher.On("AddService", uint(1)).Return()

			app := newTestApp(t, db, dispatcher)
			resp := doRequest(t, app, http.MethodPost, "/API/v1/services", tc.body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var svc model.Service
			require.NoError(t, json.Unmarshal(readBody(t, resp), &svc))
			assert.Equal(t, tc.expectedName, svc.Name)

			dispatcher.AssertCalled(t, "AddService", svc.ID)
		})
	}
}

// ---------------------------------------------------------------------------
// PATCH /API/v1/services/:service_id
// ---------------------------------------------------------------------------

func TestUpdateService(t *testing.T) {
	tests := []struct {
		name           string
		serviceIDParam string
		body           string
		seedService    *model.Service
		expectedStatus int
		expectedName   string
	}{
		{
			name:           "update name",
			serviceIDParam: "1",
			body:           `{"name":"updated-name"}`,
			seedService:    &model.Service{Name: "original", URL: "https://orig.example.com", Enabled: true},
			expectedStatus: http.StatusOK,
			expectedName:   "updated-name",
		},
		{
			name:           "update url",
			serviceIDParam: "1",
			body:           `{"url":"https://new-url.example.com"}`,
			seedService:    &model.Service{Name: "svc", URL: "https://old.example.com", Enabled: true},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "update enabled flag",
			serviceIDParam: "1",
			body:           `{"enabled":false}`,
			seedService:    &model.Service{Name: "svc", URL: "https://svc.example.com", Enabled: true},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing service",
			serviceIDParam: "999",
			body:           `{"name":"ghost"}`,
			seedService:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid id",
			serviceIDParam: "not-a-number",
			body:           `{"name":"x"}`,
			seedService:    nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedService != nil {
				require.NoError(t, db.Create(tc.seedService).Error)
			}

			dispatcher := new(MockDispatcher)
			if tc.seedService != nil {
				dispatcher.On("RestartService", tc.seedService.ID).Return()
			}

			app := newTestApp(t, db, dispatcher)
			resp := doRequest(t, app, http.MethodPatch, "/API/v1/services/"+tc.serviceIDParam, tc.body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK && tc.expectedName != "" {
				var update model.UpdateService
				require.NoError(t, json.Unmarshal(readBody(t, resp), &update))
				assert.Equal(t, tc.expectedName, *update.Name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DELETE /API/v1/services/:service_id
// ---------------------------------------------------------------------------

func TestDeleteService(t *testing.T) {
	tests := []struct {
		name           string
		serviceIDParam string
		seedService    *model.Service
		seedHeartbeats bool
		expectedStatus int
	}{
		{
			name:           "existing service",
			serviceIDParam: "1",
			seedService:    &model.Service{Name: "to-delete", URL: "https://del.example.com"},
			expectedStatus: http.StatusNoContent,
			seedHeartbeats: true,
		},
		{
			name:           "invalid id",
			serviceIDParam: "xyz",
			seedService:    nil,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedService != nil {
				require.NoError(t, db.Create(tc.seedService).Error)
			}

			if tc.seedHeartbeats && tc.seedService != nil {
				heartbeats := []model.Heartbeat{
					{ServiceID: tc.seedService.ID, Status: "up", ResponseTime: 100},
					{ServiceID: tc.seedService.ID, Status: "up", ResponseTime: 120},
				}
				for i := range heartbeats {
					require.NoError(t, db.Create(&heartbeats[i]).Error)
				}

				notif := &model.Notification{Name: "slack", CallbackType: "slack", Callback: "https://hooks.slack.com/xxx"}
				require.NoError(t, db.Create(notif).Error)
				sn := &model.ServiceNotification{
					ServiceID:        int(tc.seedService.ID),
					NotificationName: notif.Name,
				}
				require.NoError(t, db.Create(sn).Error)
			}

			dispatcher := new(MockDispatcher)
			if tc.seedService != nil {
				dispatcher.On("RemoveService", tc.seedService.ID).Return()
			}

			app := newTestApp(t, db, dispatcher)
			resp := doRequest(t, app, http.MethodDelete, "/API/v1/services/"+tc.serviceIDParam, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusNoContent {
				dispatcher.AssertCalled(t, "RemoveService", tc.seedService.ID)

				var svc model.Service
				dbErr := db.First(&svc, tc.seedService.ID).Error
				assert.Error(t, dbErr, "service should be deleted from the database")

				var heartbeatCount int64
				db.Model(&model.Heartbeat{}).Where("service_id = ?", tc.seedService.ID).Count(&heartbeatCount)
				assert.Equal(t, int64(0), heartbeatCount, "heartbeats should be deleted when service is deleted")

				var snCount int64
				db.Model(&model.ServiceNotification{}).Where("service_id = ?", tc.seedService.ID).Count(&snCount)
				assert.Equal(t, int64(0), snCount, "service notifications should be deleted when service is deleted")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// POST /API/v1/services/:service_id/notifications/:notification_name
// ---------------------------------------------------------------------------

func TestServiceAddNotification(t *testing.T) {
	tests := []struct {
		name             string
		serviceIDParam   string
		notifNameParam   string
		seedService      *model.Service
		seedNotification *model.Notification
		expectedStatus   int
	}{
		{
			name:             "success",
			serviceIDParam:   "1",
			notifNameParam:   "slack",
			seedService:      &model.Service{Name: "svc", URL: "https://svc.example.com"},
			seedNotification: &model.Notification{Name: "slack", CallbackType: "slack", Callback: "https://hooks.slack.com/xxx"},
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "service not found",
			serviceIDParam:   "999",
			notifNameParam:   "slack",
			seedService:      nil,
			seedNotification: &model.Notification{Name: "slack", CallbackType: "slack", Callback: "https://hooks.slack.com/xxx"},
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "notification not found",
			serviceIDParam:   "1",
			notifNameParam:   "nonexistent",
			seedService:      &model.Service{Name: "svc", URL: "https://svc.example.com"},
			seedNotification: nil,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "invalid service id",
			serviceIDParam:   "bad",
			notifNameParam:   "slack",
			seedService:      nil,
			seedNotification: nil,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedService != nil {
				require.NoError(t, db.Create(tc.seedService).Error)
			}
			if tc.seedNotification != nil {
				require.NoError(t, db.Create(tc.seedNotification).Error)
			}

			app := newTestApp(t, db, nil)
			path := "/API/v1/services/" + tc.serviceIDParam + "/notifications/" + tc.notifNameParam
			resp := doRequest(t, app, http.MethodPost, path, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var sn model.ServiceNotification
				require.NoError(t, json.Unmarshal(readBody(t, resp), &sn))
				assert.Equal(t, tc.notifNameParam, sn.NotificationName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DELETE /API/v1/services/:service_id/notifications/:notification_name
// ---------------------------------------------------------------------------

func TestServiceDeleteNotification(t *testing.T) {
	tests := []struct {
		name             string
		serviceIDParam   string
		notifNameParam   string
		seedService      *model.Service
		seedNotification *model.Notification
		seedLink         bool
		expectedStatus   int
	}{
		{
			name:             "success",
			serviceIDParam:   "1",
			notifNameParam:   "slack",
			seedService:      &model.Service{Name: "svc", URL: "https://svc.example.com"},
			seedNotification: &model.Notification{Name: "slack", CallbackType: "slack"},
			seedLink:         true,
			expectedStatus:   http.StatusNoContent,
		},
		{
			name:             "invalid service id",
			serviceIDParam:   "nope",
			notifNameParam:   "slack",
			seedService:      nil,
			seedNotification: nil,
			seedLink:         false,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()
			if tc.seedService != nil {
				require.NoError(t, db.Create(tc.seedService).Error)
			}
			if tc.seedNotification != nil {
				require.NoError(t, db.Create(tc.seedNotification).Error)
			}
			if tc.seedLink && tc.seedService != nil && tc.seedNotification != nil {
				sn := &model.ServiceNotification{
					ServiceID:        int(tc.seedService.ID),
					NotificationName: tc.seedNotification.Name,
				}
				require.NoError(t, db.Create(sn).Error)
			}

			app := newTestApp(t, db, nil)
			path := "/API/v1/services/" + tc.serviceIDParam + "/notifications/" + tc.notifNameParam
			resp := doRequest(t, app, http.MethodDelete, path, "")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusNoContent {
				var count int64
				db.Model(&model.ServiceNotification{}).
					Where("service_id = ? AND notification_name = ?", tc.seedService.ID, tc.notifNameParam).
					Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
