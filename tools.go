//go:build tools

package tools

import (
	_ "github.com/ogen-go/ogen/cmd/ogen"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
