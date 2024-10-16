package openrdap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/perihwk/openrdap/bootstrap"
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
	httpClient      *http.Client
	bootstrapClient *bootstrap.Client
}

func NewClient(
	httpClient *http.Client,
	bootstrapClient *bootstrap.Client,
) *Client {

	if bootstrapClient == nil {
		bootstrapClient = bootstrap.NewBootstrapClient(httpClient, "")
	}

	return &Client{
		httpClient:      httpClient,
		bootstrapClient: bootstrapClient,
	}
}

func (c *Client) GetRDAPInfoFromServer(ctx context.Context, rdapServer, query string, searchType RegistrySearchType) (any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rdapServer+fmt.Sprintf(searchType.Path(), query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("server %s returned non-200 status code: %s", rdapServer, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	var result interface{}
	switch searchType {
	case DNS:
		result = &Domain{}
	case IPv4, IPv6:
		result = &IPNetwork{}
	case ASN:
		result = &Autnum{}
	default:
		return nil, fmt.Errorf("unsupported search type")
	}

	if err = json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("error parsing RDAP response: %w", err)
	}

	return result, nil
}

func (c *Client) GetRDAPFromDomain(ctx context.Context, domain string) (*Domain, error) {
	registryServers, err := c.bootstrapClient.GetDomainRDAPServers(ctx, domain)
	if err != nil {
		return nil, err
	}

	var domainResp *Domain
	for i, u := range registryServers {
		localSrv := u.String()
		// use first https RDAP server. If no https server then use whatever the last option was
		if u.Scheme == "https" || i == len(registryServers)-1 {
			if localSrv, err = url.JoinPath(localSrv, "domain", domain); err != nil {
				return nil, err
			}

			req, err := http.NewRequestWithContext(ctx, "GET", localSrv, nil)
			if err != nil {
				return nil, err
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server %s returned non-200 status code: %s", localSrv, resp.Status)
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

func (c *Client) GetRDAPFromIP(ctx context.Context, ip string) (*IPNetwork, error) {
	registryServers, err := c.bootstrapClient.GetIPAddressRDAPServers(ctx, ip)
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
			req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
			if err != nil {
				return nil, err
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server %s returned non-200 status code: %s", u.Path, resp.Status)
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

func (c *Client) GetRDAPFromAutnum(ctx context.Context, asn string) (*Autnum, error) {
	registryServers, err := c.bootstrapClient.GetAutnumRDAPServers(ctx, asn)
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

			req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
			if err != nil {
				return nil, err
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("server %s returned non-200 status code: %s", u.Path, resp.Status)
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
