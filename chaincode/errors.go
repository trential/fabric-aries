package main

import "errors"

var (
	ErrInvalidRequest       = errors.New("invalid request")
	ErrFunctionNotSupported = errors.New("function not supported")
	ErrVerifySignature      = errors.New("failed to verify signature")
	ErrWorldstateRead       = errors.New("failed to read worldstate")
	ErrWorldstateWrite      = errors.New("failed to write worldstate")
	ErrStateNotFound        = errors.New("state not found")
	ErrUnimplemented        = errors.New("unimplemented")
)
