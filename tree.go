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
						"Wildcard segment '%s' in path '%s' conflicts with existing segment '%s'",
						segment,
						strings.Join(segments, "/"),
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
		}
		existingChild.handler = handler
		return
	}

	if !isWildSegment && isLastSegment {
		needNewChild := true
		existingChild := &node{}

		for _, child := range n.children {
			if child.isWildcard && child.handler != nil {
				panic(
					fmt.Sprintf(
						"Segment '%s' in path '%s' conflicts with existing wildcard segment '%s'",
						segment,
						strings.Join(segments, "/"),
						child.segment,
					),
				)
			} else if segment == child.segment && child.handler != nil {
				panic(
					fmt.Sprintf(
						"Segment '%s' in path '%s' conflicts with existing segment '%s'",
						segment,
						strings.Join(segments, "/"),
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
		}
		existingChild.handler = handler
		return
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
					if strings.HasPrefix(child.segment, ":") {
						// child.segment is a colon wildcard
						params[child.segment] = segment
						child.searchHelper(segments, level+1, searchResultPtr, params)
					} else {
						// child.segment is an asteroid wildcard
						params[child.segment] = strings.Join(segments[level:], "/")
						*searchResultPtr = *child
						return
					}
				}
			}
		}
	}
}

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
	if pattern == "/" {
		return
	}

	if !strings.HasPrefix(pattern, "/") {
		panic(fmt.Sprintf("Input routing pattern is %s\nIt should starts with /.", pattern))
	}

	allSegments := strings.Split(pattern, "/")
	segments := allSegments[1:] // omit the first empty string, which was "/"(root) before Split
	lenOfSegments := len(segments)
	for i, segment := range segments {
		if i == (lenOfSegments - 1) {
			if segment == "*" {
				panic(fmt.Sprint("Wildcard must have a name. e.g. *filepath"))
			}
		} else {
			if segment == "" {
				panic(fmt.Sprint("Consecutive slashes in a routing pattern are not allowed."))
			} else if segment == ":" {
				panic(fmt.Sprint("Wildcard must have a name. e.g. :goose"))
			} else if strings.HasPrefix(segment, "*") {
				panic(fmt.Sprint("The asteroid wildcard can only be used at the end of a routing pattern."))
			}
		}
	}
}
