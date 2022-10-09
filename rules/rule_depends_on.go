package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DependsOnRule makes sure that `depends_on` clause is always the last.
type DependsOnRule struct {
	tflint.DefaultRule
}

// NewDependsOnRule creates a new DependsOnRule.
func NewDependsOnRule() *DependsOnRule {
	return &DependsOnRule{}
}

// Name returns the name of the rule.
func (r *DependsOnRule) Name() string {
	return project.RuleName("depends_on")
}

// Enabled returns whether the rule is enabled by default.
func (r *DependsOnRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *DependsOnRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *DependsOnRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether `depends_on` clause is placed at the end of the
// resource definition.
func (r *DependsOnRule) Check(runner tflint.Runner) error {
	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" && b.Type != "data" {
			return nil
		}

		dependsOn, exists := b.Body.Attributes["depends_on"]
		if !exists {
			return nil
		}
		if node.LastNodeFrom(b.Body).AsAttribute() == dependsOn {
			return nil
		}

		return runner.EmitIssue(
			r,
			"`depends_on` clause must be the last one in the definition",
			dependsOn.SrcRange,
		)
	})
}
