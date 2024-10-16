package openrdap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetRDAPInfoFromServer(t *testing.T) {
	fileData, err := os.ReadFile("test/example_domain_perihwk.json")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fileData)
	}))
	defer mockServer.Close()

	client := &Client{
		httpClient: mockServer.Client(),
	}

	ctx := context.Background()
	query := "perihwk.com"
	searchType := DNS

	rdapInfo, err := client.GetRDAPInfoFromServer(ctx, mockServer.URL+"/", query, searchType)
	if err != nil {
		t.Fatalf("Failed to get RDAP info: %v", err)
	}

	domainInfo, ok := rdapInfo.(*Domain)
	if !ok {
		t.Fatalf("Expected *Domain, got %T", rdapInfo)
	}

	expectedDomainName := "PERIHWK.COM"
	if domainInfo.LDHName != expectedDomainName {
		t.Errorf("Expected domain name %s, got %s", expectedDomainName, domainInfo.LDHName)
	}
}

func TestGetRDAPFromDomain(t *testing.T) {
	fileData, err := os.ReadFile("test/example_domain_perihwk.json")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fileData)
	}))
	defer mockServer.Close()

	client := &Client{
		httpClient: mockServer.Client(),
	}

	ctx := context.Background()
	query := "perihwk.com"
	searchType := DNS

	rdapInfo, err := client.GetRDAPInfoFromServer(ctx, mockServer.URL+"/", query, searchType)
	if err != nil {
		t.Fatalf("Failed to get RDAP info: %v", err)
	}

	domainInfo, ok := rdapInfo.(*Domain)
	if !ok {
		t.Fatalf("Expected *Domain, got %T", rdapInfo)
	}

	expectedDomainName := "PERIHWK.COM"
	if domainInfo.LDHName != expectedDomainName {
		t.Errorf("Expected domain name %s, got %s", expectedDomainName, domainInfo.LDHName)
	}
}

func TestGetRDAPFromIP(t *testing.T) {
	fileData, err := os.ReadFile("test/example_ip_8888.json")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fileData)
	}))
	defer mockServer.Close()

	client := &Client{
		httpClient: mockServer.Client(),
	}

	ctx := context.Background()
	query := "8.8.8.8"
	searchType := IPv4

	rdapInfo, err := client.GetRDAPInfoFromServer(ctx, mockServer.URL+"/", query, searchType)
	if err != nil {
		t.Fatalf("Failed to get RDAP info: %v", err)
	}

	ipInfo, ok := rdapInfo.(*IPNetwork)
	if !ok {
		t.Fatalf("Expected *IPnetwork, got %T", rdapInfo)
	}

	expectedName := "GOGL"
	if ipInfo.Name != expectedName {
		t.Errorf("Expected domain name %s, got %s", expectedName, ipInfo.Name)
	}
}

func TestGetRDAPFromAutnum(t *testing.T) {
	fileData, err := os.ReadFile("test/example_asn_23552.json")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fileData)
	}))
	defer mockServer.Close()

	client := &Client{
		httpClient: mockServer.Client(),
	}

	ctx := context.Background()
	query := "23552"
	searchType := ASN

	rdapInfo, err := client.GetRDAPInfoFromServer(ctx, mockServer.URL+"/", query, searchType)
	if err != nil {
		t.Fatalf("Failed to get RDAP info: %v", err)
	}

	asnInfo, ok := rdapInfo.(*Autnum)
	if !ok {
		t.Fatalf("Expected *Autnum, got %T", rdapInfo)
	}

	expectedName := "KORNU-AS-KR-KR"
	if asnInfo.Name != expectedName {
		t.Errorf("Expected domain name %s, got %s", expectedName, asnInfo.Name)
	}
}
