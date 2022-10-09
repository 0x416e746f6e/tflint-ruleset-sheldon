package node

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Node is a drop-in replacement for HCL's nodes.
type Node interface {
	Range() hcl.Range
}

// InspectableNode provides additional information about the node.
type InspectableNode interface {
	Node

	AsAttribute() *hclsyntax.Attribute
	AsBlock() *hclsyntax.Block

	IsAttribute() bool
	IsBlock() bool
	Lines() int
	Kind() Kind
	Name() string
	Type() string
}

// Kind indicates the kind of a node (e.g. attribute, or block, or expression,
// and so on).
type Kind int

const (
	// Attribute ...
	Attribute Kind = iota

	// Block ...
	Block

	// Expression ...
	Expression

	// Composite node combines multiple nodes together.
	// For example, key + value expressions.
	Composite
)

func FirstNodeFrom(b *hclsyntax.Body) InspectableNode {
	if len(b.Attributes) == 0 && len(b.Blocks) == 0 {
		return nil
	}

	var fa *hclsyntax.Attribute
	for _, a := range b.Attributes {
		if fa == nil {
			fa = a
			continue
		}
		if fa.SrcRange.Start.Byte < a.SrcRange.Start.Byte {
			continue
		}
		fa = a
	}

	var fb *hclsyntax.Block
	if len(b.Blocks) > 0 {
		fb = b.Blocks[0]
	}
	for i := 1; i < len(b.Blocks); i++ {
		if fb == nil {
			fb = b.Blocks[i]
			continue
		}
		if fb.TypeRange.Start.Byte < b.Blocks[i].TypeRange.Start.Byte {
			continue
		}
		fb = b.Blocks[i]
	}

	if fa == nil {
		return WrapBlock(fb)
	}
	if fb == nil {
		return WrapAttribute(fa)
	}

	if fa.SrcRange.Start.Byte < fb.TypeRange.Start.Byte {
		return WrapAttribute(fa)
	}
	return WrapBlock(fb)
}

func LastNodeFrom(b *hclsyntax.Body) InspectableNode {
	if len(b.Attributes) == 0 && len(b.Blocks) == 0 {
		return nil
	}

	var fa *hclsyntax.Attribute
	for _, a := range b.Attributes {
		if fa == nil {
			fa = a
			continue
		}
		if fa.SrcRange.Start.Byte > a.SrcRange.Start.Byte {
			continue
		}
		fa = a
	}

	var fb *hclsyntax.Block
	if len(b.Blocks) > 0 {
		fb = b.Blocks[0]
	}
	for i := 1; i < len(b.Blocks); i++ {
		if fb == nil {
			fb = b.Blocks[i]
			continue
		}
		if fb.TypeRange.Start.Byte > b.Blocks[i].TypeRange.Start.Byte {
			continue
		}
		fb = b.Blocks[i]
	}

	if fa == nil {
		return WrapBlock(fb)
	}
	if fb == nil {
		return WrapAttribute(fa)
	}

	if fa.SrcRange.Start.Byte > fb.TypeRange.Start.Byte {
		return WrapAttribute(fa)
	}
	return WrapBlock(fb)
}

func OrderedNodesFrom(b *hclsyntax.Body) []Node {
	res := make(
		[]Node,
		0,
		len(b.Blocks)+len(b.Attributes),
	)

	for _, a := range b.Attributes {
		res = append(res, a)
	}
	for _, b := range b.Blocks {
		res = append(res, b)
	}

	sort.SliceStable(res, func(l, r int) bool {
		ls := res[l].Range().Start.Byte
		rs := res[r].Range().Start.Byte
		return ls < rs
	})

	return res
}

func OrderedInspecableNodesFrom(b *hclsyntax.Body) []InspectableNode {
	res := make(
		[]InspectableNode,
		0,
		len(b.Blocks)+len(b.Attributes),
	)

	for _, a := range b.Attributes {
		res = append(res, WrapAttribute(a))
	}
	for _, b := range b.Blocks {
		res = append(res, WrapBlock(b))
	}

	sort.SliceStable(res, func(l, r int) bool {
		ls := res[l].Range().Start.Byte
		rs := res[r].Range().Start.Byte
		return ls < rs
	})

	return res
}
