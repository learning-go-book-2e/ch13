package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
Write a small middleware component that uses structured logging to log the
IP address of each incoming request to your web server.
*/
func main() {
	options := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stderr, options)
	mySlog := slog.New(handler)
	r := createChiRouter(mySlog)
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}
	err := s.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func createChiRouter(logger *slog.Logger) chi.Router {
	r := chi.NewRouter().With(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ip, _, _ := strings.Cut(req.RemoteAddr, ":")
			logger.Info("incoming IP", "ip", ip)
			handler.ServeHTTP(rw, req)
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().Format(time.RFC3339)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(t))
	})
	return r
}
