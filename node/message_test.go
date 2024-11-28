package node

import (
	"fmt"
	"gvitzi/onion-routing/crypt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleLayer(t *testing.T) {

	e := crypt.BasicEncryption{}

	nodeName := "n1"
	pubKey := "abc"
	pvtKey := "abc"

	data := "this-is-my-data"
	hop := &Hop{
		Dest: "d1",
		Data: &data,
	}
	wrappedHop := WrapNextHop(hop, nodeName, pubKey, e)
	fmt.Printf("%+v\n", wrappedHop)

	unwrappedHop := DecryptNextHop(*wrappedHop.Payload, pvtKey, e)

	if unwrappedHop.Data == nil {
		t.Fatalf("expected to get ot data here but data is nil. nextHop: %++v", unwrappedHop)
	}

	if *unwrappedHop.Data != data {
		t.Fatalf("data corrupted expected: %s got: %s", data, *unwrappedHop.Data)
	}
}

func TestBuildHops(t *testing.T) {
	e := crypt.BasicEncryption{}

	n1pvt := "a"
	n2pvt := "d"
	n3pvt := "h"

	n1 := NodeKeyPair{Node: "n1", PubKey: "abc"}
	n2 := NodeKeyPair{Node: "n2", PubKey: "def"}
	n3 := NodeKeyPair{Node: "n3", PubKey: "hig"}

	nodeList := make([]NodeKeyPair, 0)
	nodeList = append(nodeList, n1)
	nodeList = append(nodeList, n2)
	nodeList = append(nodeList, n3)

	data := "this-is-my-data"
	hops := BuildHops(&data, nodeList, e)
	fmt.Printf("%+v\n", hops)

	nextHop := hops

	nextHop = DecryptNextHop(*nextHop.Payload, n3pvt, e)
	nextHop = DecryptNextHop(*nextHop.Payload, n2pvt, e)
	nextHop = DecryptNextHop(*nextHop.Payload, n1pvt, e)

	if nextHop.Data == nil {
		t.Fatalf("expected to get ot data here but data is nil. nextHop: %++v", nextHop)
	}

	if *nextHop.Data != data {
		t.Fatalf("data corrupted expected: %s got: %s", data, *nextHop.Data)
	}
}
func TestBuildMessage(t *testing.T) {
	e := crypt.BasicEncryption{}

	n1Keys := crypt.KeyPair{Pvt: "aebc", Pub: "abc"}
	n2Keys := crypt.KeyPair{Pvt: "def", Pub: "def"}

	nodeList := make([]NodeKeyPair, 0)
	nodeList = append(nodeList, NodeKeyPair{Node: "n2", PubKey: n2Keys.Pub}) //dest
	nodeList = append(nodeList, NodeKeyPair{Node: "n1", PubKey: n1Keys.Pub}) // 1st hop

	data := "this-is-data"
	msg := BuildMessage("n0", &data, nodeList, e)
	assert.Equal(t, msg.Dest, "n1")
	assert.Equal(t, msg.Source, "n0")

	fmt.Printf("parsing data \"%s\"\n", msg.Payload)
	nextHop := DecryptNextHop([]byte(msg.Payload), n1Keys.Pvt, e)
	assert.Equal(t, nextHop.Dest, "n2")

	fmt.Printf("got hop: \"%++v\"\n", nextHop)
}
