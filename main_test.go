package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFib(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       int
		want    int
		wantErr error
	}{
		{name: "zero", n: 0, want: 0},
		{name: "one", n: 1, want: 1},
		{name: "two", n: 2, want: 1},
		{name: "ten", n: 10, want: 55},
		{name: "negative", n: -1, wantErr: errNegative},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := fib(tt.n)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("fib(%d) error = %v, want %v", tt.n, err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Errorf("fib(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	newMux(time.Now()).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("body = %q, want it to contain the ok status", rec.Body.String())
	}
}

func TestFibHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		query    string
		wantCode int
		wantBody string
	}{
		{name: "ok", query: "?n=10", wantCode: http.StatusOK, wantBody: "55"},
		{name: "not a number", query: "?n=abc", wantCode: http.StatusBadRequest},
		{name: "negative", query: "?n=-3", wantCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/fib"+tt.query, nil)
			newMux(time.Now()).ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want it to contain %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
