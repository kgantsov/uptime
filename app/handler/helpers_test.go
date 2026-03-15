package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
	"github.com/kgantsov/uptime/app/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---------------------------------------------------------------------------
// MockDispatcher
// ---------------------------------------------------------------------------

// MockDispatcher is a testify/mock implementation of service.DispatcherInterface.
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
// Handler / Fiber helpers
// ---------------------------------------------------------------------------

// newTestLogger returns a logrus logger whose level is set to Panic so that
// log output does not pollute test runs.
func newTestLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}

// newTestHandler creates a Handler wired to the supplied DB and dispatcher by
// building the full repository → service → handler stack.
func newTestHandler(db *gorm.DB, dispatcher service.DispatcherInterface) *Handler {
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)

	heartbeatSvc := service.NewHeartbeatService(heartbeatRepo)
	serviceSvc := service.NewServiceService(serviceRepo, notifRepo, dispatcher)
	notifSvc := service.NewNotificationService(notifRepo, dispatcher)
	tokenSvc := service.NewTokenService(userRepo, tokenRepo, Key)

	return NewHandler(newTestLogger(), heartbeatSvc, serviceSvc, notifSvc, tokenSvc)
}

// newTestApp creates a bare Fiber+Huma application with routes registered but
// without any auth middleware, making it suitable for handler-level unit tests.
func newTestApp(t *testing.T, db *gorm.DB, dispatcher service.DispatcherInterface) *fiber.App {
	t.Helper()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := humafiber.New(app, huma.DefaultConfig("Test API", "1.0.0"))

	h := newTestHandler(db, dispatcher)
	h.RegisterRoutes(api)

	return app
}

// doRequest sends an HTTP request to the test Fiber app and returns the
// *http.Response.  body may be an empty string for requests without a payload.
func doRequest(t *testing.T, app *fiber.App, method, path, body string, headers ...map[string]string) *http.Response {
	t.Helper()

	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

// readBody is a convenience helper that reads and returns the entire response
// body as a byte slice.
func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return b
}

// ---------------------------------------------------------------------------
// JWT helpers
// ---------------------------------------------------------------------------

// makeJWTToken signs a HS256 token that contains the given tokenID and returns
// the raw token string.
func makeJWTToken(tokenID uint) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  float64(tokenID),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := t.SignedString([]byte(Key))
	if err != nil {
		panic(fmt.Sprintf("failed to sign test JWT: %v", err))
	}
	return signed
}

// bearerHeader returns an Authorization header map for use with doRequest.
func bearerHeader(tokenStr string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + tokenStr,
	}
}
