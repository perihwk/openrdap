// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package openrdap

// RDAP Conformance
// Appears in topmost JSON objects only, embedded (no separate type):
// Conformance []string `rdap:"rdapConformance"`
//
// https://tools.ietf.org/html/rfc7483#section-4.1

// Link signifies a link another resource on the Internet.
//
// https://tools.ietf.org/html/rfc7483#section-4.2
type Link struct {
	Value    string   `json:"value"`
	Rel      string   `json:"rel"`
	Href     string   `json:"href"`
	HrefLang []string `json:"hreflang"`
	Title    string   `json:"title"`
	Media    string   `json:"media"`
	Type     string   `json:"type"`
}

// Notice contains information about the entire RDAP response.
//
// https://tools.ietf.org/html/rfc7483#section-4.3
type Notice struct {
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Description []string `json:"description"`
	Links       []Link   `json:"links"`
}

// Remark contains information about the containing RDAP object.
//
// https://tools.ietf.org/html/rfc7483#section-4.3
type Remark struct {
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Description []string `json:"description"`
	Links       []Link   `json:"links"`
}

// Language Identifier
// Appears in anywhere, embedded (no separate type):
// Lang string
//
// https://tools.ietf.org/html/rfc7483#section-4.4

// Event represents some event which has occured/may occur in the future..
//
// https://tools.ietf.org/html/rfc7483#section-4.5
type Event struct {
	Action string `json:"eventAction"`
	Actor  string `json:"eventActor"`
	Date   string `json:"eventDate"`
	Links  []Link `json:"links"`
}

// Port43 indicates the IP/FQDN of a WHOIS server.
// Embedded (no separate type):
// Port43 string
//
// https://tools.ietf.org/html/rfc7483#section-4.7

// PublicID maps a public identifier to an object class.
//
// https://tools.ietf.org/html/rfc7483#section-4.8
type PublicID struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// ObjectClassName specifies the object type as a string.
// Embedded (no separate type):
// ObjectClassName string
//
// https://tools.ietf.org/html/rfc7483#section-4.9

// Common contains fields which may appear anywhere in an RDAP response.
type Common struct {
	Lang string `json:"lang"`
}
