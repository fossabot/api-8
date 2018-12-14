package api

import (
	"context"
)

// ContentTypeJSON represents content type value for json
const ContentTypeJSON = "application/json"

// Handler represents an api handler
type Handler func(context.Context, Request) Response
