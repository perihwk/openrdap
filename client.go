package openrdap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/perihwk/openrdap/bootstrap"
)

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

func (c *Client) doRequest(ctx context.Context, query string, regType bootstrap.RegistryType, registryServers []*url.URL) ([]byte, error) {
	for i, u := range registryServers {
		localSrv := u.String()
		// use first https RDAP server. If no https server then use whatever the last option was
		if u.Scheme == "https" || i == len(registryServers)-1 {
			localSrv, err := url.JoinPath(localSrv, regType.PathSegment(), query)
			if err != nil {
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

			return body, nil
		}
	}

	return nil, fmt.Errorf("no RDAP servers available")
}

func (c *Client) GetRDAPInfoFromServer(ctx context.Context, rdapServer, query string, searchType bootstrap.RegistryType) (any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rdapServer+fmt.Sprintf(searchType.Path(), query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server %s returned non-200 status code: %s", rdapServer, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result any
	switch searchType {
	case bootstrap.DNS:
		result = &Domain{}
	case bootstrap.IPv4, bootstrap.IPv6:
		result = &IPNetwork{}
	case bootstrap.ASN:
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
	body, err := c.doRequest(ctx, domain, bootstrap.DNS, registryServers)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &domainResp); err != nil {
		return nil, fmt.Errorf("error parsing RDAP response: %w", err)
	}
	return domainResp, nil
}

func (c *Client) GetRDAPFromIP(ctx context.Context, ip string) (*IPNetwork, error) {
	registryServers, err := c.bootstrapClient.GetIPAddressRDAPServers(ctx, ip)
	if err != nil {
		return nil, err
	}

	var ipAddressResp *IPNetwork
	body, err := c.doRequest(ctx, ip, bootstrap.IPv4, registryServers)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &ipAddressResp); err != nil {
		return nil, fmt.Errorf("error parsing RDAP response: %w", err)
	}
	return ipAddressResp, nil
}

func (c *Client) GetRDAPFromAutnum(ctx context.Context, asn string) (*Autnum, error) {
	registryServers, err := c.bootstrapClient.GetAutnumRDAPServers(ctx, asn)
	if err != nil {
		return nil, err
	}

	asn = strings.TrimPrefix(strings.ToUpper(asn), "AS")

	var autnumResp *Autnum
	body, err := c.doRequest(ctx, asn, bootstrap.ASN, registryServers)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &autnumResp); err != nil {
		return nil, fmt.Errorf("error parsing RDAP response: %w", err)
	}
	return autnumResp, nil
}
