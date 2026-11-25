package waroom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTimezone(t *testing.T) {
	valid := []string{
		"Asia/Tokyo",
		"America/New_York",
		"UTC",
		"Etc/UTC",
		"Etc/GMT+9",
		"Japan",
		"US/Eastern",
	}
	for _, name := range valid {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, ValidateTimezone(name))
		})
	}

	invalid := []string{
		"",
		"Local",
		"Tokyo",
		"Invalid/Timezone",
		"Asia/Tokyo ",
		" Asia/Tokyo",
		"+09:00",
		"../etc/passwd",
		"/usr/share/zoneinfo/Asia/Tokyo",
		"zone.tab",
	}
	for _, name := range invalid {
		t.Run("invalid "+name, func(t *testing.T) {
			assert.Error(t, ValidateTimezone(name))
		})
	}
}
