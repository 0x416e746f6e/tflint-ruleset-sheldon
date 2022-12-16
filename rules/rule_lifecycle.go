package rules

import (
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// LifecycleRule makes sure that `lifecycle` clause is the last one before
// `depends_on`.
type LifecycleRule struct {
	tflint.DefaultRule
}

// NewLifecycleRule creates a new LifecycleRule.
func NewLifecycleRule() *LifecycleRule {
	return &LifecycleRule{}
}

// Name returns the name of the rule.
func (r *LifecycleRule) Name() string {
	return project.RuleName("lifecycle")
}

// Enabled returns whether the rule is enabled by default.
func (r *LifecycleRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *LifecycleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *LifecycleRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether `lifecycle` clause is placed at the end of the
// resource definition (but before `depends_on`).
func (r *LifecycleRule) Check(runner tflint.Runner) error {
	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" {
			return nil
		}

		var lifecycle *hclsyntax.Block
		for _, block := range b.Body.Blocks {
			if block.Type == "lifecycle" {
				if lifecycle != nil {
					return runner.EmitIssue(
						r,
						"more than 1 `lifecycle` block found",
						block.TypeRange,
					)
				}
				lifecycle = block
			}
		}

		if lifecycle == nil {
			return nil
		}

		nodes := node.OrderedInspectableNodesFrom(b.Body)

		n := nodes[len(nodes)-1]
		if n.IsAttribute() && n.Name() == "depends_on" {
			n = nodes[len(nodes)-2]
		}
		if n.AsBlock() == lifecycle {
			return nil
		}

		return runner.EmitIssue(
			r,
			"`lifecycle` block must be at the end of the definition (but before `depends_on`)",
			lifecycle.Body.SrcRange,
		)
	})
}
