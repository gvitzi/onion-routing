package node

import "fmt"

type Router struct {
	Nodes map[string]*OnionNode
	Input chan Message
	stop  chan bool
}

func NewRouter() *Router {
	r := &Router{
		Nodes: make(map[string]*OnionNode),
		Input: make(chan Message),
		stop:  make(chan bool),
	}

	return r
}

func (r *Router) AddNode(n *OnionNode) {
	r.Nodes[n.Name] = n
}

func (r *Router) Start() {
	for {
		select {
		case msg := <-r.Input:
			fmt.Printf("Router: Routing from %s to %s data: %s\n", msg.Source, msg.Dest, msg.Payload)
			node := r.Nodes[msg.Dest]
			node.in <- msg
		case <-r.stop:
			return
		}
	}
}

func (r *Router) Stop() {
	r.stop <- true
}
