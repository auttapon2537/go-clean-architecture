package docs

import "embed"

//go:embed openapi.yaml openapi.json
var OpenAPIFS embed.FS
