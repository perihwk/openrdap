// OpenRDAP
// Copyright 2024 Paul Chihak
// MIT License, see the LICENSE file.

package bootstrap

import (
	"context"
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

func (c *Client) FetchAllRegistries(ctx context.Context) error {
	registryTypes := []RegistryType{DNS, IPv4, IPv6, ASN}
	for _, regType := range registryTypes {
		req, err := http.NewRequestWithContext(ctx, "GET", regType.ServiceRegistryIndexURL(c.serviceRegistryIndexURL), nil)
		if err != nil {
			return fmt.Errorf("unable to create request for registry %s: %w", regType.String(), err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("unable to retrieve registry %s: %w", regType.String(), err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("server returned non-200 status code: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read response for registry %s: %w", regType.String(), err)
		}

		var registry Registry
		if err = json.Unmarshal(body, &registry); err != nil {
			return err
		}

		c.registries[regType] = &registry

	}

	return nil
}

func (c *Client) FetchRegistryByType(ctx context.Context, regType RegistryType, forceUpdate bool) (*Registry, error) {
	if c.registries[regType] != nil && !forceUpdate {
		return c.registries[regType], nil
	}
	var registry Registry

	req, err := http.NewRequestWithContext(ctx, "GET", regType.ServiceRegistryIndexURL(c.serviceRegistryIndexURL), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request for registry %s: %w", regType.String(), err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve registry %s: %w", regType.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned non-200 status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response for registry %s: %w", regType.String(), err)
	}

	if err = json.Unmarshal(body, &registry); err != nil {
		return nil, err
	}

	c.registries[regType] = &registry

	return c.registries[regType], nil
}

func (c *Client) GetDomainRDAPServers(ctx context.Context, domain string) ([]*url.URL, error) {
	var err error
	if c.registries[DNS] == nil {
		c.registries[DNS], err = c.FetchRegistryByType(ctx, DNS, true)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch DNS service registry: %w", err)
		}
	}
	return c.registries[DNS].getDNSServers(domain)
}

func (c *Client) GetAutnumRDAPServers(ctx context.Context, asn string) ([]*url.URL, error) {
	var err error
	if c.registries[ASN] == nil {
		c.registries[ASN], err = c.FetchRegistryByType(ctx, ASN, true)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch ASN service registry: %w", err)
		}
	}
	return c.registries[ASN].getASNServers(asn)
}

func (c *Client) GetIPAddressRDAPServers(ctx context.Context, ip string) ([]*url.URL, error) {
	var err error
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return nil, fmt.Errorf("input %s is not an IP Address", ip)
	}
	// IPv4 address
	if ipAddress.To4() != nil {
		if c.registries[IPv4] == nil {
			c.registries[IPv4], err = c.FetchRegistryByType(ctx, IPv4, true)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch IPv4 service registry: %w", err)
			}
		}
		return c.registries[IPv4].getNetServers(ipAddress)
	} else { // IPv6 address
		if c.registries[IPv6] == nil {
			c.registries[IPv6], err = c.FetchRegistryByType(ctx, IPv6, true)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch IPv6 service registry: %w", err)
			}
		}
		return c.registries[IPv6].getNetServers(ipAddress)
	}

}
