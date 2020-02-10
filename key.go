package di

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
)

// key is a id of provider in container
type key struct {
	name string
	res  reflect.Type
	typ  providerType
}

// String represent resultKey as string.
func (k key) String() string {
	if k.name == "" {
		return fmt.Sprintf("%s", k.res)
	}
	return fmt.Sprintf("%s[%s]", k.res, k.name)
}

// IsAlwaysVisible
func (k key) IsAlwaysVisible() bool {
	return k.typ == ptConstructor
}

// Package
func (k key) SubGraph() string {
	var pkg string
	switch k.res.Kind() {
	case reflect.Slice, reflect.Ptr:
		pkg = k.res.Elem().PkgPath()
	default:
		pkg = k.res.PkgPath()
	}

	return pkg
}

// Visualize
func (k key) Visualize(node *dot.Node) {
	node.Label(k.String())
	node.Attr("fontname", "COURIER")
	node.Attr("style", "filled")
	node.Attr("fontcolor", "white")
	switch k.typ {
	case ptConstructor:
		node.Attr("shape", "box")
		node.Attr("color", "#46494C")
	case ptGroup:
		node.Attr("shape", "doubleoctagon")
		node.Attr("color", "#E54B4B")
	case ptInterface:
		node.Attr("color", "#2589BD")
	case ptEmbedParameter:
		node.Attr("shape", "box")
		node.Attr("color", "#E5984B")
	}
}
