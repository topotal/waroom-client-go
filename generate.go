package ogen

//go:generate go run ./internal/cmd/downconvert -in api/openapi.yaml -out api/openapi-3.0.yaml
//go:generate go tool ogen --clean --target internal/ogen --package ogen api/openapi-3.0.yaml
//go:generate go run ./internal/cmd/genwaroom analyze -spec api/openapi.yaml -src ./internal/ogen -config .genwaroom.yaml -o api/sdk-model.yaml
//go:generate go run ./internal/cmd/genwaroom render -model api/sdk-model.yaml -out ./waroom
