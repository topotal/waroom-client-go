package waroom

import (
	"errors"
	"fmt"
	"time"

	// Embed the IANA time zone database so that ValidateTimezone works on
	// hosts without zoneinfo files (e.g. Windows). The system database is
	// still preferred when it exists.
	_ "time/tzdata"
)

// ValidateTimezone checks that name is an IANA time zone name (e.g.
// "Asia/Tokyo") that can be sent as the timezone of an on-call escalation
// policy or rotation. It returns nil for a valid name.
//
// Aliases such as "Japan" or "Etc/UTC" are valid, and the API stores and
// returns them as sent. The empty string and "Local" are rejected, although
// time.LoadLocation accepts them.
//
// The check is best effort: the API remains the source of truth. On a
// case-insensitive file system (e.g. macOS), a name with the wrong case such
// as "asia/tokyo" passes here but is rejected by the API.
func ValidateTimezone(name string) error {
	switch name {
	case "":
		return errors.New("waroom: invalid time zone: must not be empty")
	case "Local":
		return fmt.Errorf("waroom: invalid time zone %q: must be an IANA time zone name", name)
	}
	if _, err := time.LoadLocation(name); err != nil {
		return fmt.Errorf("waroom: invalid time zone %q: %w", name, err)
	}
	return nil
}
