package logic

import "errors"

var (
	ErrTableRequired = errors.New("logic: nft inet table name required")
	ErrSetNotFound   = errors.New("logic: set not found")
)
