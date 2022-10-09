package rules

import (
	"testing"
)

func TestLifecycleRule(t *testing.T) {
	runTests(t,
		NewLifecycleRule(),
	)
}
