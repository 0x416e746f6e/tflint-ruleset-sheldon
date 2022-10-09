package rules

import (
	"testing"
)

func TestDependsOnRule(t *testing.T) {
	runTests(t,
		NewDependsOnRule(),
	)
}
