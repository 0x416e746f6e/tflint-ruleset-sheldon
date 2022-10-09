package rules

import (
	"testing"
)

func TestKeyAttributesRule(t *testing.T) {
	runTests(t,
		NewKeyAttributesRule(),
	)
}
