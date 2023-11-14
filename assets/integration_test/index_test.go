//go:build integration

package integration_test

import (
	"context"
	"github.com/carlmjohnson/requests"
	"net/http"
	"testing"
	"time"
)

var rootLocation = "http://localhost:8080"
var client = http.Client{
	Timeout: 5 * time.Second,
}

func TestHttpEndpoints(t *testing.T) {
	tests := []string{
		"/",
		"/static/bootstrap.min.css",
		"/static/bootstrap.bundle.min.js",
		"/static/htmx.min.js",
	}

	for _, tt := range tests {
		var s string

		// https://github.com/carlmjohnson/requests
		err := requests.
			URL(rootLocation + tt).
			Client(&client).
			Method(http.MethodGet).
			ToString(&s).
			Fetch(context.Background())

		if err != nil {
			t.Fatalf("%v %v", tt, err)
		}
	}
}
