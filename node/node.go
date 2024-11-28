package node

type Node struct {
	Name string
	in chan Message
	out chan Message
}

func NewNode(name string, out chan Message) *Node{
	n := &Node{
		Name: name,
		in: make(chan Message),
		out: out,

	}

	return n
}

func (n *Node) Receive() Message {
	return <-n.in
}

func (n *Node) Send(msg Message) {
	n.out<-msg
}

