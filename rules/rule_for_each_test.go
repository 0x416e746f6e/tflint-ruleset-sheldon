package rules

import (
	"testing"
)

func TestForEachRule(t *testing.T) {
	runTests(t,
		NewForEachRule(),
	)
}
