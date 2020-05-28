package goose

import (
	"fmt"
	"strings"
)

type node struct {
	segment    string
	isWildcard bool
	handler    HandlerFunc
	children   []*node
}

func (n *node) insert(pattern string, handler HandlerFunc) {
	validateRoutingPattern(pattern)

	segments := parsePattern(pattern)
	n.insertHelper(segments, 0, handler)
}

func (n *node) insertHelper(segments []string, level int, handler HandlerFunc) {
	segment := segments[level]
	isLastSegment := (level == len(segments)-1)
	isWildSegment := strings.HasPrefix(segment, ":") || strings.HasPrefix(segment, "*")

	/* if the incoming segment is a wildcard and is the last segment
	   it means that this segment acts as a wildcard endpoint

	   Note: an endpoint in the above context refers to "a node attached with non-nil handler"
	         also, the last node/segment must be an endpoint (this is the definition)
	*/
	if isWildSegment && isLastSegment {
		needNewChild := true
		existingChild := &node{}

		for _, child := range n.children {
			// make sure there's no other endpoint exists in n.children to avoid conflicts
			if child.handler != nil {
				panic(
					fmt.Sprintf(
						"Wildcard segment \"%s\" conflicts with existing segment \"%s\"",
						segment,
						child.segment,
					),
				)
			}

			if segment == child.segment {
				needNewChild = false
				existingChild = child
			}
		}

		// if we got here, it means no endpoint exists in n.children
		if needNewChild {
			newChild := &node{
				segment:    segment,
				isWildcard: isWildSegment,
				handler:    handler, // we only assign the handler when it's an endpoint
			}
			n.children = append(n.children, newChild)
			return
		} else {
			existingChild.handler = handler
			return
		}
	}

	if !isWildSegment && isLastSegment {
		needNewChild := true
		existingChild := &node{}

		for _, child := range n.children {
			if child.isWildcard && child.handler != nil {
				panic(
					fmt.Sprintf(
						"Segment \"%s\" conflicts with existing wildcard segment \"%s\"",
						segment,
						child.segment,
					),
				)
			} else if segment == child.segment && child.handler != nil {
				panic(
					fmt.Sprintf(
						"Segment \"%s\" conflicts with existing segment \"%s\"",
						segment,
						child.segment,
					),
				)
			}

			if segment == child.segment {
				needNewChild = false
				existingChild = child
			}
		}

		if needNewChild {
			newChild := &node{
				segment:    segment,
				isWildcard: isWildSegment,
				handler:    handler,
			}
			n.children = append(n.children, newChild)
			return
		} else {
			existingChild.handler = handler
			return
		}
	}

	if !isLastSegment {
		needNewChild := true
		existingChild := &node{}

		for _, child := range n.children {
			if segment == child.segment {
				needNewChild = false
				existingChild = child
				break
			}
		}

		if needNewChild {
			newChild := &node{
				segment:    segment,
				isWildcard: isWildSegment,
			}
			n.children = append(n.children, newChild)
			newChild.insertHelper(segments, level+1, handler)
		} else {
			existingChild.insertHelper(segments, level+1, handler)
		}

	}
}

// func (n *node) insertHelper(segments []string, level int, handler HandlerFunc) {
// 	segment := segments[level]
// 	child := n.matchChild(segment)

// 	// if we are processing the last segment (aka. endpoint),
// 	// assign the handler to the child node and exit
// 	if level == len(segments)-1 {
// 		// if the child node doesn't exist, create one, then assign the handler
// 		if child == nil {
// 			child = &node{
// 				segment:    segment,
// 				isWildcard: strings.HasPrefix(segment, ":") || strings.HasPrefix(segment, "*"),
// 				handler:    handler,
// 			}
// 			n.children = append(n.children, child)
// 			return
// 		}

// 		// if the child node already exists

// 		/* scenario 1:
// 		   if the handler of the child node has already been set, it means that
// 		   the child node is either a wildcard node or a repetitive routing. We should panic here.
// 		*/
// 		if child.handler != nil {
// 			panic(fmt.Sprintf(
// 				"Found conflicts between %s and existing %s",
// 				segment,
// 				child.segment,
// 			))
// 		} else {
// 			/* scenario 2:
// 			   if the handler of the child node has not been set, it means that previously
// 			   the child node is just a middle point of a URI, it's not an endpoint. That's why
// 			   it was not assigned with any handler.

