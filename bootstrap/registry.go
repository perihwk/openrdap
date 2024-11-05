package bootstrap

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const ServiceRegistryDefaultBaseURL = "https://data.iana.org/rdap/"

// A RegistryType represents a bootstrap registry type.
type RegistryType int

const (
	DNS RegistryType = iota
	IPv4
	IPv6
	ASN
)

func (r RegistryType) String() string {
	switch r {
	case DNS:
		return "dns"
	case IPv4:
		return "ipv4"
	case IPv6:
		return "ipv6"
	case ASN:
		return "asn"
	default:
		panic("Unknown RegistryType")
	}
}

// Validate the RegistryType value
func (r *RegistryType) Set(value string) error {
	switch strings.ToLower(value) {
	case "dns":
		*r = DNS
	case "ipv4":
		*r = IPv4
	case "ipv6":
		*r = IPv6
	case "asn":
		*r = ASN
	default:
		return fmt.Errorf("invalid registry-type %s, must be one of: dns, ipv4, ipv6, asn", value)
	}
	return nil
}

func (r RegistryType) ServiceRegistryIndexURL(baseURL string) string {
	if baseURL == "" {
		baseURL = ServiceRegistryDefaultBaseURL
	}

	switch r {
	case DNS:
		return baseURL + "dns.json"
	case IPv4:
		return baseURL + "ipv4.json"
	case IPv6:
		return baseURL + "ipv6.json"
	case ASN:
		return baseURL + "asn.json"
	default:
		panic("Unknown RegistryType")
	}
}

type Registry struct {
	Version     string                `json:"version"`
	Publication string                `json:"publication"`
	Description string                `json:"description"`
	Services    map[string][]*url.URL `json:"services"`
}

func (r *Registry) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	type Alias Registry

	temp := &struct {
		*Alias

		Services [][][]string `json:"services"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	r.Services = make(map[string][]*url.URL)

	for _, service := range temp.Services {
		for _, key := range service[0] {
			parsedURL, err := parseURLs(service[1])
			if err != nil {
				return fmt.Errorf("failed to parse URLs: %w", err)
			}
			r.Services[key] = parsedURL
		}
	}

	return nil
}

func (r *Registry) MarshalJSON() ([]byte, error) {
	if len(r.Services) == 0 {
		return json.Marshal("")
	}

	// Create a temporary structure to hold the formatted services
	temp := struct {
		Services [][][]string `json:"services"`
	}{
		Services: make([][][]string, 0, len(r.Services)),
	}

	// Iterate over the Services map
	for key, urls := range r.Services {
		// Convert []*url.URL back to a list of URL strings
		urlStrs := urlsToStrings(urls)

		// Append the key and URLs to the temp.Services slice in the expected format
		temp.Services = append(temp.Services, [][]string{
			{key},   // First slice contains the key
			urlStrs, // Second slice contains URLs
		})
	}

	// Marshal the entire struct into JSON
	return json.Marshal(temp)
}

func (r *Registry) getNetServers(ipAddr net.IP) ([]*url.URL, error) {
	for cidr := range r.Services {
		// Parse the CIDR block
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, ErrInvalidCIDR
		}

		// Check if the IP is in the CIDR range
		if ipNet.Contains(ipAddr) {
			return r.Services[cidr], nil
		}
	}
	return nil, ErrRDAPNotSupported
}

func (r *Registry) getDNSServers(domain string) ([]*url.URL, error) {
	parts := strings.Split(domain, ".")
	tld := parts[len(parts)-1]

	if r.Services[tld] == nil {
		return nil, fmt.Errorf("tld %s not supported: %w", tld, ErrRDAPNotSupported)
	}
	return r.Services[tld], nil
}

func (r *Registry) getASNServers(input string) ([]*url.URL, error) {
	asn, err := parseASN(input)
	if err != nil {
		return nil, err
	}

	for rangeKey, urls := range r.Services {
		rangeParts := strings.Split(rangeKey, "-")
		if len(rangeParts) > 2 {
			return nil, ErrInvalidASNRange
		}

		minASN, err := strconv.ParseUint(rangeParts[0], 10, 32)
		if err != nil {
			return nil, err
		}

		// handle case where there is only one ASN in the range or where the input is the min ASN
		if len(rangeParts) == 1 {
			if minASN == asn {
				return urls, nil
			}

			continue
		}

		maxASN, err := strconv.ParseUint(rangeParts[1], 10, 32)
		if err != nil {
			return nil, err
		}

		if asn >= minASN && asn <= maxASN {
			return urls, nil
		}
	}
	return nil, ErrRDAPNotSupported
}
