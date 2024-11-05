package openrdap

import "fmt"

func PrintDomainRDAP(domain *Domain) {
	fmt.Printf("RegistryDomainID: %s\n", domain.Handle)
	fmt.Printf("DomainName: %s\n", domain.LDHName)

	if event := domain.GetEventByName("registration"); event != nil {
		fmt.Printf("CreatedDate: %s\n", event.Date)
	}
	if event := domain.GetEventByName("last changed"); event != nil {
		fmt.Printf("UpdatedDate: %s\n", event.Date)
	}
	if event := domain.GetEventByName("expiration"); event != nil {
		fmt.Printf("RegistrarExpirationDate: %s\n", event.Date)
	}
	fmt.Printf("RegistrarWhoisServer: %s\n", domain.Port43)
	fmt.Printf("NameServer: %s\n", domain.GetNameServersDNS())
	fmt.Printf("DomainStatus: %s\n", domain.Status)

	registrar := domain.GetEntityFromRole("registrar")
	if registrar != nil {
		fmt.Printf("Registrar: %s\n", registrar.VCards[0].FullName)
		fmt.Printf("RegistrarIanaID: %s\n", registrar.Handle)
	}

	abuse := domain.GetEntityFromRole("abuse")
	if abuse != nil {
		fmt.Printf("RegistrarAbuseContactEmail: %s\n", abuse.VCards[0].Email)
		fmt.Printf("RegistrarAbuseContactPhone: %s\n", abuse.VCards[0].Telephone)
	}

	registrarURL := domain.GetRegistrarURL()
	if registrar != nil {
		fmt.Printf("RegistrarURL: %s\n", registrarURL)
	}

	registrantEntity := domain.GetEntityFromRole("registrant")
	if registrantEntity != nil {
		fmt.Printf("RegistrantOrganization: %s\n", registrantEntity.VCards[0].Org)
		fmt.Printf("RegistrantState: %+v\n", registrantEntity.VCards[0].Address)
		fmt.Printf("RegistrantCountry: %+v\n", registrantEntity.VCards[0].Address)
		fmt.Printf("RegistrantEmail: %s\n", registrantEntity.VCards[0].Email)
	}

	adminEntity := domain.GetEntityFromRole("administrative")
	if adminEntity != nil {
		fmt.Printf("AdminOrganization: %v\n", adminEntity.VCards[0].Org)
		fmt.Printf("AdminState: %v\n", adminEntity.VCards[0].Address.Region)
		fmt.Printf("AdminCountry: %v\n", adminEntity.VCards[0].Address.Country)
		fmt.Printf("AdminEmail: %v\n", adminEntity.VCards[0].Email)
	}

	techEntity := domain.GetEntityFromRole("technical")
	if techEntity != nil {
		fmt.Printf("TechOrganization: %v\n", techEntity.VCards[0].Org)
		fmt.Printf("TechState: %v\n", techEntity.VCards[0].Address.Region)
		fmt.Printf("TechCountry: %v\n", techEntity.VCards[0].Address.Country)
		fmt.Printf("TechEmail: %v\n", techEntity.VCards[0].Email)
	}
}

func PrintAutnumRDAP(asn *Autnum) {
	fmt.Printf("%+v\n", asn)
}
