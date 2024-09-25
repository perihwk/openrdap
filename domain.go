// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

import (
	"encoding/json"
	"fmt"
)

// Domain represents information about a DNS name and point of delegation.
//
// Domain is a topmost RDAP response object.
// https://tools.ietf.org/html/rfc7483
type Domain struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string   `json:"objectClassName"`

	Notices []Notice `json:"notices"`

	Handle      string `json:"handle"`
	LDHName     string `json:"ldhName"`
	UnicodeName string `json:"unicodename"`

	Variants    []Variant             `json:"variants"`
	Nameservers map[string]Nameserver `json:"nameservers"`

	SecureDNS *SecureDNS

	Entities []Entity `json:"entities"`

	// Status indicates the state of a registered object.
	// Embedded (no separate type):
	// Status []string
	//
	// https://tools.ietf.org/html/rfc7483#section-4.6
	Status []string `json:"status"`

	PublicIDs []PublicID `json:"publicIds"`
	Remarks   []Remark   `json:"remarks"`

	Links   []Link           `json:"links"`
	Port43  string           `json:"port43"`
	Events  map[string]Event `json:"events"`
	Network *IPNetwork       `json:"network"`
}

// Variant is a subfield of Domain.
type Variant struct {
	Common
	Relation     []string
	IDNTable     string `json:"idnTable"`
	VariantNames []VariantName
}

// VariantName is a subfield of Variant.
type VariantName struct {
	Common
	LDHName     string `json:"ldhName"`
	UnicodeName string
}

// SecureDNS is ia subfield of Domain.
type SecureDNS struct {
	Common
	ZoneSigned       *bool
	DelegationSigned *bool
	MaxSigLife       *uint64
	DS               []DSData  `json:"dsData"`
	Keys             []KeyData `json:"keyData"`
}

// DSData is a subfield of Domain.
type DSData struct {
	Common
	KeyTag     *uint64
	Algorithm  *uint8
	Digest     string
	DigestType *uint8

	Events []Event
	Links  []Link
}

type KeyData struct {
	Flags     *uint16
	Protocol  *uint8
	Algorithm *uint8
	PublicKey string

	Events []Event
	Links  []Link
}

func (d *Domain) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	type Alias Domain

	temp := &struct {
		*Alias

		Events      []Event      `json:"events"`
		Nameservers []Nameserver `json:"nameservers"`
	}{
		Alias: (*Alias)(d),
	}

	d.Events = make(map[string]Event)
	d.Nameservers = make(map[string]Nameserver)

	if err := json.Unmarshal(data, temp); err != nil {
		return fmt.Errorf("failed to parse events: %w", err)
	}

	for _, event := range temp.Events {
		d.Events[event.Action] = event
	}

	for _, ns := range temp.Nameservers {
		d.Nameservers[ns.LDHName] = ns
	}

	return nil
}

func (d *Domain) MarshalJSON() ([]byte, error) {
	// Convert map[string]Event back to []Event for JSON serialization
	events := make([]Event, 0, len(d.Events))
	for _, event := range d.Events {
		events = append(events, event)
	}

	// Create a temporary struct to hold the events array
	temp := struct {
		Events []Event `json:"events"`
	}{
		Events: events,
	}

	// Marshal the temp struct with the []Event
	return json.Marshal(temp)
}

func (d *Domain) GetEntityFromRole(role string) (*Entity, error) {
	for _, entity := range d.Entities {
		for _, entityRole := range entity.Roles {
			if entityRole == role {
				return &entity, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find an entity with role %s", role)
}

func (d *Domain) GetRegistrarURL() (string, error) {
	for _, link := range d.Links {
		if link.Rel == "self" {
			return link.Value, nil
		}
	}
	return "", fmt.Errorf("unable to find registrar URL")
}
