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
		if vc := registrar.FirstVCard(); vc != nil {
			fmt.Printf("Registrar: %s\n", vc.FullName)
		}
		fmt.Printf("RegistrarIanaID: %s\n", registrar.Handle)
	}

	abuse := domain.GetEntityFromRole("abuse")
	if abuse != nil {
		if vc := abuse.FirstVCard(); vc != nil {
			fmt.Printf("RegistrarAbuseContactEmail: %s\n", vc.Email)
			fmt.Printf("RegistrarAbuseContactPhone: %s\n", vc.Telephone)
		}
	}

	registrarURL := domain.GetRegistrarURL()
	if registrar != nil {
		fmt.Printf("RegistrarURL: %s\n", registrarURL)
	}

	registrantEntity := domain.GetEntityFromRole("registrant")
	if registrantEntity != nil {
		if vc := registrantEntity.FirstVCard(); vc != nil {
			fmt.Printf("RegistrantOrganization: %s\n", vc.Org)
			fmt.Printf("RegistrantState: %+v\n", vc.Address)
			fmt.Printf("RegistrantCountry: %+v\n", vc.Address)
			fmt.Printf("RegistrantEmail: %s\n", vc.Email)
		}
	}

	adminEntity := domain.GetEntityFromRole("administrative")
	if adminEntity != nil {
		if vc := adminEntity.FirstVCard(); vc != nil {
			fmt.Printf("AdminOrganization: %v\n", vc.Org)
			fmt.Printf("AdminState: %v\n", vc.Address.Region)
			fmt.Printf("AdminCountry: %v\n", vc.Address.Country)
			fmt.Printf("AdminEmail: %v\n", vc.Email)
		}
	}

	techEntity := domain.GetEntityFromRole("technical")
	if techEntity != nil {
		if vc := techEntity.FirstVCard(); vc != nil {
			fmt.Printf("TechOrganization: %v\n", vc.Org)
			fmt.Printf("TechState: %v\n", vc.Address.Region)
			fmt.Printf("TechCountry: %v\n", vc.Address.Country)
			fmt.Printf("TechEmail: %v\n", vc.Email)
		}
	}
}

func PrintAutnumRDAP(asn *Autnum) {
	fmt.Printf("%+v\n", asn)
}
