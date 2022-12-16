package rules

import (
	"fmt"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/node"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/project"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/visit"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// SpacingRule makes sure that there is consistent spacing between attributes
// blocks and.
type SpacingRule struct {
	tflint.DefaultRule
}

// NewSpacingRule creates a new SpacingRule.
func NewSpacingRule() *SpacingRule {
	return &SpacingRule{}
}

// Name returns the name of the rule.
func (r *SpacingRule) Name() string {
	return project.RuleName("spacing")
}

// Enabled returns whether the rule is enabled by default.
func (r *SpacingRule) Enabled() bool {
	return true
}

// Severity returns the severity of the rule.
func (r *SpacingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the reference link for the rule.
func (r *SpacingRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check verifies whether all attributes and blocks are properly sorted.
func (r *SpacingRule) Check(runner tflint.Runner) error {
	return visit.Files(r, runner, func(b *hclsyntax.Body, src []byte) error {
		return r.checkNodes(runner, src, 0, node.OrderedInspectableNodesFrom(b), b)
	})
}

func (r *SpacingRule) checkNodes(
	runner tflint.Runner,
	src []byte,
	level int,
	nodes []node.InspectableNode,
	parent node.Node,
) error {
	if len(nodes) == 0 {
		return r.checkEmptySpace(runner, src, parent)
	}

	if level == 0 {
		if err := r.checkLeadingSpace(runner, src, nodes[0], parent); err != nil {
			return nil
		}
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(nodes); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			nodes[i],
			nodes[i-1],
			isolateMultiliners,
		); err != nil {
			return err
		}
	}

	// Step inside
	for _, n := range nodes {
		if b := n.AsBlock(); b != nil {
			if err := r.checkBlock(runner, src, level+1, b); err != nil {
				return err
			}
		} else if a := n.AsAttribute(); a != nil {
			if err := r.checkExpression(runner, src, level+1, a.Expr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *SpacingRule) checkBlock(
	runner tflint.Runner,
	src []byte,
	level int,
	b *hclsyntax.Block,
) error {
	nodes := node.OrderedInspectableNodesFrom(b.Body)

	// Check the leading empty lines
	if len(nodes) > 0 {
		if err := r.checkLeadingSpace(runner, src, nodes[0], b.Body); err != nil {
			return err
		}
	}

	// Check special treatment of `count`, `for_each`, `depends_on`, and `lifecycle`
	if level == 1 && len(nodes) > 1 && (b.Type == "resource" || b.Type == "data") {
		if a := nodes[0].AsAttribute(); a != nil && (a.Name == "count" || a.Name == "for_each") {
			n := nodes[1]
			lines := n.Range().Start.Line - a.Range().End.Line
			if lines < 2 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"attribute `%s` must be separated from the rest of the definition by an extra line",
						a.Name,
					),
					a.Range(),
				); err != nil {
					return err
				}
			}
		}

		// TODO: depends_on, lifecycle
	}

	// Check special treatment of `source`
	if level == 1 && len(nodes) > 1 && b.Type == "module" {
		if a := nodes[0].AsAttribute(); a != nil && a.Name == "source" {
			n := nodes[1]
			lines := n.Range().Start.Line - a.Range().End.Line
			if lines < 2 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"attribute `%s` must be separated from the rest of the definition by an extra line",
						a.Name,
					),
					a.Range(),
				); err != nil {
					return err
				}
			}
		}
	}

	// Check inside the nested nodes
	if err := r.checkNodes(runner, src, level, nodes, b.Body); err != nil {
		return err
	}

	// Check the trailing empty lines
	if len(nodes) > 0 {
		if err := r.checkTrailingSpace(runner, src, nodes[len(nodes)-1], b.Body); err != nil {
			return err
		}
	}

	return nil
}

func (r *SpacingRule) checkExpression(
	runner tflint.Runner,
	src []byte,
	level int,
	expression hclsyntax.Expression,
) error {
	switch x := expression.(type) {
	case *hclsyntax.ForExpr:
		return r.checkForExpr(runner, src, level, x)
	case *hclsyntax.FunctionCallExpr:
		return r.checkFunctionCallExpr(runner, src, level, x)
	case *hclsyntax.ObjectConsExpr:
		return r.checkObjectConsExpr(runner, src, level, x)
	case *hclsyntax.ParenthesesExpr:
		return r.checkParenthesesExpr(runner, src, level, x)
	case *hclsyntax.TupleConsExpr:
		return r.checkTupleConsExpr(runner, src, level, x)
	default:
		return nil
	}
}

