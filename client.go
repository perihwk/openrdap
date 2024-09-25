package openrdap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetRDAPInfoFromServer(rdapServer, query string, searchType RegistrySearchType) (any, error) {
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

	fmt.Printf("body: %s\n", body)

	return nil, nil
}

func (c *Client) GetRDAPFromDomain(query string) (*Domain, error) {
	registryServers, err := c.bootstrapClient.GetDomainRDAPServers(query)
	if err != nil {
		return nil, err
	}

	var domainResp *Domain
	for _, u := range registryServers {
		u.Path = u.Path + "/domain/" + query
		resp, err := c.httpClient.Get(u.String())
		if err != nil {
			fmt.Println(err)
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
	return domainResp, nil
}
