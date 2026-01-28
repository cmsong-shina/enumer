package main

import (
	"bytes"
	"fmt"
	"strings"
)

type openAPITypeInfo struct {
	typeName string
	values   []string // transformed enum value names
}

// generateOpenAPIYAML produces OpenAPI YAML schema content for the collected enum types.
// If codegen is true, oapi-codegen extensions (x-go-type, x-go-type-import) are included.
func (g *Generator) generateOpenAPIYAML(codegen bool) []byte {
	var buf bytes.Buffer

	for i, info := range g.openAPITypes {
		if i > 0 {
			buf.WriteString("---\n")
		}

		buf.WriteString("type: string\n")
		buf.WriteString(fmt.Sprintf("enum: [%s]\n", strings.Join(info.values, ", ")))

		if codegen {
			pkgName := g.pkg.name
			pkgPath := ""
			if g.pkg.typesPkg != nil {
				pkgPath = g.pkg.typesPkg.Path()
			}
			buf.WriteString(fmt.Sprintf("x-go-type: %s.%s\n", pkgName, info.typeName))
			buf.WriteString("x-go-type-import:\n")
			buf.WriteString(fmt.Sprintf("  path: %s\n", pkgPath))
			buf.WriteString(fmt.Sprintf("  name: %s\n", pkgName))
		}
	}

	return buf.Bytes()
}
