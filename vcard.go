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
			jcard.Version = propertyValue.(string)
		case "fn":
			jcard.FullName = propertyValue.(string)
		// adr property is parsed according to the following specification
		// https://datatracker.ietf.org/doc/html/rfc6350#section-6.3.1
		case "adr":
			label, ok := propArray[1].(map[string]interface{})["label"].(string)
			if ok {
				parsedAddr, err := parseAddressFromLabel(label)
				if err == nil {
					jcard.Address = *parsedAddr
				}
				jcard.Address.Label = label
			} else {
				if structuredAddress, ok := propArray[3].([]interface{}); ok {
					for i, v := range structuredAddress {
						str, ok := v.(string)
						if !ok {
							continue
						}

						switch i {
						case 0:
							jcard.Address.PostOfficeBox = str
						case 1:
							jcard.Address.ExtendedAddress = str
						case 2:
							jcard.Address.StreetAddress = str
						case 3:
							jcard.Address.Locality = str
						case 4:
							jcard.Address.Region = str
						case 5:
							jcard.Address.PostalCode = str
						case 6:
							jcard.Address.Country = str

						}
					}
				}
			}
		case "kind":
			jcard.Kind = propertyValue.(string)
		case "email":
			jcard.Email = propertyValue.(string)
		case "tel":
			jcard.Telephone = propertyValue.(string)
		case "org":
			jcard.Org = propertyValue.(string)
		}
	}

	return jcard, nil
}

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
			} else {
				addr.StreetAddress = part
			}
		case 5: // ExtendedAddress (apartment or suite)
			upper := strings.ToUpper(part)
			if strings.Contains(upper, "PO") || strings.Contains(upper, "P.O") {
				addr.PostOfficeBox = part
			} else {
				addr.ExtendedAddress = part
			}
		case 6: // PostOfficeBox
			addr.PostOfficeBox = part
		}
	}

	return addr, nil
}
