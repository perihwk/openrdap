# openrdap

openrdap is a Go library designed to interact with RDAP (Registration Data Access Protocol) servers. It provides a simple and efficient way to query domain, IP, and ASN registration data from RDAP-compliant registries. With a modern and easy-to-use API, this library simplifies the process of querying and parsing RDAP responses in Go.
Features

    RDAP Queries: Perform RDAP queries for domain names, IP addresses, and ASNs.
    Automatic JSON Parsing: Parse RDAP JSON responses into native Go structures.
    Multiple Registry Support: Interact with various global RDAP registries.
    Error Handling: Built-in error handling and validation for RDAP responses.
    Go 1.22+: Designed to take advantage of features in Go 1.22 and later.

Requirements

    Go 1.22+ (Should work on older versions but has only been tested on 1.22+)

Installation

To install the openrdap library in your Go project, run:

```go get github.com/perihwk/openrdap```

Alternatively, you can manually clone the repository:

```
git clone https://github.com/perihwk/openrdap.git
cd openrdap
go mod tidy
make build
```

## Example Usage
### Query a Domain

```package main

import (
    "fmt"
    "log"

    "github.com/perihwk/openrdap"
)

func main() {
    rdapClient := openrdap.NewClient(httpClient, nil)

    ctx := context.Background()

    domain, err := rdapClient.GetRDAPFromDomain(ctx, "example.com")
    if err != nil {
        log.Fatalf("Error querying domain: %v", err)
    }

    fmt.Printf("Domain Info: %+v\n", domain)
}```
