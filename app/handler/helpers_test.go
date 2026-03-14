package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt"
	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---------------------------------------------------------------------------
// MockDispatcher
// ---------------------------------------------------------------------------

// MockDispatcher is a testify/mock implementation of DispatcherInterface.
type MockDispatcher struct {
	mock.Mock
}

func (m *MockDispatcher) AddService(serviceID uint) {
	m.Called(serviceID)
}

func (m *MockDispatcher) RemoveService(serviceID uint) {
	m.Called(serviceID)
}

func (m *MockDispatcher) RestartService(serviceID uint) {
	m.Called(serviceID)
}

func (m *MockDispatcher) Start() {
	m.Called()
}

func (m *MockDispatcher) Stop() {
	m.Called()
}

// ---------------------------------------------------------------------------
// DB helpers
// ---------------------------------------------------------------------------

// newTestDB opens an in-memory SQLite database and runs all migrations.
// Each call returns a completely isolated database instance.
func newTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open test database: %v", err))
	}
	model.MigrateDB(db)
	return db
}

// ---------------------------------------------------------------------------
// Handler / Echo helpers
// ---------------------------------------------------------------------------

// newTestLogger returns a logrus logger whose level is set to Panic so that
// log output does not pollute test runs.
func newTestLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}

// newTestHandler creates a Handler wired to the supplied DB and dispatcher.
func newTestHandler(db *gorm.DB, dispatcher DispatcherInterface) *Handler {
	return NewHandler(newTestLogger(), db, dispatcher)
}

// echoCtx builds a minimal echo.Context around a plain HTTP request/recorder
// pair. body may be an empty string for requests that carry no payload.
func echoCtx(method, path, body, contentType string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if contentType != "" {
		req.Header.Set(echo.HeaderContentType, contentType)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// echoCtxJSON is a convenience wrapper that sets Content-Type to
// application/json.
func echoCtxJSON(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	return echoCtx(method, path, body, echo.MIMEApplicationJSON)
}

// ---------------------------------------------------------------------------
// JWT helpers
// ---------------------------------------------------------------------------

// makeJWTToken signs a HS256 token that contains the given tokenID and returns
// the raw token string.
func makeJWTToken(tokenID uint) string {
	t := jwt.New(jwt.SigningMethodHS256)
	claims := t.Claims.(jwt.MapClaims)
	claims["id"] = float64(tokenID)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signed, err := t.SignedString([]byte(Key))
	if err != nil {
		panic(fmt.Sprintf("failed to sign test JWT: %v", err))
	}
	return signed
}

// setJWTContext parses a signed token string and stores the resulting
// *jwt.Token in the echo.Context under the "token" key, mimicking what the
// JWT middleware does at runtime.
func setJWTContext(c echo.Context, tokenStr string) {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(Key), nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to parse test JWT: %v", err))
	}
	c.Set("token", parsed)
}
