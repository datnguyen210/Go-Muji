package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


func newTestApplication(t *testing.T) *application{
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string){
	res, err := ts.Client().Get(urlPath)
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