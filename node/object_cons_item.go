package node

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type objectConsItem struct {
	item *hclsyntax.ObjectConsItem
}

func WrapObjectConsItem(i *hclsyntax.ObjectConsItem) InspectableNode {
	return &objectConsItem{
		item: i,
	}
}

func (i objectConsItem) Range() hcl.Range {
	return hcl.RangeBetween(
		i.item.KeyExpr.Range(),
		i.item.ValueExpr.Range(),
	)
}

func (i objectConsItem) Kind() Kind {
	return Composite
}

func (i objectConsItem) AsAttribute() *hclsyntax.Attribute {
	return nil
}

func (i objectConsItem) AsBlock() *hclsyntax.Block {
	return nil
}

func (i objectConsItem) IsAttribute() bool {
	return false
}

func (i objectConsItem) IsBlock() bool {
	return false
}

func (i objectConsItem) Type() string {
	return "key"
}

func (i objectConsItem) Name() string {
	var names []string

	switch x := i.item.KeyExpr.(type) {
	case *hclsyntax.ObjectConsKeyExpr:
		switch w := x.Wrapped.(type) {
		case *hclsyntax.ScopeTraversalExpr:
			for _, t := range w.Traversal {
				names = make([]string, 0, len(w.Traversal))
				switch k := t.(type) {
				// TODO: Figure out other possible types
				case hcl.TraverseRoot:
					names = append(names, k.Name)
				case hcl.TraverseAttr:
					names = append(names, k.Name)
				case *hcl.TraverseRoot:
					names = append(names, k.Name)
				case *hcl.TraverseAttr:
					names = append(names, k.Name)
				}
			}
		case *hclsyntax.TemplateExpr:
			names = make([]string, 0, len(w.Parts))
			for _, _p := range w.Parts {
				switch p := _p.(type) {
				case *hclsyntax.LiteralValueExpr:
					names = append(names, p.Val.AsString())
				}
			}
		}
	}

	pos := 0
	for _, n := range names {
		if len(n) > 0 {
			names[pos] = n
			pos++
		}
	}
	names = names[:pos]

	return strings.Join(names, ".")
}

func (i objectConsItem) Lines() int {
	r := i.Range()
	return r.End.Line - r.Start.Line + 1
}
