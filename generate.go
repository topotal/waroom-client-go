package ogen

//go:generate go run ./internal/cmd/downconvert -in api/openapi.yaml -out api/openapi-3.0.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --clean --target internal/ogen --package ogen api/openapi-3.0.yaml
