// Package config parses configuration documents without requiring a schema.
package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

// Format identifies the syntax detected in a document.
type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
)

// Document holds a single configuration document. Root preserves scalar tags,
// mapping order, aliases, and source positions. JSON uses the same YAML node
// representation, as JSON syntax is supported by the YAML decoder.
type Document struct {
	Format Format
	Root   *yaml.Node
}

// Parse reads one JSON or YAML document. Any root type is accepted, including
// scalars and null. Valid JSON takes precedence; otherwise YAML is attempted.
// Empty input and streams containing multiple documents are rejected.
// The caller owns r and remains responsible for closing it.
func Parse(r io.Reader) (*Document, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}

	format := YAML
	if json.Valid(data) {
		format = JSON
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var node yaml.Node
	if err := decoder.Decode(&node); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("configuration is empty")
		}
		return nil, fmt.Errorf("parse %s configuration: %w", format, err)
	}

	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return nil, fmt.Errorf("parse trailing configuration: %w", err)
		}
		return nil, errors.New("expected a single configuration document")
	}

	if err := validate(node.Content[0], make(map[*yaml.Node]bool), make(map[*yaml.Node]bool)); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	return &Document{Format: format, Root: node.Content[0]}, nil
}

// Validate the graph without converting scalars to Go types: configuration
// numbers may exceed Go's numeric ranges, and custom tags must remain intact.
func validate(node *yaml.Node, active, done map[*yaml.Node]bool) error {
	if active[node] {
		return fmt.Errorf("recursive alias at line %d", node.Line)
	}
	if done[node] {
		return nil
	}
	active[node] = true
	defer delete(active, node)

	if node.Kind == yaml.MappingNode {
		type scalarKey struct {
			tag   string
			value string
		}
		seen := make(map[scalarKey]bool)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			if key.Kind != yaml.ScalarNode {
				continue
			}
			id := scalarKey{tag: key.Tag, value: key.Value}
			if seen[id] {
				return fmt.Errorf("duplicate mapping key at line %d", key.Line)
			}
			seen[id] = true
		}
	}
	if node.Kind == yaml.AliasNode {
		if err := validate(node.Alias, active, done); err != nil {
			return err
		}
	}
	for _, child := range node.Content {
		if err := validate(child, active, done); err != nil {
			return err
		}
	}
	done[node] = true
	return nil
}
