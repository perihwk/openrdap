package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/perihwk/openrdap"
	"github.com/perihwk/openrdap/bootstrap"
)

func main() {
	// Define flags for RDAP server URL and service registry URL
	serviceRegistryURL := flag.String("service-registry-url", "", "The URL of the service registry (optional)")
	query := flag.String("query", "", "Name to query")
	registryType := flag.String("registry-type", "dns", "Type of registry to query (dns, ipv4, ipv6, asn)")

	// Parse command-line flags
	flag.Parse()
	ctx := context.Background()

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	bClient := bootstrap.NewBootstrapClient(httpClient, *serviceRegistryURL)

	rdapClient := openrdap.NewClient(httpClient, bClient)

	var regType bootstrap.RegistryType
	if err := regType.Set(*registryType); err != nil {
		fmt.Println("Error:", err)
		flag.Usage()
		return
	}

	switch regType {
	case bootstrap.DNS:
		domain, err := rdapClient.GetRDAPFromDomain(ctx, *query)
		if err != nil {
			fmt.Println(err)
		}
		openrdap.PrintDomainRDAP(domain)
	case bootstrap.ASN:
		autnum, err := rdapClient.GetRDAPFromAutnum(ctx, *query)
		if err != nil {
			fmt.Println(err)
		}
		openrdap.PrintAutnumRDAP(autnum)
	}
}
