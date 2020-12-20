package tree

import (
	"net/http"

	"github.com/hupf3/RestAPI/help/handlers"
	"github.com/hupf3/RestAPI/help/params"
	"github.com/hupf3/RestAPI/help/syntax"
)

// 树的形式存储路由
type Tree struct {
	node
	disableOptions bool
	disableHead    bool
}

// 声明一个新的树
func New(disableOptions, disableHead bool) *Tree {
	return &Tree{
		node:           node{segment: syntax.NewSegment("")},
		disableOptions: disableOptions,
		disableHead:    disableHead,
	}
}

// 添加路由项
func (tree *Tree) Add(pattern string, h http.Handler, methods ...string) error {
	n, err := tree.getNode(pattern)
	if err != nil {
		return err
	}

	if n.handlers == nil {
		n.handlers = handlers.New(tree.disableOptions, tree.disableHead)
	}

	return n.handlers.Add(h, methods...)
}

// 移除路由项
func (tree *Tree) Remove(pattern string, methods ...string) {
	child := tree.find(pattern)
	if child == nil || child.handlers == nil {
		return
	}

	if child.handlers.Remove(methods...) && len(child.children) == 0 {
		child.parent.children = removeNodes(child.parent.children, child.segment.Value)
		child.parent.buildIndexes()
	}
}

// 获取指定的节点
func (tree *Tree) getNode(pattern string) (*node, error) {
	segs, err := syntax.Split(pattern)
	if err != nil {
		return nil, err
	}

	return tree.node.getNode(segs), nil
}

// 设置指定节点的报头
func (tree *Tree) SetAllow(pattern, allow string) error {
	n, err := tree.getNode(pattern)
	if err != nil {
		return err
	}

	if n.handlers == nil {
		n.handlers = handlers.New(tree.disableOptions, tree.disableHead)
	}

	n.handlers.SetAllow(allow)
	return nil
}

// 生成URL地址
func (tree *Tree) URL(pattern string, params map[string]string) (string, error) {
	node, err := tree.getNode(pattern)
	if err != nil {
		return "", err
	}

	return node.url(params)
}

// 找到匹配的handle
func (tree *Tree) Handler(path string) (*handlers.Handlers, params.Params) {
	ps := make(params.Params, 3)
	node := tree.match(path, ps)

	if node == nil || node.handlers == nil || node.handlers.Len() == 0 {
		return nil, nil
	}
	return node.handlers, ps
}

// 获取所有的路由项
func (tree *Tree) All(ignoreHead, ignoreOptions bool) map[string][]string {
	routes := make(map[string][]string, 100)
	tree.all(ignoreHead, ignoreOptions, "", routes)
	return routes
}

// 清除路由项
func (tree *Tree) Clean(prefix string) {
	tree.clean(prefix)
}
