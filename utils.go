package goose

import (
	"fmt"
	"log"
)

// DrawRoutingTree is used for visualizing the routing tree
func (g *Goose) DrawRoutingTree(method string) {
	log.Println("WARNING: your routing system crashes when you use this function!!!")
	log.Println("WARNING: your routing system crashes when you use this function!!!")
	log.Println("WARNING: your routing system crashes when you use this function!!!")

	dummyHead := g.router.routers[method]
	root := dummyHead.children[0]

	// actually we need a deep copy of variable "root" here to avoid modifying the routing tree itself
	// but I tried a lot of answers in StackOverFlow and blogs, none of them works
	// and I don't want to introduce a third-party library just to fix this minor "bug"
	// since you should never use this function in production, it should only be used for
	// visualizing your registered routers during development

	queue := []*node{root}
	results := [][]string{}

	for len(queue) > 0 {
		currRes := []string{}
		nextLvl := []*node{}

		for _, node := range queue {
			currRes = append(currRes, node.segment)
			if len(node.children) > 0 {
				for _, child := range node.children {
					child.segment += fmt.Sprintf("(->%s)", node.segment)
				}
				nextLvl = append(nextLvl, node.children...)
			}
		}
		results = append(results, currRes)
		queue = nextLvl
	}

	for i, segments := range results {
		fmt.Printf("Level %d | ", i)
		for _, segment := range segments {
			fmt.Printf("  %s", segment)
		}
		fmt.Println()
	}
}
