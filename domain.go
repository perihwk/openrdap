package openrdap

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

	Variants    []Variant    `json:"variants"`
	Nameservers []Nameserver `json:"nameservers"`

	SecureDNS *SecureDNS `json:"secureDNS"`

	Entities []Entity `json:"entities"`

	// Status indicates the state of a registered object.
	// Embedded (no separate type):
	// Status []string
	//
	// https://tools.ietf.org/html/rfc7483#section-4.6
	Status []string `json:"status"`

	PublicIDs []PublicID `json:"publicIds"`
	Remarks   []Remark   `json:"remarks"`

	Links   []Link     `json:"links"`
	Port43  string     `json:"port43"`
	Events  []Event    `json:"events"`
	Network *IPNetwork `json:"network"`
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
	ZoneSigned       bool      `json:"zoneSigned"`
	DelegationSigned bool      `json:"delegationSigned"`
	MaxSigLife       uint64    `json:"maxSigLife"`
	DS               []DSData  `json:"dsData"`
	Keys             []KeyData `json:"keyData"`
}

// DSData is a subfield of Domain.
type DSData struct {
	Common
	KeyTag     uint64 `json:"keyTag"`
	Algorithm  uint8  `json:"algorithm"`
	Digest     string `json:"digest"`
	DigestType uint8  `json:"digestType"`

	Events []Event `json:"events"`
	Links  []Link  `json:"links"`
}

type KeyData struct {
	Flags     uint16 `json:"flags"`
	Protocol  uint8  `json:"protocol"`
	Algorithm uint8  `json:"algorithm"`
	PublicKey string `json:"publicKey"`

	Events []Event `json:"events"`
	Links  []Link  `json:"links"`
}

func (d *Domain) GetEventByName(name string) *Event {
	for _, event := range d.Events {
		if event.Action == name {
			return &event
		}
	}
	return nil
}

func (d *Domain) GetEntityFromRole(role string) *Entity {
	for _, entity := range d.Entities {
		for _, entityRole := range entity.Roles {
			if entityRole == role {
				return &entity
			}
		}
		for _, nestedEntity := range entity.Entities {
			for _, nestedEntityRole := range nestedEntity.Roles {
				if nestedEntityRole == role {
					return &nestedEntity
				}
			}
		}
	}
	return nil
}

func (d *Domain) GetRegistrarURL() string {
	for _, link := range d.Links {
		if link.Rel == "self" {
			return link.Value
		}
	}
	return ""
}

func (d *Domain) GetNameServersDNS() []string {
	var nameservers []string
	for _, ns := range d.Nameservers {
		nameservers = append(nameservers, ns.LDHName)
	}
	return nameservers
}
