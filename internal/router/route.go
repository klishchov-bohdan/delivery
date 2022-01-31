package router

import (
	"net/http"
	"strings"
)

type Node struct {
	childHandlers map[string]*Node
	parentHandler map[string]*Node
	fullPath      string
	handler       http.Handler
}

func (n *Node) CreateChild(fullPath string, handler http.Handler) (*Node, error) {
	strings.Split(fullPath)
}
