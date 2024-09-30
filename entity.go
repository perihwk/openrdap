// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

import (
	"encoding/json"
	"fmt"
)

// Entity represents information of an organisation or person.
//
// Entity is a topmost RDAP response object.
type Entity struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string
	Notices         []Notice

	Handle       string
	VCards       []VCard `json:"vcardArray"`
	Roles        []string
	PublicIDs    []PublicID
	Entities     []Entity
	Remarks      []Remark
	Links        []Link
	Events       []Event
	AsEventActor []Event
	Status       []string
	Port43       string
	Networks     []IPNetwork
	Autnums      []Autnum
}

// UnmarshalJSON for Entity to handle custom vCard processing
func (e *Entity) UnmarshalJSON(data []byte) error {
	// Create a temporary type to handle raw unmarshal first
	type Alias Entity
	aux := &struct {
		*Alias
		RawVCard []interface{} `json:"vcardArray"`
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal entity: %w", err)
	}

	// Process the rawVCard data into the structured VCard type
	var err error
	parsedJCard, err := parseJCard(aux.RawVCard)
	if err != nil {
		return err
	}

	e.VCards = append(e.VCards, parsedJCard)

	return nil
}

// returns the first VCard that contains the specified role
func (e *Entity) GetVCardFromRole(role string) (*VCard, error) {
	for _, entityRole := range e.Roles {
		if entityRole == role {
			return &e.VCards[0], nil
		}
	}

	var vCard *VCard
	for _, entity := range e.Entities {
		vCard, err := entity.GetVCardFromRole(role)
		if err != nil {
			return nil, err
		}
		return vCard, nil

	}
	return vCard, nil
}

func (e *Entity) ContainsRole(role string) bool {
	if e.Roles == nil {
		return false
	}
	for _, entityRole := range e.Roles {
		if entityRole == role {
			return true
		}
	}
	return false
}
