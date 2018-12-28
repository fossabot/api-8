package api

// ContentTypeJSON represents content type value for json
const ContentTypeJSON = "application/json"

// Handler represents an api handler
type Handler func(Request) Response
