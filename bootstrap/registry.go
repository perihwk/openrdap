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
	Type        RegistryType          `json:"-"`
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

func (r *Registry) getNetServers(ipAddr string) ([]*url.URL, error) {
	// Parse the IP address
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipAddr)
	}

	for cidr, _ := range r.Services {
		// Parse the CIDR block
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR block: %s", cidr)
		}

		// Check if the IP is in the CIDR range
		if ipNet.Contains(ip) {
			return r.Services[cidr], nil
		}
	}
	return nil, fmt.Errorf("RDAP server for ip address %s cannot be found", ipAddr)
}

func (r *Registry) getDNSServers(domain string) ([]*url.URL, error) {
	parts := strings.Split(domain, ".")
	tld := parts[len(parts)-1]

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
			return nil, fmt.Errorf("invalid ASN range %s: %w", rangeKey, err)
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

		if asn > minASN && asn <= maxASN {
			return urls, nil
		}
	}
	return nil, fmt.Errorf("could not find proper range for ASN %s", input)
}
