// OpenRDAP
// Copyright 2024 Paul Chihak
// MIT License, see the LICENSE file.

package bootstrap

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

// Client implements an RDAP bootstrap client
type Client struct {
	httpClient              *http.Client
	serviceRegistryIndexURL string
	registries              map[RegistryType]*Registry
}

func NewBootstrapClient(httpClient *http.Client, serviceRegistryIndexURL string) *Client {
	return &Client{
		httpClient:              httpClient,
		serviceRegistryIndexURL: serviceRegistryIndexURL,
		registries:              make(map[RegistryType]*Registry),
	}
}

func (c *Client) FetchRegistryByType(regType RegistryType, forceUpdate bool) (*Registry, error) {
	if c.registries[regType] != nil && !forceUpdate {
		return c.registries[regType], nil
	}
	var registry Registry

	resp, err := c.httpClient.Get(regType.ServiceRegistryIndexURL(c.serviceRegistryIndexURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned non-200 status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &registry); err != nil {
		return nil, err
	}

	c.registries[regType] = &registry

	return c.registries[regType], nil
}

func (c *Client) GetDomainRDAPServers(domain string) ([]*url.URL, error) {
	var err error
	if c.registries[DNS] == nil {
		c.registries[DNS], err = c.FetchRegistryByType(DNS, true)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch DNS service registry: %w", err)
		}
	}
	return c.registries[DNS].getDNSServers(domain)
}

func (c *Client) GetAutnumRDAPServers(asn string) ([]*url.URL, error) {
	var err error
	if c.registries[ASN] == nil {
		c.registries[ASN], err = c.FetchRegistryByType(ASN, true)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch ASN service registry: %w", err)
		}
	}
	return c.registries[ASN].getASNServers(asn)
}

func (c *Client) GetIPAddressRDAPServers(ip string) ([]*url.URL, error) {
	var err error
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return nil, fmt.Errorf("input %s is not an IP Address", ip)
	}
	// IPv4 address
	if ipAddress.To4() != nil {
		if c.registries[IPv4] == nil {
			c.registries[IPv4], err = c.FetchRegistryByType(IPv4, true)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch IPv4 service registry: %w", err)
			}
		}
		return c.registries[IPv4].getNetServers(ipAddress)
	} else { // IPv6 address
		if c.registries[IPv6] == nil {
			c.registries[IPv6], err = c.FetchRegistryByType(IPv6, true)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch IPv6 service registry: %w", err)
			}
		}
		return c.registries[IPv6].getNetServers(ipAddress)
	}

}
