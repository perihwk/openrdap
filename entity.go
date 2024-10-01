package openrdap

import (
	"encoding/json"
	"fmt"
)

// Entity represents information of an organisation or person.
//
// Entity is a topmost RDAP response object.
// https://datatracker.ietf.org/doc/html/rfc7483#section-5.1
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
