package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecVersion(t *testing.T) {
	v, err := specVersion([]byte("openapi: 3.1.0\ninfo:\n  title: Waroom\n  version: 1.2.3\n"))
	require.NoError(t, err)
	assert.Equal(t, "1.2.3", v)
}

func TestSpecVersionEmpty(t *testing.T) {
	_, err := specVersion([]byte("openapi: 3.1.0\ninfo:\n  title: Waroom\n"))
	assert.Error(t, err)
}

func TestRender(t *testing.T) {
	src, err := render("1.2.3")
	require.NoError(t, err)
	assert.Contains(t, string(src), `const SpecVersion = "1.2.3"`)
}
