package rules

import (
	"testing"
)

func TestCountRule(t *testing.T) {
	runTests(t,
		NewCountRule(),
	)
}
