package main

import (
	"fmt"
	"gvitzi/onion-routing/crypt"
	"gvitzi/onion-routing/node"
	"time"
)

func main() {
	fmt.Println("main")

	router := node.NewRouter()

	n1 := node.NewOnionNode("n1", crypt.KeyPair{Pvt: "aebc", Pub: "abc"}, router.Input)
	n2 := node.NewOnionNode("n2", crypt.KeyPair{Pvt: "def", Pub: "def"}, router.Input)
	n3 := node.NewOnionNode("n3", crypt.KeyPair{Pvt: "hig", Pub: "hig"}, router.Input)
	n4 := node.NewOnionNode("n4", crypt.KeyPair{Pvt: "xyz", Pub: "xyz"}, router.Input)

	router.AddNode(n1)
	router.AddNode(n2)
	router.AddNode(n3)
	router.AddNode(n4)

	go router.Start()
	go n1.Start()
	go n2.Start()
	go n3.Start()
	go n4.Start()

	nodeList := make([]node.NodeKeyPair, 0)
	nodeList = append(nodeList, n1.GetNodeKeyPair()) // dest
	nodeList = append(nodeList, n2.GetNodeKeyPair()) // hop 3
	nodeList = append(nodeList, n3.GetNodeKeyPair()) // hop 2
	nodeList = append(nodeList, n4.GetNodeKeyPair()) // hop 1

	e := crypt.BasicEncryption{}
	data := "this-is-my-data"
	msg := node.BuildMessage("n0", &data, nodeList, e)

	router.Input <- *msg

	fmt.Println("sleeping")

	time.Sleep(4 * time.Second)

	n1.Stop()
	n2.Stop()
	n3.Stop()
	n4.Stop()
	router.Stop()
	fmt.Println("finished")
}
