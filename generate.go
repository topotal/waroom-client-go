package ogen

//go:generate go run ./internal/cmd/genwaroom analyze -spec api/openapi.yaml -ogen-config .ogen.yaml -ogen-target internal/ogen -config .genwaroom.yaml -o api/sdk-model.yaml
//go:generate go run ./internal/cmd/genwaroom render -model api/sdk-model.yaml -out ./waroom
