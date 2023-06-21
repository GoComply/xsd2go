package templates

import "embed"

//go:embed all:*/*.tmpl */*.incl */*.elem
var Templates embed.FS
