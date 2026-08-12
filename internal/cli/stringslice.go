package cli

import "strings"

// stringSlice is a flag.Value that collects every occurrence of a
// repeatable flag, e.g. -t rest -t static.iris.
type stringSlice []string

// String implements flag.Value.
func (s *stringSlice) String() string { return strings.Join(*s, ", ") }

// Set implements flag.Value.
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}
