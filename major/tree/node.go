package tree

import (
	"fmt"
	"sort"
	"strings"

	"github.com/issue9/errwrap"

	"github.com/hupf3/RestAPI/major/handlers"
	"github.com/hupf3/RestAPI/major/syntax"
	"github.com/hupf3/RestAPI/params"
)

// 建立索引表大小
const indexesSize = 5

// 路由中的节点
type node struct {
	parent   *node
	handlers *handlers.Handlers
	children []*node
	segment  *syntax.Segment

	indexes map[byte]int
}

// 构建索引表
func (n *node) buildIndexes() {
	if len(n.children) < indexesSize {
		n.indexes = nil
		return
	}

	if n.indexes == nil {
		n.indexes = make(map[byte]int, indexesSize)
	}

	for index, node := range n.children {
		if node.segment.Type == syntax.String {
			n.indexes[node.segment.Value[0]] = index
		}
	}
}

// 返回优先级
func (n *node) priority() int {
	ret := int(n.segment.Type) * 10

	if len(n.children) > 0 || n.segment.Endpoint {
		return ret + 1
	}

	return ret
}

// 获取指定路径下的节点
func (n *node) getNode(segments []*syntax.Segment) *node {
	child := n.addSegment(segments[0])

	if len(segments) == 1 {
		return child
	}

	return child.getNode(segments[1:])
}

// 在制定节点添加Segment
func (n *node) addSegment(seg *syntax.Segment) *node {
	var child *node
	var l int
	for _, c := range n.children {
		l1 := c.segment.Similarity(seg)

		if l1 == -1 {
			return c
		}

		if l1 > l {
			l = l1
			child = c
		}
	}

	if l <= 0 {
		return n.newChild(seg)
	}

	parent := splitNode(child, l)

	if len(seg.Value) == l {
		return parent
	}

	return parent.addSegment(syntax.NewSegment(seg.Value[l:]))
}

// 产生新的子节点
func (n *node) newChild(s *syntax.Segment) *node {
	child := &node{
		parent:  n,
		segment: s,
	}

	n.children = append(n.children, child)
	sort.SliceStable(n.children, func(i, j int) bool {
		return n.children[i].priority() < n.children[j].priority()
	})
	n.buildIndexes()

	return child
}

// 查找路由项
func (n *node) find(pattern string) *node {
	for _, child := range n.children {
		if child.segment.Value == pattern {
			return child
		}

		if strings.HasPrefix(pattern, child.segment.Value) {
			nn := child.find(pattern[len(child.segment.Value):])
			if nn != nil {
				return nn
			}
		}
	}

	return nil
}

// 清除路由项
func (n *node) clean(prefix string) {
	if len(prefix) == 0 {
		n.children = n.children[:0]
		return
	}

	dels := make([]string, 0, len(n.children))
	for _, child := range n.children {
		if len(child.segment.Value) < len(prefix) {
			if strings.HasPrefix(prefix, child.segment.Value) {
				child.clean(prefix[len(child.segment.Value):])
			}
		}

		if strings.HasPrefix(child.segment.Value, prefix) {
			dels = append(dels, child.segment.Value)
		}
	}

	for _, del := range dels {
		n.children = removeNodes(n.children, del)
	}
	n.buildIndexes()
}

// 查找当前路径匹配的节点
func (n *node) match(path string, params params.Params) *node {
	if len(n.indexes) > 0 && len(path) > 0 {
		node := n.children[n.indexes[path[0]]]
		if node == nil {
			goto LOOP
		}

		index := node.segment.Match(path, params)
		if index < 0 {
			goto LOOP
		}

		if nn := node.match(path[index:], params); nn != nil {
			return nn
		}
	}

LOOP:
	for i := len(n.indexes); i < len(n.children); i++ {
		node := n.children[i]

		index := node.segment.Match(path, params)
		if index < 0 {
			continue
		}

		if nn := node.match(path[index:], params); nn != nil {
			return nn
		}

		delete(params, n.segment.Name)
	}

	if len(path) == 0 && n.handlers != nil && n.handlers.Len() > 0 {
		return n
	}

	return nil
}

// 生成URL地址
func (n *node) url(params map[string]string) (string, error) {
	nodes := make([]*node, 0, 5)
	for curr := n; curr.parent != nil; curr = curr.parent {
		nodes = append(nodes, curr)
	}

	var buf errwrap.StringBuilder
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		switch node.segment.Type {
		case syntax.String:
			buf.WString(node.segment.Value)
		case syntax.Named, syntax.Regexp:
			param, exists := params[node.segment.Name]
			if !exists {
				return "", fmt.Errorf("未找到参数 %s 的值！", node.segment.Name)
			}
			buf.WString(param).
				WString(node.segment.Suffix)
		}
	}

	if buf.Err != nil {
		return "", buf.Err
	}
	return buf.String(), nil
}

// 获取所有的地址
func (n *node) all(ignoreHead, ignoreOptions bool, parent string, routes map[string][]string) {
	path := parent + n.segment.Value

	if n.handlers != nil && n.handlers.Len() > 0 {
		routes[path] = n.handlers.Methods(ignoreHead, ignoreOptions)
	}

	for _, v := range n.children {
		v.all(ignoreHead, ignoreOptions, path, routes)
	}
}

// 删除指定节点
func removeNodes(nodes []*node, pattern string) []*node {
	for index, n := range nodes {
		if n.segment.Value == pattern {
			return append(nodes[:index], nodes[index+1:]...)
		}
	}

	return nodes
}

// 拆分制定位置，并返回节点
func splitNode(n *node, pos int) *node {
	if len(n.segment.Value) <= pos {
		return n
	}

	p := n.parent
	if p == nil {
		panic("节点必须要有一个有效的父节点，才能进行拆分！")
	}

	p.children = removeNodes(p.children, n.segment.Value)
	p.buildIndexes()

	segs := n.segment.Split(pos)
	ret := p.newChild(segs[0])
	c := ret.newChild(segs[1])
	c.handlers = n.handlers
	c.children = n.children
	c.indexes = n.indexes
	for _, item := range c.children {
		item.parent = c
	}

	return ret
}
