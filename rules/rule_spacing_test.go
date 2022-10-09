package rules

import (
	"testing"
)

func TestSpacingRule(t *testing.T) {
	runTests(t,
		NewSpacingRule(),
	)
}
