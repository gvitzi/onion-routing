package node

import (
	"encoding/json"
	"fmt"
	"gvitzi/onion-routing/crypt"
)

type Hop struct {
	Payload *[]byte // Marshalled 'Hop' struct encrypted with nextHops pub key - on last layer this is empty string ''
	Dest    string  // clear Dest - on last layer this is empty string ''
	Data    *string // Only set on last hop - Encrypted with Dest pub-key
}

// General message that can be routed
// The router simulates IP Routers and should have an unencrypted Dest value
type Message struct {
	Source  string // Source node name - for debugging only
	Dest    string // unencrypted NextHop
	Payload string // Marshalled 'Hop' struct encrypted with nextHops pub key
}

type NodeKeyPair struct {
	Node   string
	PubKey string
}

func FromJSON(jsonData string) Message {
	var msg Message
	byteData := []byte(jsonData)

	if err := json.Unmarshal(byteData, &msg); err != nil {
		panic(err)
	}
	fmt.Printf("FromJson: %s\n", msg)
	return msg
}

func (m Message) ToJSON() string {
	jsonDataBytes, _ := json.Marshal(m)
	jsonData := string(jsonDataBytes)
	fmt.Printf("ToJson: %s\n", jsonData)
	return jsonData
}

func BuildHops(data *string, hops []NodeKeyPair, e crypt.Encryption) *Hop {
	nextHop := &Hop{
		Payload: nil,
		Dest:    "",
		Data:    data,
	}

	for _, nodeKeyPair := range hops {
		nextHop = WrapNextHop(nextHop, nodeKeyPair.Node, nodeKeyPair.PubKey, e)
	}

	return nextHop
}

func WrapNextHop(hop *Hop, dest string, pubKey string, e crypt.Encryption) *Hop {
	encryptedHop := EncryptHop(hop, pubKey, e)
	return &Hop{
		Payload: &encryptedHop,
		Dest:    dest,
		Data:    nil,
	}
}

func EncryptHop(hop *Hop, pubKey string, e crypt.Encryption) []byte {
	jsonDataBytes, _ := json.Marshal(hop)
	encryptedbytes := e.Encrypt(jsonDataBytes, pubKey)
	return encryptedbytes
}

func DecryptNextHop(payload []byte, pvtKey string, e crypt.Encryption) *Hop {
	if payload == nil {
		return nil
	}

	decryptedBytes := e.Decrypt(payload, pvtKey)
	nextHop := Hop{}
	if err := json.Unmarshal(decryptedBytes, &nextHop); err != nil {
		panic(err)
	}
	fmt.Printf("nextHop dest: %s, NextHop: %+v\n", nextHop.Dest, nextHop)

	return &nextHop
}

// Builds a Message to send on the wire
// nodeList is a list of hops
func BuildMessage(src string, data *string, nodeList []NodeKeyPair, e crypt.Encryption) *Message {
	// the last node is the first hop
	firstHop := nodeList[len(nodeList)-1]

	// remove it from the list
	nodeList = nodeList[:len(nodeList)-1]

	// generate recursive hops strucure
	hops := BuildHops(data, nodeList, e)

	// encrypt structure
	encryptedHop := EncryptHop(hops, firstHop.PubKey, e)

	// create a message to dest with the encrypted structure
	msg := Message{Source: src, Dest: firstHop.Node, Payload: string(encryptedHop[:])}
	return &msg
}
