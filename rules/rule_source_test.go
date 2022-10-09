package rules

import (
	"testing"
)

func TestSourceRule(t *testing.T) {
	runTests(t,
		NewSourceRule(),
	)
}