func (r *SpacingRule) checkForExpr(
	runner tflint.Runner,
	src []byte,
	level int,
	x *hclsyntax.ForExpr,
) error {
	// Step inside
	return r.checkExpression(runner, src, level+1, x.ValExpr)
}

func (r *SpacingRule) checkFunctionCallExpr(
	runner tflint.Runner,
	src []byte,
	level int,
	x *hclsyntax.FunctionCallExpr,
) error {
	// Check empty arguments list
	if len(x.Args) == 0 {
		return r.checkEmptySpace(runner, src, x)
	}

	if err := r.checkLeadingSpace(runner, src, x.Args[0], x); err != nil {
		return err
	}

	if err := r.checkTrailingSpace(runner, src, x.Args[len(x.Args)-1], x); err != nil {
		return err
	}

	for i := 1; i < len(x.Args); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			x.Args[i],
			x.Args[i-1],
			dontIsolateMultiliners,
		); err != nil {
			return err
		}
	}

	// Step inside
	for _, e := range x.Args {
		if err := r.checkExpression(runner, src, level+1, e); err != nil {
			return err
		}
	}

	return nil
}

func (r *SpacingRule) checkObjectConsExpr(
	runner tflint.Runner,
	src []byte,
	level int,
	x *hclsyntax.ObjectConsExpr,
) error {
	// Check empty object
	if len(x.Items) == 0 {
		return r.checkEmptySpace(runner, src, x)
	}

	// Check the leading empty lines
	if err := r.checkLeadingSpace(runner, src, x.Items[0].KeyExpr, x); err != nil {
		return err
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(x.Items); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			node.WrapObjectConsItem(&x.Items[i]),
			node.WrapObjectConsItem(&x.Items[i-1]),
			isolateMultiliners,
		); err != nil {
			return err
		}
	}

	// Check the trailing empty lines
	if err := r.checkTrailingSpace(runner, src, x.Items[len(x.Items)-1].ValueExpr, x); err != nil {
		return err
	}

	// Step inside
	for _, e := range x.Items {
		if err := r.checkExpression(
			runner,
			src,
			level+1,
			e.ValueExpr,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *SpacingRule) checkParenthesesExpr(
	runner tflint.Runner,
	src []byte,
	level int,
	x *hclsyntax.ParenthesesExpr,
) error {
	if err := r.checkLeadingSpace(runner, src, x.Expression, x); err != nil {
		return err
	}
	if err := r.checkTrailingSpace(runner, src, x.Expression, x); err != nil {
		return err
	}
	return nil
}

func (r *SpacingRule) checkTupleConsExpr(
	runner tflint.Runner,
	src []byte,
	level int,
	x *hclsyntax.TupleConsExpr,
) error {
	// Check empty list
	if len(x.Exprs) == 0 {
		return r.checkEmptySpace(runner, src, x)
	}

	// Check the leading empty lines
	if err := r.checkLeadingSpace(runner, src, x.Exprs[0], x); err != nil {
		return err
	}

	// Check the spacing in-between the elements
	for i := 1; i < len(x.Exprs); i++ {
		if err := r.checkMiddleSpace(
			runner,
			src,
			x.Exprs[i],
			x.Exprs[i-1],
			dontIsolateMultiliners,
		); err != nil {
			return err
		}
	}

	// Check the trailing empty lines
	if err := r.checkTrailingSpace(runner, src, x.Exprs[len(x.Exprs)-1], x); err != nil {
		return err
	}

	// Step inside
	for _, e := range x.Exprs {
		if err := r.checkExpression(runner, src, level+1, e); err != nil {
			return err
		}
	}

	return nil
}

// checkEmptySpace inspects the empty space within the range of the node.
func (r *SpacingRule) checkEmptySpace(
	runner tflint.Runner,
	src []byte,
	node node.Node,
) error {
	lines := node.Range().End.Line - node.Range().Start.Line

	if lines > 1 {
		_src := src[node.Range().Start.Byte:node.Range().End.Byte]
		_rng := hcl.Range{
			Filename: node.Range().Filename,
			Start:    node.Range().Start,
			End:      node.Range().End,
		}

		if err := r.checkLines(runner, _src, _rng, func(lines int, _ hcl.Range, _ bool) error {
			if lines > 1 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"%d redundant blank %s between the braces",
						lines-1,
						r.lineOrLines(lines-1),
					),
					node.Range(),
				); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// checkLeadingSpace inspects the empty space in front of the first node inside
// the parent's range.
func (r *SpacingRule) checkLeadingSpace(
	runner tflint.Runner,
	src []byte,
	node node.Node,
	parent node.Node,
) error {
	lines := node.Range().Start.Line - parent.Range().Start.Line
	if lines == 0 {
		return nil
	}

	_src := src[parent.Range().Start.Byte:node.Range().Start.Byte]
	_rng := hcl.Range{
		Filename: parent.Range().Filename,
		Start:    parent.Range().Start,
		End:      node.Range().Start,
	}
	if err := r.checkLines(runner, _src, _rng, func(lines int, _ hcl.Range, sawComments bool) error {
		threshold := 0
		if sawComments {
			threshold = 1
		}
		if lines > threshold {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"%d redundant blank %s in front",
					lines-threshold,
					r.lineOrLines(lines-threshold),
				),
				node.Range(),
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

type optIsolateMultiliners int

const (
	isolateMultiliners optIsolateMultiliners = iota
	dontIsolateMultiliners
)

// checkMiddleSpace inspects the contents in-between the nodes.
func (r *SpacingRule) checkMiddleSpace(
	runner tflint.Runner,
	src []byte,
	node node.Node,
	prev node.Node,
	opt optIsolateMultiliners,
) error {
	lines := node.Range().Start.Line - prev.Range().End.Line

	if lines > 1 {
		_src := src[prev.Range().End.Byte:node.Range().Start.Byte]
		_rng := hcl.Range{
			Filename: prev.Range().Filename,
			Start:    prev.Range().End,
			End:      node.Range().Start,
		}
		if err := r.checkLines(runner, _src, _rng, func(lines int, _ hcl.Range, _ bool) error {
			if lines > 1 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"%d redundant blank %s in front",
						lines-1,
						r.lineOrLines(lines-1),
					),
					node.Range(),
				); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if opt == isolateMultiliners {
		nl := node.Range().End.Line - node.Range().Start.Line
		pl := prev.Range().End.Line - prev.Range().Start.Line

		if nl > 0 && lines < 2 {
			if err := runner.EmitIssue(
				r,
				"multi-line element must be separated from the previous one by an extra line",
				node.Range(),
			); err != nil {
				return err
			}
		}

		if pl > 0 && nl == 0 && lines < 2 {
			if err := runner.EmitIssue(
				r,
				"single-line element must be separated from the preceding multi-line one by an extra line",
				node.Range(),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkTrailingSpace inspects the empty space after the last node inside the
// parent's range.
func (r *SpacingRule) checkTrailingSpace(
	runner tflint.Runner,
	src []byte,
	node node.Node,
	parent node.Node,
) error {
	lines := parent.Range().End.Line - node.Range().End.Line
	if lines <= 1 {
		return nil
	}

	_src := src[node.Range().End.Byte:parent.Range().End.Byte]
	_rng := hcl.Range{
		Filename: node.Range().Filename,
		Start:    node.Range().End,
		End:      parent.Range().End,
	}
	if err := r.checkLines(runner, _src, _rng, func(lines int, lastToken hcl.Range, sawComments bool) error {
		threshold := 0
		if sawComments {
			threshold = 1
		}
		if lines > threshold {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf(
					"%d redundant blank %s in front",
					lines-threshold,
					r.lineOrLines(lines-threshold),
				),
				lastToken,
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// checkLines inspects the content within a range from the blank-line
// perspective (taking into consideration the comments).
func (r *SpacingRule) checkLines(
	runner tflint.Runner,
	src []byte,
	rng hcl.Range,
	check func(lines int, lastToken hcl.Range, sawComments bool) error,
) error {
	tokens, err := hclsyntax.LexConfig(src, rng.Filename, rng.Start)
	if err != nil {
		return err
	}

	newLines := 0
	sawComments := false
	var lastToken hcl.Range
	for _, t := range tokens {
		switch t.Type {
		case hclsyntax.TokenNewline:
			newLines++
		case hclsyntax.TokenComment:
			lastToken = t.Range
			sawComments = true
			if newLines > 2 {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf(
						"%d redundant blank %s in front",
						newLines-2,
						r.lineOrLines(newLines-2),
					),
					t.Range,
				); err != nil {
					return err
				}
			}
			newLines = 1 // Comment token implies `\n` at the end
		case hclsyntax.TokenEOF:
			// noop
		default:
			lastToken = t.Range
		}
	}

	if err := check(newLines-1, lastToken, sawComments); err != nil {
		return err
	}

	return nil
}

func (r *SpacingRule) lineOrLines(count int) string {
	if count == 1 {
		return "line"
	}
	return "lines"
}
