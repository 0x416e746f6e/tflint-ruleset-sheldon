package rules

import (
	"fmt"
	"reflect"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// KeyAttributesRule makes sure that key-attributes (those that uniquely
// identify the resource) are put on top of the resource definition.
type KeyAttributesRule struct {
	tflint.DefaultRule
}

// NewKeyAttributesRule creates a new KeyAttributesRule.
func NewKeyAttributesRule() *KeyAttributesRule {
	return &KeyAttributesRule{}
}

// Name returns the name of the rule.
func (r *KeyAttributesRule) Name() string {
	return project.RuleName("key_attributes")
}

// Enabled returns whether the rule is enabled by default.
func (r *KeyAttributesRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *KeyAttributesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *KeyAttributesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether the key-attributes (those that uniquely identify the
// resource) are put on top of the resource definition.
func (r *KeyAttributesRule) Check(rr tflint.Runner) error {
	runner, ok := rr.(*custom.Runner)
	if !ok {
		return fmt.Errorf("unexpected runner type: %s", reflect.TypeOf(rr))
	}

	return visit.Blocks(r, runner, func(b *hclsyntax.Block, _ []byte) error {
		if b.Type != "resource" && b.Type != "data" {
			return nil
		}

		// Get key attributes for the resource.
		kind := b.Labels[0]
		resource, kindIsKnown := runner.Resources[kind]
		if !kindIsKnown {
			return nil
		}
		knownKeyAttributes := resource.KeyAttributes

		if len(knownKeyAttributes) == 0 {
			return nil
		}

		body := b.Body
		for _, keyBlockName := range resource.KeyBlocks {
			for _, block := range body.Blocks {
				if block.Type != keyBlockName {
					continue
				}
				body = block.Body
			}
		}

		// Extract key attributes that are present in resource definition.
		kaList := make([]*hclsyntax.Attribute, 0, len(knownKeyAttributes))
		kaSet := map[string]struct{}{}
		for _, attrName := range knownKeyAttributes {
			attr, exists := body.Attributes[attrName]
			if !exists {
				continue
			}
			kaList = append(kaList, attr)
			kaSet[attrName] = struct{}{}
		}

		pos := 0
		for _, n := range node.OrderedInspecableNodesFrom(body) {
			if n.IsAttribute() && n.Name() == "for_each" {
				continue
			}
			if pos == len(kaList) {
				break
			}
			k := kaList[pos]
			if n.IsAttribute() && n.Name() != k.Name {
				if _, isKey := kaSet[n.Name()]; isKey {
					return runner.EmitIssue(
						r,
						fmt.Sprintf(
							"higher-priority key-attribute `%s` should be defined before `%s`",
							k.Name,
							n.Name(),
						),
						k.SrcRange,
					)
				} else {
					return runner.EmitIssue(
						r,
						fmt.Sprintf(
							"key-attribute `%s` should be defined before non-key `%s`",
							k.Name,
							n.Name(),
						),
						k.SrcRange,
					)
				}
			}
			pos++
		}

		return nil
	})
}
