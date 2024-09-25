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
	Address   string
	Kind      string
	Email     string
	Telephone string
	Org       string
}

// ParseJCard parses a jCard from its JSON representation
func parseJCard(jcardData []interface{}) (VCard, error) {
	var jcard VCard
	if len(jcardData) < 2 {
		return jcard, fmt.Errorf("invalid jCard format")
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
		propertyValue := propArray[3].(string)

		switch propertyName {
		case "version":
			jcard.Version = propertyValue
		case "fn":
			jcard.FullName = propertyValue
		case "adr":
			label, ok := propArray[1].(map[string]interface{})["label"].(string)
			if ok {
				jcard.Address = label
			}
		case "kind":
			jcard.Kind = propertyValue
		case "email":
			jcard.Email = propertyValue
		case "tel":
			jcard.Telephone = propertyValue
		case "org":
			jcard.Org = propertyValue
		}
	}

	return jcard, nil
}
