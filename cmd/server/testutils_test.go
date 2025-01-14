package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/datnguyen210/go-blog/internal/models/mocks"
	"github.com/go-playground/form/v4"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		blogs:          &mocks.BlogModel{},
		users:          &mocks.UserModel{},
		formDecoder:    formDecoder,
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal()
	}
	bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t testing.TB, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if matches == nil {
		t.Fatal("no CRSF token found in body")
	}
	// Because the CSRF token is a base64 encoded string it will
	// potentially include the + character, and this will be escaped to
	// &#43;. So after extracting the token from the HTML we need
	// to run it through html.UnescapeString() to get the original token value.
	return html.UnescapeString(string(matches[1]))
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	res, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}
