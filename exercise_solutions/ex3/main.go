package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
Add the ability to return the time as JSON. Use the +accept+ header to control
whether JSON or text is returned (default to text).

The JSON should be structured as:

	{
	    "day_of_week": "Monday",
	    "day_of_month": 10,
	    "month": "April",
	    "year": 2023,
	    "hour": 20,
	    "minute": 15,
	    "second": 20
	}
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
		now := time.Now()
		var out string
		if r.Header.Get("Accept") == "application/json" {
			out = buildJSON(now)
		} else {
			out = buildText(now)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(out))
	})
	return r
}

func buildText(now time.Time) string {
	return now.Format(time.RFC3339)
}

/*
	  {
	    "day_of_week": "Monday",
	    "day_of_month": 10,
	    "month": "April",
	    "year": 2023,
	    "hour": 20,
	    "minute": 15,
	    "second": 20
	}
*/
func buildJSON(now time.Time) string {
	timeOut := struct {
		DayOfWeek  string `json:"day_of_week"`
		DayOfMonth int    `json:"day_of_month"`
		Month      string `json:"month"`
		Year       int    `json:"year"`
		Hour       int    `json:"hour"`
		Minute     int    `json:"minute"`
		Second     int    `json:"second"`
	}{
		DayOfWeek:  now.Weekday().String(),
		DayOfMonth: now.Day(),
		Month:      now.Month().String(),
		Year:       now.Year(),
		Hour:       now.Hour(),
		Minute:     now.Minute(),
		Second:     now.Second(),
	}
	out, _ := json.Marshal(timeOut)
	return string(out)
}
