package openrdap

// Nameserver represents information of a DNS nameserver.
//
// Nameserver is a topmost RDAP response object.
// https://datatracker.ietf.org/doc/html/rfc7483#section-5.2
type Nameserver struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string
	Notices         []Notice

	Handle      string
	LDHName     string `json:"ldhName"`
	UnicodeName string

	IPAddresses *IPAddressSet `json:"ipAddresses"`

	Entities []Entity
	Status   []string
	Remarks  []Remark
	Links    []Link
	Port43   string
	Events   []Event
}

// IPAddressSet is a subfield of Nameserver.
type IPAddressSet struct {
	Common
	V6 []string
	V4 []string
}
