package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type optSorting uint

const (
	dontSortMultiliners optSorting = 1 << iota
)

func (o optSorting) dontSortMultiliners() bool {
	return o&dontSortMultiliners != 0
}

// SortingRule makes sure that all attributes and blocks are properly sorted.
type SortingRule struct {
	tflint.DefaultRule
}

// NewSortingRule creates a new SortingRule.
func NewSortingRule() *SortingRule {
	return &SortingRule{}
}

// Name returns the name of the rule.
func (r *SortingRule) Name() string {
	return project.RuleName("sorting")
}

// Enabled returns whether the rule is enabled by default.
func (r *SortingRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *SortingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SortingRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether all attributes and blocks are properly sorted.
func (r *SortingRule) Check(rr tflint.Runner) error {
	runner, ok := rr.(*custom.Runner)
	if !ok {
		return fmt.Errorf("unexpected runner type: %s", reflect.TypeOf(rr))
	}

	return visit.Files(r, runner, func(b *hclsyntax.Body, src []byte) error {
		return r.checkNodes(runner, src, 0, node.OrderedInspecableNodesFrom(b))
	})
}

func (r *SortingRule) checkNodes(
	runner *custom.Runner,
	src []byte,
	level int,
	nodes []node.InspectableNode,
) error {
	for i := 1; i < len(nodes); i++ {
		n := nodes[i]
		p := nodes[i-1]

		// Lint: single-line node should precede multi-line
		if p.Lines() > 1 && n.Lines() == 1 {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"single-line node `%s` should precede multi-line `%s`",
					n.Name(),
					p.Name(),
				),
				n.Range(),
			); err != nil {
				return err
			}
		}

		// Lint: attribute should precede block
		if p.IsBlock() && n.IsAttribute() {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"attribute `%s` should precede block `%s`",
					n.Name(),
					p.Name(),
				),
				n.Range(),
			); err != nil {
				return err
			}
		}

		if p.Kind() != n.Kind() {
			continue
		}

		// Enforce sorting inside consecutive segments of single-line elements
		if n.Range().Start.Line-p.Range().End.Line == 1 &&
			p.Lines() == 1 &&
			n.Lines() == 1 &&
			n.Name() < p.Name() {
			// ---
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("%s `%s` should be placed after `%s` (alphabetical sorting)",
					p.Type(),
					p.Name(),
					n.Name(),
				),
				p.Range(),
			); err != nil {
				return err
			}
		}

		// Enforce sorting of multi-liners regardless of spacing between them
		// (but only if there are no comments between them)
		if level > 0 &&
			p.Lines() > 1 &&
			n.Lines() > 1 &&
			n.Name() < p.Name() {
			// ---
			hasComments, err := r.hasCommentsInBetween(src, p, n)
			if err != nil {
				return err
			}
			if !hasComments {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"%s `%s` should be placed after `%s` (alphabetical sorting)",
						p.Type(),
						p.Name(),
						n.Name(),
					),
					p.Range(),
				); err != nil {
					return err
				}
			}
		}
	}

	// Step inside
	for _, n := range nodes {
		if a := n.AsAttribute(); a != nil {
			if err := r.checkExpression(runner, level+1, a.Expr); err != nil {
				return err
			}
		} else if b := n.AsBlock(); b != nil {
			if err := r.checkBlock(runner, src, level+1, b); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *SortingRule) checkBlock(
	runner *custom.Runner,
	src []byte,
	level int,
	b *hclsyntax.Block,
) error {
	nodes := node.OrderedInspecableNodesFrom(b.Body)
	if len(nodes) == 0 {
		return nil
	}

	if level == 1 {
		// Top-level `locals` block: just check the expressions and exit
		if b.Type == "locals" {
			for _, _n := range nodes {
				if _a := _n.AsAttribute(); _a != nil {
					if err := r.checkExpression(runner, level+1, _a.Expr, dontSortMultiliners); err != nil {
						return err
					}
				}
			}
			return nil
		}

		// Top-level `resource` or `data` block
		if b.Type == "resource" || b.Type == "data" {
			var err error
			nodes, err = r.preprocessResourceOrData(runner, src, level, b.Labels[0], nodes)
			if err != nil {
				return err
			}
		}

		// Top-level `module` block
		if b.Type == "module" {
			var err error
			nodes, err = r.preprocessModule(runner, level, b.Labels[0], nodes)
			if err != nil {
				return err
			}
		}
	}

	if err := r.checkNodes(runner, src, level+1, nodes); err != nil {
		return err
	}

	return nil
}

func (r *SortingRule) preprocessResourceOrData(
	runner *custom.Runner,
	src []byte,
	level int,
	kind string,
	nodes []node.InspectableNode,
) ([]node.InspectableNode, error) {
	// Drop leading `for_each`
	if disabled, exists := runner.Disabled["for_each"]; !exists || !disabled {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "for_each" {
			nodes = nodes[1:]
		}
		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop leading `count`
	if disabled, exists := runner.Disabled["count"]; !exists || !disabled {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "count" {
			nodes = nodes[1:]
		}
		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop trailing `depends_on`
	if disabled, exists := runner.Disabled["depends_on"]; !exists || disabled {
		if a := nodes[len(nodes)-1].AsAttribute(); a != nil && a.Name == "depends_on" {
			nodes = nodes[:len(nodes)-1]
		}
		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop trailing `lifecycle`
	if disabled, exists := runner.Disabled["lifecycle"]; !exists || !disabled {
		if b := nodes[len(nodes)-1].AsBlock(); b != nil && b.Type == "lifecycle" {
			nodes = nodes[:len(nodes)-1]
		}
		if len(nodes) == 0 {
			return nodes, nil
		}
	}

	// Drop key-attributes
	if disabled, exists := runner.Disabled["key_attributes"]; !exists || !disabled {
		// Are there key-attributes?
		resource, ok := runner.Resources[kind]
		if !ok {
			return nodes, nil
		}
		if len(resource.KeyAttributes) == 0 {
			return nodes, nil
		}

		// TODO: Allow for nested key attrs (e.g. `kubernetes_manifest.manifest.metadata.*`)

		// Drop key-attributes
		if level > len(resource.KeyBlocks) {
			for _, ka := range resource.KeyAttributes {
				if len(nodes) == 0 {
					break
				}
				if a := nodes[0].AsAttribute(); a != nil {
					if a.Name == ka {
						nodes = nodes[1:] // Drop
					}
				}
			}
			return nodes, nil
		}

		// Drop the key block, but inspect its contents
		if kb := nodes[0].AsBlock(); kb != nil && kb.Type == resource.KeyBlocks[level-1] {
			nodes = nodes[1:] // Drop
			knodes := node.OrderedInspecableNodesFrom(kb.Body)
			// Drop leading key attributes inside the key block
			if len(resource.KeyBlocks) == level {
				for _, ka := range resource.KeyAttributes {
					if len(knodes) == 0 {
						break
					}
					if a := knodes[0].AsAttribute(); a != nil {
						if a.Name == ka {
							knodes = knodes[1:] // Drop
						}
					}
				}
			}
			// Step inside
			if err := r.checkNodes(runner, src, level+1, knodes); err != nil {
				return nil, err
			}
		}
	}

	return nodes, nil
}

func (r *SortingRule) preprocessModule(
	runner *custom.Runner,
	level int,
	kind string,
	nodes []node.InspectableNode,
) ([]node.InspectableNode, error) {
	// Drop leading `source`
	if disabled, exists := runner.Disabled["source"]; !exists || !disabled {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "source" {
			nodes = nodes[1:]
		}
	}
	return nodes, nil
}

func (r *SortingRule) checkExpression(
	runner *custom.Runner,
	level int,
	expression hclsyntax.Expression,
	opts ...optSorting,
) error {
	var opt optSorting
	for _, o := range opts {
		opt = opt | o
	}

	switch x := expression.(type) {
	case *hclsyntax.ForExpr:
		return r.checkForExpr(runner, level, x, opt)
	case *hclsyntax.FunctionCallExpr:
		return r.checkFunctionCallExpr(runner, level, x, opt)
	case *hclsyntax.ObjectConsExpr:
		return r.checkObjectConsExpr(runner, level, x, opt)
	case *hclsyntax.ParenthesesExpr:
		return r.checkParenthesesExpr(runner, level, x, opt)
	case *hclsyntax.TupleConsExpr:
		return r.checkTupleConsExpr(runner, level, x, opt)
	default:
		return nil
	}
}

func (r *SortingRule) checkForExpr(
	runner *custom.Runner,
	level int,
	x *hclsyntax.ForExpr,
	opt optSorting,
) error {
	// Step inside
	return r.checkExpression(runner, level+1, x.ValExpr, opt)
}

func (r *SortingRule) checkFunctionCallExpr(
	runner *custom.Runner,
	level int,
	x *hclsyntax.FunctionCallExpr,
	opt optSorting,
) error {
	// Step inside
	for _, e := range x.Args {
		if err := r.checkExpression(runner, level+1, e, opt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SortingRule) checkObjectConsExpr(
	runner *custom.Runner,
	level int,
	x *hclsyntax.ObjectConsExpr,
	opt optSorting,
) error {
	for i := 1; i < len(x.Items); i++ {
		e := x.Items[i]
		p := x.Items[i-1]

		el := e.ValueExpr.Range().End.Line - e.KeyExpr.Range().Start.Line
		pl := p.ValueExpr.Range().End.Line - p.KeyExpr.Range().Start.Line

		ek := strings.ToLower(node.WrapObjectConsItem(&e).Name())
		pk := strings.ToLower(node.WrapObjectConsItem(&p).Name())

		if pl > 0 && el == 0 {
			if err := runner.EmitIssue(
				r,
				"singleline key/value pair should be placed before multiline",
				e.KeyExpr.Range(),
			); err != nil {
				return err
			}
		}

		if pl == 0 && el == 0 &&
			e.KeyExpr.Range().Start.Line-p.ValueExpr.Range().End.Line == 1 &&
			ek < pk {
			// ---
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"key `%s` is out of order (should not follow alphabetically greater `%s`)",
					ek,
					pk,
				),
				e.KeyExpr.Range(),
			); err != nil {
				return err
			}
		}

		if pl > 0 && el > 0 && !opt.dontSortMultiliners() &&
			ek < pk {
			// ---
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"key `%s` is out of order (should not follow alphabetically greater `%s`)",
					ek,
					pk,
				),
				e.KeyExpr.Range(),
			); err != nil {
				return err
			}
		}
	}

	// Step inside
	for _, e := range x.Items {
		if err := r.checkExpression(runner, level+1, e.ValueExpr, opt); err != nil {
			return err
		}
	}

	return nil
}

func (r *SortingRule) checkTupleConsExpr(
	runner *custom.Runner,
	level int,
	x *hclsyntax.TupleConsExpr,
	opt optSorting,
) error {
	// Step inside
	for _, e := range x.Exprs {
		if err := r.checkExpression(runner, level+1, e, opt); err != nil {
			return err
		}
	}
	return nil
}

func (r *SortingRule) checkParenthesesExpr(
	runner *custom.Runner,
	level int,
	x *hclsyntax.ParenthesesExpr,
	opt optSorting,
) error {
	// Step inside
	return r.checkExpression(runner, level+1, x.Expression, opt)
}

func (r *SortingRule) hasCommentsInBetween(
	src []byte,
	p node.Node,
	n node.Node,
) (bool, error) {
	rng := hcl.Range{
		Filename: n.Range().Filename,
		Start: hcl.Pos{
			Line:   p.Range().End.Line,
			Column: p.Range().End.Column,
			Byte:   p.Range().End.Byte,
		},
		End: hcl.Pos{
			Line:   n.Range().Start.Line,
			Column: n.Range().Start.Column,
			Byte:   n.Range().Start.Byte,
		},
	}

	tokens, err := hclsyntax.LexConfig(
		src[rng.Start.Byte:rng.End.Byte],
		rng.Filename,
		rng.Start,
	)
	if err != nil {
		return false, err
	}

	for _, t := range tokens {
		switch t.Type {
		case hclsyntax.TokenComment:
			return true, nil
		}
	}

	return false, nil
}
