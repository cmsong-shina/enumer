package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// openAPITypeInfo holds enum information needed for OpenAPI schema generation.
type openAPITypeInfo struct {
	typeName string
	values   []string // transformed enum value names
}

// generateOpenAPIYAML produces OpenAPI YAML schema content for the given enum type.
// If codegen is true, oapi-codegen extensions (x-go-type, x-go-type-import) are included.
// This is a pure function with no side effects on Generator state.
func (g *Generator) generateOpenAPIYAML(
	values []Value,
	typeName string,
	outPath string,
	codegen bool,
) {
	var b strings.Builder

	b.WriteString("type: string\n")
	// Write enum field (available values) as JSON representaion
	strRepr := []string{}
	for _, v := range values {
		strRepr = append(strRepr, v.name)
	}
	jsonRepr, err := json.Marshal(strRepr)
	if err != nil {
		log.Fatalf("Failed to encoding enum values: %s", err)
	}
	b.WriteString(fmt.Sprintf("enum: %s\n", jsonRepr))

	// Write x-go-type extension
	if codegen {
		b.WriteString(fmt.Sprintf("x-go-type: %s.%s\n", g.pkg.name, typeName))
		b.WriteString("x-go-type-import:\n")
		b.WriteString(fmt.Sprintf("  path: %s\n", g.pkg.typesPkg.Path()))
		b.WriteString(fmt.Sprintf("  name: %s\n", g.pkg.name))
	}

	writeOpenAPIFile(outPath, []byte(b.String()))
}

// writeOpenAPIFile writes the given content to the specified path atomically.
func writeOpenAPIFile(path string, content []byte) {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("creating directory for OpenAPI output: %s", err)
		}
	}

	tmpFile, err := os.CreateTemp(dir, "openapi_")
	if err != nil {
		log.Fatalf("creating temporary file for OpenAPI output: %s", err)
	}
	_, err = tmpFile.Write(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		log.Fatalf("writing OpenAPI output: %s", err)
	}
	tmpFile.Close()

	err = os.Rename(tmpFile.Name(), path)
	if err != nil {
		log.Fatalf("moving tempfile to OpenAPI output file: %s", err)
	}
}
