package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配，part 含有 : 或 * 时为true，即为模糊匹配
}

// String 用于调试
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// matchChild 用于在子节点中匹配指定 part 的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 用于在子节点中匹配所有满足条件的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 用于插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件 ：parts 长度与 height 相等
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil { // 不存在则创建
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 递归调用
	child.insert(pattern, parts, height+1)
}

// search 用于搜索节点
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件 ：parts 长度与 height 相等 或 part 为 * 时
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // * 匹配所有
		if n.pattern == "" {
			return nil
		}
		// 找到匹配项
		return n
	}
	part := parts[height]
	children := n.matchChildren(part) // 匹配所有满足条件的节点
	for _, child := range children {
		// 递归调用
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	// 未找到匹配项
	return nil
}

// 遍历所有节点找到存在pattern的节点
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
