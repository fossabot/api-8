package testhelper

import "context"

// NewContext creates context for testing
func NewContext() context.Context {
	ctx := context.Background()
	return ctx
}
