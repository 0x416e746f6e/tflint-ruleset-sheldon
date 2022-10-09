package node

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type block struct {
	block *hclsyntax.Block
}

func WrapBlock(b *hclsyntax.Block) InspectableNode {
	return &block{
		block: b,
	}
}

func (b block) Range() hcl.Range {
	return b.block.Range()
}

func (b block) Kind() Kind {
	return Block
}

func (b block) AsAttribute() *hclsyntax.Attribute {
	return nil
}

func (b block) AsBlock() *hclsyntax.Block {
	return b.block
}

func (b block) IsAttribute() bool {
	return false
}

func (b block) IsBlock() bool {
	return true
}

func (b block) Name() string {
	names := []string{}
	if b.block.Type != "dynamic" {
		names = append(names, b.block.Type)
	}
	names = append(names, b.block.Labels...)
	return strings.Join(names, " ")
}

func (b block) Type() string {
	return "block"
}

func (b block) Lines() int {
	r := b.Range()
	return r.End.Line - r.Start.Line + 1
}
