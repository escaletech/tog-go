package flags

import (
	"fmt"
	"strings"
)

type ClientOptions struct {
	Addr    string
	Cluster bool
}

type Flag struct {
	Namespace   string    `json:"namespace,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Timestamp   int64     `json:"timestamp,omitempty"`
	Rollout     []Rollout `json:"rollout,omitempty"`
}

type Rollout struct {
	Value      bool     `json:"value,omitempty"`
	Percentage *int     `json:"percentage,omitempty"`
	Traits     []string `json:"traits,omitempty"`
}

type MultiError []error

func (me MultiError) Error() string {
	msgs := make([]string, len(me))
	for i, err := range me {
		msgs[i] = fmt.Sprintf("  %v. %v", i, err.Error())
	}

	return fmt.Sprintf("multiple errors:\n%v", strings.Join(msgs, "\n"))
}
