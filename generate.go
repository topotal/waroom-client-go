package ogen

//go:generate go tool ogen --clean --target internal/ogen --package ogen api/openapi.yaml
//go:generate go run ./internal/cmd/specversion api/openapi.yaml internal/ogen/spec_version_gen.go