// 			   Therefore, in this case, we assign the handler to it.
// 			   (from now on, the child node is no longer a "nobody cares" middle point. We assigned it
// 			    a handler, which makes the node an endpoint)
// 			*/
// 			child.handler = handler
// 		}
// 		return
// 	}

// 	// if we are not processing the last segment

// 	// if the segment doesn't match any child node, add a new child
// 	if child == nil {
// 		child = &node{
// 			segment:    segment,
// 			isWildcard: strings.HasPrefix(segment, ":") || strings.HasPrefix(segment, "*"),
// 		}
// 		n.children = append(n.children, child)
// 	}
// 	// move on to the next level
// 	child.insertHelper(segments, level+1, handler)
// }

func (n *node) search(pattern string) (*node, map[string]string) {
	segments := parsePattern(pattern)
	searchResultPtr := &node{}
	params := make(map[string]string)
	n.searchHelper(segments, 0, searchResultPtr, params)
	return searchResultPtr, params
}

func (n *node) searchHelper(segments []string, level int, searchResultPtr *node, params map[string]string) {
	segment := segments[level]
	isLastSegment := (level == len(segments)-1)

	if isLastSegment {
		for _, child := range n.children {
			if (child.segment == segment || child.isWildcard) && child.handler != nil {
				if child.isWildcard {
					params[child.segment] = segment
				}
				*searchResultPtr = *child
				return
			}
		}
	} else {
		skipWildcardMatching := false

		for _, child := range n.children {
			if child.segment == segment {
				skipWildcardMatching = true
				child.searchHelper(segments, level+1, searchResultPtr, params)
			}
		}

		if !skipWildcardMatching {
			for _, child := range n.children {
				if child.isWildcard {
					params[child.segment] = segment
					child.searchHelper(segments, level+1, searchResultPtr, params)
				}
			}
		}
	}
}

// func (n *node) searchHelper(segments []string, level int, searchResultPtr *node, params map[string]string) {
// 	segment := segments[level]
// 	child := n.matchChild(segment)

// 	// if there is no match for the incoming segment, exit
// 	if child == nil {
// 		return
// 	}

// 	// else
// 	// => save the "wildcard: real-value" mapping in map:params
// 	if child.isWildcard {
// 		params[child.segment] = segment
// 	}
// 	// => if the incoming segment is an asteriod wildcard, or we are dealing with the last segment
// 	//    save the search result and return
// 	if strings.HasPrefix(child.segment, "*") || (level == len(segments)-1) {
// 		*searchResultPtr = *child
// 		return
// 	}
// 	child.searchHelper(segments, level+1, searchResultPtr, params)
// }

// func (n *node) matchChild(segment string) *node {
// 	// if the child node already exists, and the incoming segment is a wildcard,
// 	// they will conflict with each other
// 	/* Note 1: node & segment are interchangeable concepts in the above context
// 	   Note 2: if the incoming segment is a wildcard, it also implies that we are
// 			   using func:matchChild in func:insert, not func:search, which means we
// 			   are in the process of constructing the tree. Since the request URL
// 			   will not contain wildcard.
// 	*/
// 	// if (len(n.children) > 0) && (strings.HasPrefix(segment, ":") || strings.HasPrefix(segment, "*")) {
// 	// 	panic(fmt.Sprintf(
// 	// 		"Wildcard segment %s conflicts with existing routers",
// 	// 		segment,
// 	// 	))
// 	// }

// 	// if the incoming segment is not wildcard
// 	for _, child := range n.children {
// 		if child.segment == segment || child.isWildcard {
// 			return child
// 		}
// 	}
// 	return nil
// }

// param:pattern must be a valid URL
func parsePattern(pattern string) (segments []string) {
	if len(pattern) == 1 {
		segments = []string{"root"}
	} else {
		segments = strings.Split(pattern, "/")
		segments[0] = "root"
		if segments[len(segments)-1] == "" {
			segments[len(segments)-1] = "/"
		}
	}
	return
}

func validateRoutingPattern(pattern string) {
	if !strings.HasPrefix(pattern, "/") {
		panic(fmt.Sprintf("Input routing pattern is %s\nIt should starts with /.", pattern))
	}
}
