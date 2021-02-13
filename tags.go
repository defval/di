package di

import (
	"reflect"
	"sort"
	"strings"
)

var iTaggable = reflect.TypeOf(new(taggable)).Elem()

// Tags is a string representation of key value pairs.
//
// 	type Server struct {
//		di.Tags `http:"true" server:"true"`
//	}
//	_, err := di.New(
//		di.Provide(func() *Server { return &Server{} }),
//	)
//  var s *Server
//  c.Resolve(&s, di.Tags{"http": "true", "server": "true"})
//	require.NoError(t, err)
type Tags map[string]string

// injectable interface needs to struct fields injection functional.
type taggable interface {
	isTaggable()
}

func (t Tags) isTaggable() {
	bug()
}

// haveTags checks that typ is taggable
func haveTags(typ reflect.Type) bool {
	return typ.Implements(iTaggable)
}

func (t Tags) applyProvide(params *ProvideParams) {
	if params.Tags == nil {
		params.Tags = map[string]string{}
	}

	for k, v := range t {
		params.Tags[k] = v
	}
}

func (t Tags) applyResolve(params *ResolveParams) {
	if params.Tags == nil {
		params.Tags = map[string]string{}
	}
	for k, v := range t {
		params.Tags[k] = v
	}
}

// String is a tags string representation.
func (t Tags) String() string {
	var keys []string
	for k, _ := range t {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i := 0; i < len(keys); i++ {
		keys[i] = keys[i] + ":" + t[keys[i]]
	}
	if len(keys) == 0 {
		return ""
	}
	return "[" + strings.Join(keys, ";") + "]"
}

// matchTags checks that all of key value pairs exists in t.
func (t Tags) match(tags Tags) bool {
	for k, v := range tags {
		tv, ok := t[k]
		if !ok {
			return false
		}
		if v == "*" {
			continue
		}
		if tv != v {
			return false
		}
	}
	return true
}

// isEmpty returns true if tags empty.
func (t Tags) isEmpty() bool {
	return len(t) == 0
}

func matchTags(nodes []*node, tags Tags) []*node {
	matched := make([]*node, 0, 1)
	for i := 0; i < len(nodes); i++ {
		if nodes[i].tags.match(tags) {
			matched = append(matched, nodes[i])
		}
	}
	return matched
}
