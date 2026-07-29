// Command action-test is a minimal HTTP service used to exercise the shared
// CI workflows in duynhlab/gha-workflows against a real Go build.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

var errNegative = errors.New("n must not be negative")

// fib returns the nth Fibonacci number. It exists to give the test suite
// something with real branches to cover.
func fib(n int) (int, error) {
	if n < 0 {
		return 0, errNegative
	}
	if n < 2 {
		return n, nil
	}
	prev, cur := 0, 1
	for i := 2; i <= n; i++ {
		prev, cur = cur, prev+cur
	}
	return cur, nil
}

func newMux(started time.Time) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(healthResponse{
			Status:  "ok",
			Version: version,
			Uptime:  time.Since(started).Truncate(time.Second).String(),
		}); err != nil {
			log.Printf("encode health response: %v", err)
		}
	})

	mux.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.Query().Get("n")
		n, err := strconv.Atoi(raw)
		if err != nil {
			http.Error(w, "n must be an integer", http.StatusBadRequest)
			return
		}
		result, err := fib(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%d\n", result)
	})

	return mux
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           newMux(time.Now()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("action-test %s listening on %s", version, addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
