package gee

import (
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool // 表示该节点部分的路由是否允许模糊匹配
}

// 匹配子树中第一个节点
func (node *node) matchChild(part string) *node {
	for _, child := range node.children {
		// 匹配逻辑
		if part == child.part || child.isWild {
			return child
		}
	}
	// 未匹配到该节点的任何一子节点，匹配失败
	return nil
}

// 匹配子树中所有节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		// 匹配逻辑
		if part == child.part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 根据路由构建前缀树节点 前缀树的插入
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件
	// 同一层只会取最后一个模糊匹配节点，即覆盖该层先前注册的所有模糊匹配节点
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	child := n.matchChild(parts[height])
	if child == nil {
		child = &node{part: parts[height], isWild: parts[height][0] == ':' || parts[height][0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 根据路由匹配结果 前缀树的查找
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}

	return nil
}
