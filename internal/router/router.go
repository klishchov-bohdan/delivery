package router

import (
	"net/http"
)

type Router struct {
	root *Node
}

func New(root http.Handler) (*Router, error) {
	return &Router{
		root: &Node{
			fullPath: "/",
			handler:  root,
		},
	}, nil
}
