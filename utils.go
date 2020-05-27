package goose

import "fmt"

// DrawRoutingTree is used for visualizing the routing tree
func (g *Goose) DrawRoutingTree(method string) {
	dummyHead := g.router.routers[method]
	root := dummyHead.children[0]

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