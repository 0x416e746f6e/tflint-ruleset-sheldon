package node

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type attribute struct {
	attribute *hclsyntax.Attribute
}

func WrapAttribute(a *hclsyntax.Attribute) InspectableNode {
	return &attribute{
		attribute: a,
	}
}

func (a attribute) Range() hcl.Range {
	return a.attribute.Range()
}

func (a attribute) Kind() Kind {
	return Attribute
}

func (a attribute) AsAttribute() *hclsyntax.Attribute {
	return a.attribute
}

func (a attribute) AsBlock() *hclsyntax.Block {
	return nil
}

func (a attribute) IsAttribute() bool {
	return true
}

func (a attribute) IsBlock() bool {
	return false
}

func (a attribute) Name() string {
	return a.attribute.Name
}

func (a attribute) Type() string {
	return "attribute"
}

func (a attribute) Lines() int {
	r := a.Range()
	return r.End.Line - r.Start.Line + 1
}
