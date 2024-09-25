// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

// Autnum represents information of Autonomous System registrations.
//
// Autnum is a topmost RDAP response object.
type Autnum struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string
	Notices         []Notice

	Handle      string
	StartAutnum *uint32
	EndAutnum   *uint32
	IPVersion   string `json:"ipVersion"`
	Name        string
	Type        string
	Status      []string
	Country     string
	Entities    []Entity
	Remarks     []Remark
	Links       []Link
	Port43      string
	Events      []Event
}
