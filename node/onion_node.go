package node

import (
	"fmt"
	"gvitzi/onion-routing/crypt"
	"time"
)

type OnionNode struct {
	Node
	KeyPair    crypt.KeyPair
	encryptAlg crypt.Encryption
	stop       chan bool
}

func NewOnionNode(name string, keyPair crypt.KeyPair, out chan Message) *OnionNode {
	base := NewNode(name, out)
	n := &OnionNode{
		*base,
		keyPair,
		crypt.BasicEncryption{},
		make(chan bool),
	}

	return n
}

func (n *OnionNode) Start() {
	for {
		fmt.Printf("%s: waiting for next msg\n", n.Name)
		select {
		case msg := <-n.in:
			fmt.Printf("%s: received \"%++v\"\n", n.Name, msg)
			go n.handleMessage(msg)
		case <-n.stop:
			return
		}
	}
}

func (n *OnionNode) handleMessage(msg Message) {
	fmt.Printf("%s: parsing data \"%s\"\n", n.Name, msg.Payload)

	nextHop := DecryptNextHop([]byte(msg.Payload), n.KeyPair.Pvt, n.encryptAlg)

	if nextHop.Data != nil {
		fmt.Printf("%s: I got a message for me !!! \"%s\"\n", n.Name, *nextHop.Data)
	} else {
		fmt.Printf("%s: This message is for someone else - routing to: \"%s\"\n", n.Name, nextHop.Dest)
		newMsg := Message{
			Source:  n.Name,
			Dest:    nextHop.Dest,
			Payload: string(*nextHop.Payload),
		}

		time.Sleep(1 * time.Second)
		n.Send(newMsg)
	}
}

func (n *OnionNode) Stop() {
	n.stop <- true
	return
}

func (n *OnionNode) GetNodeKeyPair() NodeKeyPair {
	return NodeKeyPair{
		Node:   n.Name,
		PubKey: n.KeyPair.Pub,
	}
}
