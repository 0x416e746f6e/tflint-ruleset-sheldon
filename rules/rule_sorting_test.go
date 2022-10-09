package rules

import (
	"testing"
)

func TestSortingRule(t *testing.T) {
	runTests(t,
		NewSortingRule(),
	)
}
