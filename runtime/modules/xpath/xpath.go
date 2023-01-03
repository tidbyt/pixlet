package xpath

import (
	"fmt"
	"strings"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type XPath struct {
	doc           *xmlquery.Node
	query         *starlark.Builtin
	queryAll      *starlark.Builtin
	queryNode     *starlark.Builtin
	queryAllNodes *starlark.Builtin
}

var (
	xPathOnce   sync.Once
	xPathModule starlark.StringDict
)

func LoadXPathModule() (starlark.StringDict, error) {
	xPathOnce.Do(func() {
		xPathModule = starlark.StringDict{
			"xpath": &starlarkstruct.Module{
				Name: "xpath",
				Members: starlark.StringDict{
					"loads": starlark.NewBuiltin("loads", xPathLoads),
				},
			},
		}
	})

	return xPathModule, nil
}

func xPathLoads(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var xml starlark.String

	if err := starlark.UnpackArgs(
		"loads",
		args, kwargs,
		"xml", &xml,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for cache.get: %v", err)
	}

	doc, err := xmlquery.Parse(strings.NewReader(xml.GoString()))
	if err != nil {
		return nil, fmt.Errorf("parsing XML: %v", err)
	}

	x := &XPath{
		doc:           doc,
		query:         starlark.NewBuiltin("query", xPathQuery),
		queryAll:      starlark.NewBuiltin("query_all", xPathQueryAll),
		queryNode:     starlark.NewBuiltin("query_node", xPathQueryNode),
		queryAllNodes: starlark.NewBuiltin("query_all_node", xPathQueryAllNodes),
	}

	return x, nil
}

func xPathQueryNode(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path starlark.String

	if err := starlark.UnpackArgs(
		"query",
		args, kwargs,
		"path", &path,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for query: %v", err)
	}

	x := b.Receiver().(*XPath)

	node, err := xmlquery.Query(x.doc, path.GoString())
	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}

	if node == nil {
		return starlark.None, nil
	}

	result := &XPath{
		doc:           node,
		query:         starlark.NewBuiltin("query", xPathQuery),
		queryAll:      starlark.NewBuiltin("query_all", xPathQueryAll),
		queryNode:     starlark.NewBuiltin("query_node", xPathQueryNode),
		queryAllNodes: starlark.NewBuiltin("query_all_node", xPathQueryAllNodes),
	}
	return result, nil
}

func xPathQuery(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path starlark.String

	if err := starlark.UnpackArgs(
		"query",
		args, kwargs,
		"path", &path,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for query: %v", err)
	}

	x := b.Receiver().(*XPath)

	node, err := xmlquery.Query(x.doc, path.GoString())
	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}

	if node == nil {
		return starlark.None, nil
	}
	return starlark.String(node.InnerText()), nil
}

func xPathQueryAllNodes(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path starlark.String

	if err := starlark.UnpackArgs(
		"query_all",
		args, kwargs,
		"path", &path,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for query_all: %v", err)
	}

	x := b.Receiver().(*XPath)

	nodes, err := xmlquery.QueryAll(x.doc, path.GoString())
	if err != nil {
		return nil, fmt.Errorf("querying all: %v", err)
	}

	results := make([]starlark.Value, 0, len(nodes))
	for _, n := range nodes {
		result := &XPath{
			doc:           n,
			query:         starlark.NewBuiltin("query", xPathQuery),
			queryAll:      starlark.NewBuiltin("query_all", xPathQueryAll),
			queryNode:     starlark.NewBuiltin("query_node", xPathQueryNode),
			queryAllNodes: starlark.NewBuiltin("query_all_node", xPathQueryAllNodes),
		}
		results = append(results, result)
	}

	return starlark.NewList(results), nil
}

func xPathQueryAll(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path starlark.String

	if err := starlark.UnpackArgs(
		"query_all",
		args, kwargs,
		"path", &path,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for query_all: %v", err)
	}

	x := b.Receiver().(*XPath)

	nodes, err := xmlquery.QueryAll(x.doc, path.GoString())
	if err != nil {
		return nil, fmt.Errorf("querying all: %v", err)
	}

	nodeTexts := make([]starlark.Value, 0, len(nodes))
	for _, n := range nodes {
		nodeTexts = append(nodeTexts, starlark.String(n.InnerText()))
	}

	return starlark.NewList(nodeTexts), nil
}

func (x *XPath) AttrNames() []string {
	return []string{
		"query",
		"query_all",
		"query_node",
		"query_all_nodes",
	}
}

func (x *XPath) Attr(name string) (starlark.Value, error) {
	switch name {

	case "query":
		return x.query.BindReceiver(x), nil

	case "query_all":
		return x.queryAll.BindReceiver(x), nil

	case "query_node":
		return x.queryNode.BindReceiver(x), nil

	case "query_all_nodes":
		return x.queryAllNodes.BindReceiver(x), nil

	default:
		return nil, nil
	}
}

func (x *XPath) String() string       { return "XPath(...)" }
func (x *XPath) Type() string         { return "XPath" }
func (x *XPath) Freeze()              {}
func (x *XPath) Truth() starlark.Bool { return true }

func (x *XPath) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(x, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
