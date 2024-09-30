// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

import "fmt"

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
	CountryName     string
}

// ParseJCard parses a jCard from its JSON representation
func parseJCard(jcardData []interface{}) (VCard, error) {
	var jcard VCard

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
		propertyType := propArray[2].(string)

		switch propertyName {
		case "version":
			jcard.Version = propertyValue.(string)
		case "fn":
			jcard.FullName = propertyValue.(string)
		case "adr":
			if propertyType == "text" {
				label, ok := propArray[1].(map[string]interface{})["label"].(string)
				if ok {
					jcard.Address.Label = label
				}
			} else if propertyType == "array" {
				structuredAddress, ok := propArray[3].([]string)
				if !ok {
					continue
				}
				jcard.Address.PostOfficeBox = structuredAddress[0]
				jcard.Address.ExtendedAddress = structuredAddress[1]
				jcard.Address.StreetAddress = structuredAddress[2]
				jcard.Address.Locality = structuredAddress[3]
				jcard.Address.Region = structuredAddress[4]
				jcard.Address.PostalCode = structuredAddress[5]
				jcard.Address.CountryName = structuredAddress[6]
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
