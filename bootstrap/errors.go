package bootstrap

import (
	"errors"
)

var (
	ErrInvalidCIDR      = errors.New("invalid CIDR block")
	ErrRDAPNotSupported = errors.New("RDAP server not found")
	ErrInvalidASNRange  = errors.New("invalid ASN range")
)
