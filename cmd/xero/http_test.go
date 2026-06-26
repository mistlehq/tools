package main

import (
	"strings"
	"testing"
)

func TestXeroAPIPathMapsSupportedFamilies(t *testing.T) {
	testCases := []struct {
		family   string
		endpoint string
		want     string
	}{
		{family: "accounting", endpoint: "/Invoices", want: "/api.xro/2.0/Invoices"},
		{family: "files", endpoint: "/Files", want: "/files.xro/1.0/Files"},
		{family: "assets", endpoint: "/Assets", want: "/assets.xro/1.0/Assets"},
		{family: "projects", endpoint: "/Projects", want: "/projects.xro/2.0/Projects"},
	}

	for _, tc := range testCases {
		t.Run(tc.family, func(t *testing.T) {
			got, err := xeroAPIPath(tc.family, tc.endpoint)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestXeroAPIPathRejectsInvalidInputs(t *testing.T) {
	testCases := []struct {
		name     string
		family   string
		endpoint string
		want     string
	}{
		{name: "family", family: "payroll", endpoint: "/Employees", want: "unsupported Xero API family"},
		{name: "empty endpoint", family: "accounting", endpoint: "", want: "endpoint is required"},
		{name: "relative endpoint", family: "accounting", endpoint: "Invoices", want: "endpoint must start with '/'"},
		{name: "url endpoint", family: "accounting", endpoint: "/https://example.com", want: "endpoint must be a path, not a URL"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := xeroAPIPath(tc.family, tc.endpoint)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error to contain %q, got %v", tc.want, err)
			}
		})
	}
}
