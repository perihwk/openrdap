package openrdap

// Autnum represents information of Autonomous System registrations.
//
// Autnum is a topmost RDAP response object.
// https://datatracker.ietf.org/doc/html/rfc7483#section-5.5
type Autnum struct {
	Common
	Conformance     []string `json:"rdapConformance"`
	ObjectClassName string   `json:"objectClassName"`
	Notices         []Notice `json:"notices"`

	Handle      string   `json:"handle"`
	StartAutnum uint32   `json:"startAutnum"`
	EndAutnum   uint32   `json:"endAutnum"`
	IPVersion   string   `json:"ipVersion"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Status      []string `json:"status"`
	Country     string   `json:"country"`
	Entities    []Entity `json:"entities"`
	Remarks     []Remark `json:"remarks"`
	Links       []Link   `json:"links"`
	Port43      string   `json:"port43"`
	Events      []Event  `json:"events"`
}
