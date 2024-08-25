package templates

import "embed"

//go:embed service.go.tmpl
var ServiceFile embed.FS

//go:embed handler.go.tmpl
var HandlerFile embed.FS
