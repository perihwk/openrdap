package openrdap

// IPNetwork represents information of an IP Network.
//
// IPNetwork is a topmost RDAP response object.
// https://datatracker.ietf.org/doc/html/rfc7483#section-5.4
type IPNetwork struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string   `json:"objectClassName"`
	Notices         []Notice `json:"notices"`

	Handle       string   `json:"handle"`
	StartAddress string   `json:"startAddress"`
	EndAddress   string   `json:"endAddress"`
	IPVersion    string   `json:"ipVersion"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Country      string   `json:"country"`
	ParentHandle string   `json:"parentHandle"`
	Status       []string `json:"status"`
	Entities     []Entity `json:"entities"`
	Remarks      []Remark `json:"remarks"`
	Links        []Link   `json:"links"`
	Port43       string   `json:"port43"`
	Events       []Event  `json:"events"`
}
