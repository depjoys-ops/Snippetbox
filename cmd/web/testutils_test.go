package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/depjoys-ops/Snippetbox/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {

	templateCacheMock, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sessionManagerMock := scs.New()
	sessionManagerMock.Lifetime = 12 * time.Hour
	sessionManagerMock.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetModelMock{},
		users:          &mocks.UserModelMock{},
		templateCache:  templateCacheMock,
		sessionManager: sessionManagerMock,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored and sent with subsequent requests when using this client.
	ts.Client().Jar = jar

	// Disabling redirect following for the test server client by setting a custom
	// CheckRedirect function.
	// http.ErrUseLastResponse error forces the client to immediately return
	// the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
