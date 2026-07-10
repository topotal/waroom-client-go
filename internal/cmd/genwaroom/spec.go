package main

import (
	"fmt"
	"os"

	"github.com/go-faster/yaml"
)

var httpMethods = []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}

// loadTags extracts operationId → first tag from the OpenAPI spec. Tags
// are the only spec input; everything else comes from the ogen package.
func loadTags(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Paths map[string]map[string]yaml.Node `yaml:"paths"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	tags := map[string]string{}
	for p, item := range doc.Paths {
		for _, method := range httpMethods {
			node, ok := item[method]
			if !ok {
				continue
			}
			var op struct {
				OperationID string   `yaml:"operationId"`
				Tags        []string `yaml:"tags"`
			}
			if err := node.Decode(&op); err != nil {
				return nil, fmt.Errorf("parse %s %s: %w", method, p, err)
			}
			if op.OperationID == "" {
				return nil, fmt.Errorf("%s %s: missing operationId", method, p)
			}
			if len(op.Tags) == 0 {
				continue // fall back to path-segment grouping in load.go
			}
			tags[op.OperationID] = op.Tags[0]
		}
	}
	return tags, nil
}
