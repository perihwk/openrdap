package openrdap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"openrdap/bootstrap"
)

// A RegistrySearch represents a registry search type.
type RegistrySearchType int

const (
	DNS RegistrySearchType = iota
	IPv4
	IPv6
	ASN
	ENTITY
)

func (r RegistrySearchType) Path() string {
	switch r {
	case DNS:
		return "domain/%s"
	case IPv4:
		return "ip/%s"
	case IPv6:
		return "ip/%s"
	case ASN:
		return "autnum/%s"
	case ENTITY:
		return "entity/%s"
	default:
		panic("Unknown RegistrySearchType")
	}
}

type Client struct {
	ctx             *context.Context
	httpClient      *http.Client
	bootstrapClient *bootstrap.Client
}

func NewClient(
	ctx *context.Context,
	httpClient *http.Client,
	bootstrapClient *bootstrap.Client,
) *Client {

	if bootstrapClient == nil {
		bootstrapClient = bootstrap.NewBootstrapClient(ctx, httpClient, "")
	}

	return &Client{
		ctx:             ctx,
		httpClient:      httpClient,
		bootstrapClient: bootstrapClient,
	}
}

func (c *Client) GetRDAPInfoFromServer(rdapServer, query string, searchType RegistrySearchType) (*Domain, error) {
	resp, err := c.httpClient.Get(rdapServer + fmt.Sprintf(searchType.Path(), query))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("server returned non-200 status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	var domainResp *Domain
	if err = json.Unmarshal(body, &domainResp); err != nil {
		return nil, fmt.Errorf("error parsing RDAP response: %w", err)
	}

	return domainResp, nil
}

func (c *Client) GetRDAPFromDomain(domain string) (*Domain, error) {
	registryServers, err := c.bootstrapClient.GetDomainRDAPServers(domain)
	if err != nil {
		return nil, err
	}

	var domainResp *Domain
	for i, u := range registryServers {
		// use first https RDAP server. If no https server then use whatever the last option was
		if u.Scheme == "https" || i == len(registryServers)-1 {
			if u.Path, err = url.JoinPath(u.Path, "domain", domain); err != nil {
				return nil, err
			}
			resp, err := c.httpClient.Get(u.String())
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server returned non-200 status code: %s", resp.Status)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading RDAP response: %w", err)
			}

			if err = json.Unmarshal(body, &domainResp); err != nil {
				return nil, fmt.Errorf("error parsing RDAP response: %w", err)
			}
		}
	}
	return domainResp, nil
}

func (c *Client) GetRDAPFromIP(ip string) (*IPNetwork, error) {
	registryServers, err := c.bootstrapClient.GetIPAddressRDAPServers(ip)
	if err != nil {
		return nil, err
	}

	var ipAddressResp *IPNetwork
	for i, u := range registryServers {
		// use first https RDAP server. If no https server then use whatever the last option was
		if u.Scheme == "https" || i == len(registryServers)-1 {
			if u.Path, err = url.JoinPath(u.Path, "ip", ip); err != nil {
				return nil, err
			}
			resp, err := c.httpClient.Get(u.String())
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server returned non-200 status code: %s", resp.Status)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading RDAP response: %w", err)
			}

			if err = json.Unmarshal(body, &ipAddressResp); err != nil {
				return nil, fmt.Errorf("error parsing RDAP response: %w", err)
			}
		}
	}
	return ipAddressResp, nil
}

func (c *Client) GetRDAPFromAutnum(asn string) (*Autnum, error) {
	registryServers, err := c.bootstrapClient.GetAutnumRDAPServers(asn)
	if err != nil {
		return nil, err
	}

	var autnumResp *Autnum
	for i, u := range registryServers {
		// use first https RDAP server. If no https server then use whatever the last option was
		if u.Scheme == "https" || i == len(registryServers)-1 {
			if u.Path, err = url.JoinPath(u.Path, "autnum", asn); err != nil {
				return nil, err
			}
			resp, err := c.httpClient.Get(u.String())
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server returned non-200 status code: %s", resp.Status)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading RDAP response: %w", err)
			}

			if err = json.Unmarshal(body, &autnumResp); err != nil {
				return nil, fmt.Errorf("error parsing RDAP response: %w", err)
			}
		}
	}
	return autnumResp, nil
}
