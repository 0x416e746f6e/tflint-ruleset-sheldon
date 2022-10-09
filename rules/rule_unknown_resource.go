package rules

import (
	"fmt"
	"reflect"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// UnknownResourceRule warns if the linter encounters unknown resource.
type UnknownResourceRule struct {
	tflint.DefaultRule
}

// NewUnknownResourceRule creates a new UnknownResourceRule.
func NewUnknownResourceRule() *UnknownResourceRule {
	return &UnknownResourceRule{}
}

// Name returns the name of the rule.
func (r *UnknownResourceRule) Name() string {
	return project.RuleName("unknown_resource")
}

// Enabled returns whether the rule is enabled by default.
func (r *UnknownResourceRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *UnknownResourceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the reference link for the rule.
func (r *UnknownResourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the key-attributes (those that uniquely identify the
// resource) are put on top of the resource definition.
func (r *UnknownResourceRule) Check(rr tflint.Runner) error {
	runner, ok := rr.(*custom.Runner)
	if !ok {
		return fmt.Errorf("unexpected runner type: %s", reflect.TypeOf(rr))
	}

	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" && b.Type != "data" {
			return nil
		}

		kind := b.Labels[0]
		if _, kindIsKnown := runner.Resources[kind]; kindIsKnown {
			return nil
		}

		return runner.EmitIssue(
			r,
			fmt.Sprintf("key-attributes for resource type `%s` are not configured", kind),
			b.LabelRanges[0],
		)
	})
}
