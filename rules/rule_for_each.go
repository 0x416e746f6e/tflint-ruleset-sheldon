package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ForEachRule makes sure that `for_each` meta-attribute is always on top.
type ForEachRule struct {
	tflint.DefaultRule
}

// NewForEachRule creates a new ForEachRule.
func NewForEachRule() *ForEachRule {
	return &ForEachRule{}
}

// Name returns the name of the rule.
func (r *ForEachRule) Name() string {
	return project.RuleName("for_each")
}

// Enabled returns whether the rule is enabled by default.
func (r *ForEachRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *ForEachRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *ForEachRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether `for_each` clause is placed on the top of the resource
// definition.
func (r *ForEachRule) Check(runner tflint.Runner) error {
	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" && b.Type != "data" {
			return nil
		}

		forEach, exists := b.Body.Attributes["for_each"]
		if !exists {
			return nil
		}
		if node.FirstNodeFrom(b.Body).AsAttribute() == forEach {
			return nil
		}

		return runner.EmitIssue(
			r,
			"`for_each` must be the top-most attribute",
			forEach.SrcRange,
		)
	})
}
