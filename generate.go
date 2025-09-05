package ogen

// official //go run github.com/ogen-go/ogen/cmd/ogen@latest --target petstore --clean petstore.yml

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --clean --target internal/ogen --package ogen api/openapi.yaml
