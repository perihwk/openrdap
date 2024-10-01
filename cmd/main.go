package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
	fmt.Printf("NameServer: %s\n", domain.GetNameServersDNS())

	registrar, err := domain.GetEntityFromRole("registrar")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Registrar: %s\n", registrar.VCards[0].FullName)
		fmt.Printf("RegistrarIanaID: %s\n", registrar.Handle)
	}

	abuse, err := domain.GetEntityFromRole("abuse")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("RegistrarAbuseContactEmail: %s\n", abuse.VCards[0].Email)
		fmt.Printf("RegistrarAbuseContactPhone: %s\n", abuse.VCards[0].Telephone)
	}

	registrarURL, err := domain.GetRegistrarURL()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("RegistrarURL: %s\n", registrarURL)

	}
	fmt.Printf("DomainStatus: %s\n", domain.Status)

	registrantEntity, err := domain.GetEntityFromRole("registrant")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("RegistrantOrganization: %s\n", registrantEntity.VCards[0].Org)
		fmt.Printf("RegistrantState: %+v\n", registrantEntity.VCards[0].Address)
		fmt.Printf("RegistrantCountry: %+v\n", registrantEntity.VCards[0].Address)
		fmt.Printf("RegistrantEmail: %s\n", registrantEntity.VCards[0].Email)
	}

	adminEntity, err := domain.GetEntityFromRole("administrative")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("AdminOrganization: %v\n", adminEntity.VCards[0].Org)
		fmt.Printf("AdminState: %v\n", adminEntity.VCards[0].Address.Region)
		fmt.Printf("AdminCountry: %v\n", adminEntity.VCards[0].Address.Country)
		fmt.Printf("AdminEmail: %v\n", adminEntity.VCards[0].Email)
	}

	techEntity, err := domain.GetEntityFromRole("technical")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("AdminOrganization: %v\n", techEntity.VCards[0].Org)
		fmt.Printf("AdminState: %v\n", techEntity.VCards[0].Address.Region)
		fmt.Printf("AdminCountry: %v\n", techEntity.VCards[0].Address.Country)
		fmt.Printf("AdminEmail: %v\n", techEntity.VCards[0].Email)
	}
}
