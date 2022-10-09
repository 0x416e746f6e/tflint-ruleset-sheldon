// Package rules ...
package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

// All returns all the rules from the package.
func All() []tflint.Rule {
	return []tflint.Rule{
		NewCountRule(),
		NewDependsOnRule(),
		NewForEachRule(),
		NewKeyAttributesRule(),
		NewLifecycleRule(),
		NewSortingRule(),
		NewSourceRule(),
		NewSpacingRule(),
		NewUnknownResourceRule(),
	}
}
