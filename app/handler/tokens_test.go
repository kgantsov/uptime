package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// POST /tokens
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
			name:           "empty body",
			body:           `{}`,
			seedEmail:      "user@example.com",
			seedPassword:   "secret123",
			expectedStatus: http.StatusBadRequest,
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

			h := newTestHandler(db, nil)
			c, rec := echoCtxJSON(http.MethodPost, "/API/v1/tokens", tc.body)

			err = h.CreateToken(c)
			if tc.expectedStatus == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)

				var token model.Token
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &token))
				assert.NotEmpty(t, token.Token, "JWT string should be present in the response")

				// Token record must be persisted in the database.
				var dbToken model.Token
				require.NoError(t, db.First(&dbToken, token.ID).Error)
				assert.Equal(t, user.ID, dbToken.UserID)
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
// DELETE /tokens
// ---------------------------------------------------------------------------

func TestDeleteToken(t *testing.T) {
	// NOTE: DeleteToken silently discards the error from GetCurrentTokenID and
	// relies on the JWT middleware (configured in ConfigureMiddleware) to reject
	// unauthenticated requests before the handler is reached.  In unit tests we
	// call the handler directly, bypassing middleware, so both cases below
	// result in a 204 – the difference is only whether an actual token row was
	// removed from the database.
	tests := []struct {
		name           string
		seedToken      bool // whether to insert a Token row and inject its JWT
		injectJWT      bool // whether to put a parsed JWT into the echo context
		expectedStatus int
	}{
		{
			name:           "valid token – successful deletion",
			seedToken:      true,
			injectJWT:      true,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "no JWT in context – handler is a no-op (middleware guards in production)",
			seedToken:      false,
			injectJWT:      false,
			expectedStatus: http.StatusNoContent,
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

			h := newTestHandler(db, nil)
			c, rec := echoCtxJSON(http.MethodDelete, "/API/v1/tokens", "")

			if tc.injectJWT && tokenID != 0 {
				setJWTContext(c, makeJWTToken(tokenID))
			}

			err := h.DeleteToken(c)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.seedToken && tokenID != 0 {
				// The token row should be soft-deleted.
				var count int64
				db.Model(&model.Token{}).Where("id = ?", tokenID).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
