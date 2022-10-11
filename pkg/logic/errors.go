package logic

import "errors"

var (
	ErrTableRequired   = errors.New("logic: nft inet table name required")
	ErrSetNotFound     = errors.New("logic: set not found")
	ErrIp6NotSupported = errors.New("logic: ip6 not supported")
)
