package di

import (
	"fmt"
	"sort"
	"strings"
)

// Tags are embedded expanders for injecting type. Embed it in your injecting type:
//
//  type ListConsoleCommand struct {
//    di.Tags `console.command:"list"`
//  }
//
// And use Resolve() with WithTag() option for fetch all instances
// have concrete tag-value combination:
//
//  var command Command
//  container.Resolve(&command, WithTag("console.command", "list"))
type Tags interface{}

func tagsToString(tags map[string]string) string {
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var kvs []string
	for _, key := range keys {
		kvs = append(kvs, fmt.Sprintf("%s:%s", key, tags[key]))
	}
	return strings.Join(kvs, ";")
}
