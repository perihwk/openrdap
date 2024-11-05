package openrdap

import (
	"fmt"
	"strings"
)

// A jCard consists of an array of properties (e.g. "fn", "tel") describing the
// individual or entity. Properties may be repeated, e.g. to represent multiple
// telephone numbers. RFC6350 documents a set of standard properties.
//
// RFC7095 describes the JSON document format, which looks like:
//
//	["vcard", [
//	  [
//	    ["version", {}, "text", "4.0"],
//	    ["fn", {}, "text", "Joe Appleseed"],
//	    ["tel", {
//	          "type":["work", "voice"],
//	        },
//	        "uri",
//	        "tel:+1-555-555-1234;ext=555"
//	    ],
//	    ...
//	  ]
//	]

type VCard struct {
	Version   string
	FullName  string
	Address   Address
	Kind      string
	Email     string
	Telephone string
	Org       string
}

type Address struct {
	Label           string
	PostOfficeBox   string
	ExtendedAddress string // apartment or suite number
	StreetAddress   string
	Locality        string // city
	Region          string // state or province
	PostalCode      string
	Country         string
}

// ParseJCard parses a jCard from its JSON representation
func parseJCard(jcardData []interface{}) (VCard, error) {
	var jcard VCard

	if jcardData == nil {
		return jcard, nil
	}

	if jcardData[0] != "vcard" {
		return jcard, fmt.Errorf("not a vcard")
	}

	// The second element should be an array of jCard properties
	properties, ok := jcardData[1].([]interface{})
	if !ok {
		return jcard, fmt.Errorf("invalid jCard properties format")
	}

	// Iterate over the properties array
	for _, prop := range properties {
		propArray, ok := prop.([]interface{})
		if !ok || len(propArray) < 4 {
			continue
		}

		// Parse the jCard field based on the property type (first element)
		propertyName := propArray[0].(string)
		propertyValue := propArray[3]

		switch propertyName {
		case "version":
			if version, ok := propertyValue.(string); ok {
				jcard.Version = version
			}
		case "fn":
			if fn, ok := propertyValue.(string); ok {
				jcard.FullName = fn
			}
		// adr property is parsed according to the following specification
		// https://datatracker.ietf.org/doc/html/rfc6350#section-6.3.1
		case "adr":
			if propMap, ok := propArray[1].(map[string]interface{}); ok {
				if label, ok := propMap["label"].(string); ok {
					parsedAddr, err := parseAddressFromLabel(label)
					if err == nil {
						jcard.Address = *parsedAddr
					}
					jcard.Address.Label = label
				}
			}

			if addr := parseStructuredAddress(propArray); addr != nil {
				jcard.Address = *addr
			}
		case "kind":
			if kind, ok := propertyValue.(string); ok {
				jcard.Kind = kind
			}
		case "email":
			if email, ok := propertyValue.(string); ok {
				jcard.Email = email
			}
		case "tel":
			if telephone, ok := propertyValue.(string); ok {
				jcard.Telephone = telephone
			}
		case "org":
			if org, ok := propertyValue.(string); ok {
				jcard.Org = org
			}
		}
	}

	return jcard, nil
}

// attempts to parse the address from the label string using its position
func parseAddressFromLabel(label string) (*Address, error) {
	parts := strings.Split(label, "\n")

	addr := &Address{}

	// Iterate from the end of the slice to the beginning
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i]) // Clean up any leading/trailing whitespace

		switch len(parts) - 1 - i {
		case 0: // CountryName
			addr.Country = part
		case 1: // PostalCode
			addr.PostalCode = part
		case 2: // Region (state or province)
			addr.Region = part
		case 3: // Locality (city)
			addr.Locality = part
		case 4: // StreetAddress
			upper := strings.ToUpper(part)
			if strings.Contains(upper, "PO") || strings.Contains(upper, "P.O") {
				addr.PostOfficeBox = part
			} else if strings.HasPrefix(upper, "SUITE") || strings.HasPrefix(upper, "APT") {
				addr.ExtendedAddress = part
			} else {
				addr.StreetAddress = part
			}
		case 5: // ExtendedAddress (apartment or suite)
			upper := strings.ToUpper(part)
			if strings.Contains(upper, "PO") || strings.Contains(upper, "P.O") {
				addr.PostOfficeBox = part
			} else if addr.ExtendedAddress == "" {
				addr.ExtendedAddress = part
			} else {
				addr.StreetAddress = part
			}
		case 6: // PostOfficeBox
			addr.PostOfficeBox = part
		}
	}

	return addr, nil
}

// parseStructuredAddress parses a structured address if available.
func parseStructuredAddress(propArray []interface{}) *Address {
	if structuredAddress, ok := propArray[3].([]interface{}); ok {
		addr := &Address{}
		emptyAddr := true // track if all the fields are ""
		for i, v := range structuredAddress {
			if str, ok := v.(string); ok {
				switch i {
				case 0:
					addr.PostOfficeBox = str
				case 1:
					addr.ExtendedAddress = str
				case 2:
					addr.StreetAddress = str
				case 3:
					addr.Locality = str
				case 4:
					addr.Region = str
				case 5:
					addr.PostalCode = str
				case 6:
					addr.Country = str
				}
				if str != "" {
					emptyAddr = false
				}
			}
		}

		if emptyAddr {
			return nil
		}
		return addr
	}
	return nil
}
