package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// CountRule makes sure that `count` attribute is always on top.
type CountRule struct {
	tflint.DefaultRule
}

// NewCountRule creates a new CountRule.
func NewCountRule() *CountRule {
	return &CountRule{}
}

// Name returns the name of the rule.
func (r *CountRule) Name() string {
	return project.RuleName("count")
}

// Enabled returns whether the rule is enabled by default.
func (r *CountRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *CountRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *CountRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether `count` clause is placed on the top of the resource
// definition.
func (r *CountRule) Check(runner tflint.Runner) error {
	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" && b.Type != "data" {
			return nil
		}

		count, exists := b.Body.Attributes["count"]
		if !exists {
			return nil
		}
		if node.FirstNodeFrom(b.Body).AsAttribute() == count {
			return nil
		}

		return runner.EmitIssue(
			r,
			"`count` must be the top-most attribute",
			count.SrcRange,
		)
	})
}
