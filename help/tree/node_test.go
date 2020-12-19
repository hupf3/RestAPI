package tree

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/hupf3/RestAPI/help/handlers"
	"github.com/hupf3/RestAPI/help/syntax"
)

func TestNode_find(t *testing.T) {
	a := assert.New(t)
	node := &node{}

	addNode := func(p string, code int, methods ...string) {
		segs, err := syntax.Split(p)
		a.NotError(err).NotNil(segs)
		nn := node.getNode(segs)
		a.NotNil(nn)

		if nn.handlers == nil {
			nn.handlers = handlers.New(false, false)
		}

		a.NotError(nn.handlers.Add(buildHandler(code), methods...))
	}

	addNode("/", 1, http.MethodGet)
	addNode("/posts/{id}", 1, http.MethodGet)
	a.Equal(node.find("/").segment.Value, "/")
	a.Equal(node.find("/posts/{id}").segment.Value, "{id}")
}

func TestRemoveNodes(t *testing.T) {
	a := assert.New(t)
	newNode := func(str string) *node {
		return &node{segment: syntax.NewSegment(str)}
	}

	n1 := newNode("/1")
	n2 := newNode("/2")
	n3 := newNode("/2")
	n4 := newNode("/3")
	n5 := newNode("/4")

	nodes := []*node{n1, n2, n3, n4, n5}

	nodes = removeNodes(nodes, "")
	a.Equal(len(nodes), 5)

	nodes = removeNodes(nodes, "/4")
	a.Equal(len(nodes), 4)
}

func TestSplitNode(t *testing.T) {
	a := assert.New(t)
	newNode := func(str string) *node {
		return &node{segment: syntax.NewSegment(str)}
	}
	p := newNode("/blog")

	a.Panic(func() {
		nn := splitNode(p, 1)
		a.Nil(nn)
	})

	node := p.newChild(syntax.NewSegment("/posts/{id}/author"))
	a.NotNil(node)

	nn := splitNode(node, 7)
	a.NotNil(nn)
	a.Equal(len(nn.children), 1).
		Equal(nn.children[0].segment.Value, "{id}/author")
	a.Equal(nn.parent, p)
}
