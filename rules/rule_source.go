package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// SourceRule makes sure that `source` meta-attribute is always on top.
type SourceRule struct {
	tflint.DefaultRule
}

// NewSourceRule creates a new SourceRule.
func NewSourceRule() *SourceRule {
	return &SourceRule{}
}

// Name returns the name of the rule.
func (r *SourceRule) Name() string {
	return project.RuleName("source")
}

// Enabled returns whether the rule is enabled by default.
func (r *SourceRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *SourceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether `source` clause is placed on the top of the module
// definition.
func (r *SourceRule) Check(runner tflint.Runner) error {
	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "module" {
			return nil
		}

		source, exists := b.Body.Attributes["source"]
		if !exists {
			return nil
		}
		if node.FirstNodeFrom(b.Body).AsAttribute() == source {
			return nil
		}

		return runner.EmitIssue(
			r,
			"`source` must be the top-most attribute",
			source.SrcRange,
		)
	})
}
