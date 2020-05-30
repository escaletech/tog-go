package flags

import (
	"fmt"
	"regexp"
)

var validIdentifier = regexp.MustCompile(`^[\w_-]{1,30}$`)

func validate(f Flag) error {
	if !validIdentifier.Match([]byte(f.Namespace)) {
		return fmt.Errorf("invalid flag namespace, must conform to %v", validIdentifier.String())
	}

	if !validIdentifier.Match([]byte(f.Name)) {
		return fmt.Errorf("invalid flag name, must conform to %v", validIdentifier.String())
	}

	return nil
}
