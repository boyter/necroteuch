//go:build integration

package integration_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/carlmjohnson/requests"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestJsonEndpoints(t *testing.T) {
	type jsonCase struct {
		Url         string
		ExistValues []string
	}

	tests := []jsonCase{
		{
			Url: "/health-check/",
			ExistValues: []string{
				"success",
				"messages",
				"time",
				"timing",
				"timing.0.timeMillis",
				"timing.0.source",
				"response.ipAddress",
				"response.memUsage",
			},
		},
		{
			Url: "/.well-known/webfinger?resource=acct:test@test.com",
			ExistValues: []string{
				"subject",
				"links",
				"links.0.rel",
				"links.0.type",
				"links.0.href",
			},
		},
	}

	for _, tt := range tests {
		var s string

		err := requests.
			URL(rootLocation + tt.Url).
			Client(&client).
			Method(http.MethodGet).
			ToString(&s).
			Fetch(context.Background())

		if err != nil {
			t.Fatalf("%v %v", err, tt)
		}

		for _, e := range tt.ExistValues {
			value := gjson.Get(s, e)

			if !value.Exists() {
				t.Errorf("want %v in %v was not found", e, s)
			}
		}
	}
}

func TestWebFingerEdgeCases(t *testing.T) {
	type EdgeCases struct {
		AccountName        string
		ExpectedStatusCode int
	}

	tests := []EdgeCases{
		{
			AccountName:        "acct:test@test.com",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			AccountName:        "acct:stuff",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			AccountName:        "acct:",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			AccountName:        "",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			AccountName:        "acct:test@test.com@something.com",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			AccountName:        "acct:" + gofakeit.Adjective(),
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			AccountName:        "acct:" + gofakeit.Vegetable(),
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequestWithContext(context.Background(),
			http.MethodGet, rootLocation+"/.well-known/webfinger?resource="+tt.AccountName, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		err = res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != tt.ExpectedStatusCode {
			t.Errorf("want %v for %v got %v", tt.ExpectedStatusCode, tt.AccountName, res.StatusCode)
		}
	}
}
