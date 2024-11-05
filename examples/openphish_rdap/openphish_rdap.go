package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/perihwk/openrdap"
	"github.com/perihwk/openrdap/bootstrap"
)

func main() {

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
	bClient := bootstrap.NewBootstrapClient(httpClient, "")

	rdapClient := openrdap.NewClient(httpClient, bClient)

	// URL of the file to download
	fileURL := "https://raw.githubusercontent.com/openphish/public_feed/refs/heads/main/feed.txt"

	// Step 1: Download the file
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return
	}

	// Step 2: Use a scanner to read the response body line-by-line
	scanner := bufio.NewScanner(resp.Body)

	// Step 3: Process each URL in the response body
	for scanner.Scan() {
		line := scanner.Text()
		parsedURL, err := url.Parse(line)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			fmt.Printf("Invalid URL: %s\n", line)
			continue
		}
		fmt.Printf("Processing URL: %s\n", parsedURL)
		domain := GetTLDPlusOne(parsedURL.Host)
		domainInfo, err := rdapClient.GetRDAPFromDomain(ctx, domain)
		if err != nil {
			fmt.Println("BIG ASS ERROR: ", err)
			continue
		}
		fmt.Printf("\tRegistryDomainID: %s\n", domainInfo.Handle)
		fmt.Printf("\tDomainName: %s\n", domainInfo.LDHName)
		if event := domainInfo.GetEventByName("registration"); event != nil {
			fmt.Printf("\tCreatedDate: %s\n", event.Date)
		}
		if event := domainInfo.GetEventByName("last changed"); event != nil {
			fmt.Printf("\tUpdatedDate: %s\n", event.Date)
		}
		if event := domainInfo.GetEventByName("expiration"); event != nil {
			fmt.Printf("\tRegistrarExpirationDate: %s\n", event.Date)
		}
		fmt.Printf("\tRegistrarWhoisServer: %s\n", domainInfo.Port43)
		fmt.Printf("\tNameServer: %s\n", domainInfo.GetNameServersDNS())
		fmt.Printf("\tDomainStatus: %s\n", domainInfo.Status)

		registrar := domainInfo.GetEntityFromRole("registrar")
		if registrar != nil {
			fmt.Printf("\tRegistrar: %s\n", registrar.VCards[0].FullName)
			fmt.Printf("\tRegistrarIanaID: %s\n", registrar.Handle)
		}

		abuse := domainInfo.GetEntityFromRole("abuse")
		if abuse != nil {
			fmt.Printf("\tRegistrarAbuseContactEmail: %s\n", abuse.VCards[0].Email)
			fmt.Printf("\tRegistrarAbuseContactPhone: %s\n", abuse.VCards[0].Telephone)
		}

		registrarURL := domainInfo.GetRegistrarURL()
		if registrar != nil {
			fmt.Printf("\tRegistrarURL: %s\n", registrarURL)
		}

		registrantEntity := domainInfo.GetEntityFromRole("registrant")
		if registrantEntity != nil {
			fmt.Printf("\tRegistrantOrganization: %s\n", registrantEntity.VCards[0].Org)
			fmt.Printf("\tRegistrantState: %+v\n", registrantEntity.VCards[0].Address)
			fmt.Printf("\tRegistrantCountry: %+v\n", registrantEntity.VCards[0].Address)
			fmt.Printf("\tRegistrantEmail: %s\n", registrantEntity.VCards[0].Email)
		}

		adminEntity := domainInfo.GetEntityFromRole("administrative")
		if adminEntity != nil {
			fmt.Printf("\tAdminOrganization: %v\n", adminEntity.VCards[0].Org)
			fmt.Printf("\tAdminState: %v\n", adminEntity.VCards[0].Address.Region)
			fmt.Printf("\tAdminCountry: %v\n", adminEntity.VCards[0].Address.Country)
			fmt.Printf("\tAdminEmail: %v\n", adminEntity.VCards[0].Email)
		}

		techEntity := domainInfo.GetEntityFromRole("technical")
		if techEntity != nil {
			fmt.Printf("\tTechOrganization: %v\n", techEntity.VCards[0].Org)
			fmt.Printf("\tTechState: %v\n", techEntity.VCards[0].Address.Region)
			fmt.Printf("\tTechCountry: %v\n", techEntity.VCards[0].Address.Country)
			fmt.Printf("\tTechEmail: %v\n", techEntity.VCards[0].Email)
		}

	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response body:", err)
	}
}

func GetTLDPlusOne(domain string) string {
	parts := strings.Split(domain, ".")

	// If it's a single part (e.g. "localhost") return itself
	if len(parts) < 2 {
		return domain
	}

	// Return the last two parts (TLD + 1)
	return strings.Join(parts[len(parts)-2:], ".")
}
