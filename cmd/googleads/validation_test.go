package main

import "testing"
import "net/http"
import "net/http/httptest"

func TestRequestValidation(t *testing.T) {
	client := NewGoogleAdsClient(Config{BaseURL: "https://googleads.googleapis.com/v24"})
	if _, err := client.Request(GoogleAdsRequest{Method: "TRACE", Path: "/customers:listAccessibleCustomers"}); err == nil {
		t.Fatal("expected unsupported method to fail")
	}
	if _, err := client.Request(GoogleAdsRequest{Method: "GET"}); err == nil {
		t.Fatal("expected missing path to fail")
	}
	if _, err := client.Request(GoogleAdsRequest{Method: "GET", Path: "customers:listAccessibleCustomers"}); err == nil {
		t.Fatal("expected relative path to fail")
	}
}

func TestGAQLValidation(t *testing.T) {
	if _, err := gaqlRequestBody(GoogleAdsGAQLInput{Query: "SELECT customer.id FROM customer"}, true); err == nil {
		t.Fatal("expected missing customer id to fail")
	}
	if _, err := gaqlRequestBody(GoogleAdsGAQLInput{CustomerID: "123"}, true); err == nil {
		t.Fatal("expected missing query to fail")
	}
	body, err := gaqlRequestBody(GoogleAdsGAQLInput{CustomerID: "123", Query: "SELECT customer.id FROM customer", PageSize: "10"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if body["query"] != "SELECT customer.id FROM customer" || body["pageSize"] != "10" {
		t.Fatalf("unexpected GAQL body: %#v", body)
	}
}

func TestRequestSetsLoginCustomerIDHeaderOnlyWhenRequested(t *testing.T) {
	receivedHeaders := make(chan string, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders <- r.Header.Get("login-customer-id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(server.Close)

	client := NewGoogleAdsClient(Config{BaseURL: server.URL})
	if _, err := client.Request(GoogleAdsRequest{Method: "GET", Path: "/without"}); err != nil {
		t.Fatal(err)
	}
	if got := <-receivedHeaders; got != "" {
		t.Fatalf("expected no login-customer-id header, got %q", got)
	}

	if _, err := client.Request(GoogleAdsRequest{Method: "GET", Path: "/with", LoginCustomerID: "1234567890"}); err != nil {
		t.Fatal(err)
	}
	if got := <-receivedHeaders; got != "1234567890" {
		t.Fatalf("expected login-customer-id header, got %q", got)
	}
}
