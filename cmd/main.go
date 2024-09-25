package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"openrdap"
	"openrdap/bootstrap"
	"time"
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
	bClient := bootstrap.NewBootstrapClient(&ctx, httpClient, "")

	rdapClient := openrdap.NewClient(&ctx, httpClient, bClient)

	domain, err := rdapClient.GetRDAPFromDomain("x.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("RegistryDomainID: %s\n", domain.Handle)
	fmt.Printf("DomainName: %s\n", domain.LDHName)
	fmt.Printf("CreatedDate: %s\n", domain.Events["registration"].Date)
	fmt.Printf("UpdatedDate: %s\n", domain.Events["last changed"].Date)
	fmt.Printf("RegistrarExpirationDate: %s\n", domain.Events["expiration"].Date)
	fmt.Printf("RegistrarWhoisServer: %s\n", domain.Port43)
	fmt.Printf("NameServer: %v\n", domain.Nameservers)
	fmt.Printf("DNSSec: %+v\n", domain.SecureDNS)

	registrar, err := domain.GetEntityFromRole("registrar")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Registrar: %+v\n", registrar.VCards)

	registrarURL, err := domain.GetRegistrarURL()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("RegistrarURL: %+v\n", registrarURL)
		fmt.Printf("RegistrarIanaID: %s\n", registrar.Handle)

		abuse, err := registrar.GetVCardFromRole("abuse")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("RegistrarAbuseContactEmail: %s\n", abuse.Email)
		fmt.Printf("RegistrarAbuseContactPhone: %s\n", abuse.Telephone)
	}
	fmt.Printf("DomainStatus: %s\n", domain.Status)

	registrantEntity, err := domain.GetEntityFromRole("registrant")
	if err != nil {
		fmt.Println(err)
	} else {
		registrantInfo, err := registrantEntity.GetVCardFromRole("registrant")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("RegistrantOrganization: %v\n", registrantInfo.Org)
		fmt.Printf("RegistrantState: %v\n", registrantInfo.Address)
		fmt.Printf("RegistrantCountry: %v\n", registrantInfo.Address)
		fmt.Printf("RegistrantEmail: %v\n", registrantInfo.Email)
	}

	adminEntity, err := domain.GetEntityFromRole("administrative")
	if err != nil {
		fmt.Println(err)
	} else {
		adminInfo, err := adminEntity.GetVCardFromRole("administrative")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("AdminOrganization: %v\n", adminInfo.Org)
		fmt.Printf("AdminState: %v\n", adminInfo.Address)
		fmt.Printf("AdminCountry: %v\n", adminInfo.Address)
		fmt.Printf("AdminEmail: %v\n", adminInfo.Email)
	}

	techEntity, err := domain.GetEntityFromRole("technical")
	if err != nil {
		fmt.Println(err)
	} else {
		techInfo, err := techEntity.GetVCardFromRole("technical")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("TechOrganization: %v\n", techInfo.Org)
		fmt.Printf("TechState: %v\n", techInfo.Address)
		fmt.Printf("TechCountry: %v\n", techInfo.Address)
		fmt.Printf("TechEmail: %v\n", techInfo.Email)
	}
}
