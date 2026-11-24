package waroom

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestAPISpecVersion catches a spec sync that was not followed by go generate.
func TestAPISpecVersion(t *testing.T) {
	b, err := os.ReadFile("../api/openapi.yaml")
	require.NoError(t, err)
	var doc struct {
		Info struct {
			Version string `yaml:"version"`
		} `yaml:"info"`
	}
	require.NoError(t, yaml.Unmarshal(b, &doc))
	assert.Equal(t, doc.Info.Version, APISpecVersion)
}
