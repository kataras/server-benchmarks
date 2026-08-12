package report

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// templates holds the parsed report templates. It is initialized once and
// never mutated afterwards, so it is safe for concurrent use.
var templates = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))
