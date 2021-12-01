package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast  bool
	segment string
	handler ControllerHandler
	childs  []*node
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}

//判断一个segment是否是通用的segment,即 :开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	//如果segment是通配符，则下一层所有子节点都满足
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	//过滤下一层子节点
	for _, cnode := range n.childs {
		//如果下一层子节点有通配符，则满足所有需求
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			//如果下一层节点没有通配符，但是文本完全匹配，则满足所有需求
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	//使用分隔符将uri切割为2个部分
	segments := strings.SplitN(uri, "/", 2)
	//第一个部分用于匹配下一层子节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	//匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)

	// 如果当前子节点没有一个符合，那么说明这个uri一定是之前不存在, 直接返回nil
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	//如果只有一个segment ，则是最后一个标记
	if len(segments) == 1 {
		// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}

		// 都不是最后一个节点
		return nil
	}

	// 如果有2个segment, 递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// 增加路由节点
func (tree *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := tree.root
	// 确认路由是否冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}
	segments := strings.Split(uri, "/")
	//判断每一个segment是否存在
	for index, segment := range segments {
		//最终进入Node segment的字段
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}

		isLast := index == len(segments)-1

		//标记是否有合适的子节点
		var objNode *node

		childNodes := n.filterChildNodes(segment)
		//如果有匹配的子节点
		if len(childNodes) > 0 {
			//如果有segment相同的子节点 则选择这个子节点
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			//创建当前节点
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}

		n = objNode
	}

	return nil

}

// 匹配uri
func (tree *Tree) FindHandler(uri string) ControllerHandler {
	// 直接复用matchNode函数，uri是不带通配符的地址
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}

// // 匹配uri
// func (tree *Tree) PrintTree(root *node, level int) {

// 	if root == nil {
// 		fmt.Println("is null ?", root)
// 		return
// 	}
// 	for i := 0; i < len(root.childs); i++ {
// 		fmt.Printf("当前第%d层 segment= %s \n", level, root.segment)
// 		tree.PrintTree(root.childs[i], level+1)
// 	}
// }
