package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/llan0/go-social/internal/auth"
	"github.com/llan0/go-social/internal/store"
	"github.com/llan0/go-social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// logger := zap.Must(zap.NewProduction()).Sugar() // for logs when testing
	mockStore := store.NewMockStorage()
	mockCacheStore := cache.NewMockStorage()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStore:    mockCacheStore,
		authenticator: testAuth,
		// config:        cfg,
	}
}
func executeRequests(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, want, got int) {
	if want != got {
		t.Errorf("expected response code %v, got %v", want, got)
	}
}
