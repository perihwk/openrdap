// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

// IPNetwork represents information of an IP Network.
//
// IPNetwork is a topmost RDAP response object.
type IPNetwork struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string
	Notices         []Notice

	Handle       string
	StartAddress string
	EndAddress   string
	IPVersion    string `json:"ipVersion"`
	Name         string
	Type         string
	Country      string
	ParentHandle string
	Status       []string
	Entities     []Entity
	Remarks      []Remark
	Links        []Link
	Port43       string
	Events       []Event
}
