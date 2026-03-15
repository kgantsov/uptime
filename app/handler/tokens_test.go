package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// POST /API/v1/tokens
// ---------------------------------------------------------------------------

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		seedEmail      string
		seedPassword   string // plain-text; will be hashed before storing
		expectedStatus int
	}{
		{
			name:           "valid credentials",
			body:           `{"email":"user@example.com","password":"secret123"}`,
			seedEmail:      "user@example.com",
			seedPassword:   "secret123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wrong password",
			body:           `{"email":"user@example.com","password":"wrongpassword"}`,
			seedEmail:      "user@example.com",
			seedPassword:   "secret123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unknown email",
			body:           `{"email":"nobody@example.com","password":"secret123"}`,
			seedEmail:      "user@example.com",
			seedPassword:   "secret123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body – Huma rejects missing required fields",
			body:           `{}`,
			seedEmail:      "user@example.com",
			seedPassword:   "secret123",
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()

			// Seed the user with a bcrypt-hashed password.
			hashed, err := auth.HashPassword(tc.seedPassword)
			require.NoError(t, err)
			user := &model.User{
				FirstName: "Test",
				LastName:  "User",
				Email:     tc.seedEmail,
				Password:  hashed,
			}
			require.NoError(t, db.Create(user).Error)

			app := newTestApp(t, db, nil)
			resp := doRequest(t, app, http.MethodPost, "/API/v1/tokens", tc.body)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var token model.Token
				require.NoError(t, json.Unmarshal(readBody(t, resp), &token))
				assert.NotEmpty(t, token.Token, "JWT string should be present in the response")

				// Token record must be persisted in the database.
				var dbToken model.Token
				require.NoError(t, db.First(&dbToken, token.ID).Error)
				assert.Equal(t, user.ID, dbToken.UserID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DELETE /API/v1/tokens
// ---------------------------------------------------------------------------

func TestDeleteToken(t *testing.T) {
	tests := []struct {
		name           string
		seedToken      bool // whether to insert a Token row and build a JWT for it
		sendJWT        bool // whether to include the Authorization header
		expectedStatus int
	}{
		{
			name:           "valid token – successful deletion",
			seedToken:      true,
			sendJWT:        true,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "no JWT in request – Huma rejects missing required header with 422",
			seedToken:      false,
			sendJWT:        false,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB()

			// Seed a user so the token FK is satisfied.
			user := &model.User{FirstName: "Del", LastName: "User", Email: "del@example.com", Password: "x"}
			require.NoError(t, db.Create(user).Error)

			var tokenID uint
			if tc.seedToken {
				tok := &model.Token{UserID: user.ID}
				require.NoError(t, db.Create(tok).Error)
				tokenID = tok.ID
			}

			app := newTestApp(t, db, nil)

			var resp *http.Response
			if tc.sendJWT && tokenID != 0 {
				jwtStr := makeJWTToken(tokenID)
				resp = doRequest(t, app, http.MethodDelete, "/API/v1/tokens", "", bearerHeader(jwtStr))
			} else {
				resp = doRequest(t, app, http.MethodDelete, "/API/v1/tokens", "")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.seedToken && tokenID != 0 && tc.expectedStatus == http.StatusNoContent {
				// The token row should be soft-deleted (not visible without Unscoped).
				var count int64
				db.Model(&model.Token{}).Where("id = ?", tokenID).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
