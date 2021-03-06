package main

import (
	"fmt"

	"github.com/hashicorp/raft-boltdb"
	"github.com/komuw/kshaka"
)

func main() {
	// The store should, ideally be disk persisted.
	// Any that implements hashicorp/raft StableStore interface will suffice
	boltStore, err := raftboltdb.NewBoltStore("/tmp/bolt.db")
	if err != nil {
		panic(err)
	}

	// The function that will be applied by CASPaxos.
	// This will be applied to the current value stored
	// under the key passed into the Propose method of the proposer.
	var setFunc = func(val []byte) kshaka.ChangeFunction {
		return func(current []byte) ([]byte, error) {
			return val, nil
		}
	}

	// Note that, in practice, nodes ideally should be
	// in different machines each with its own store.
	node1 := kshaka.NewNode(1, boltStore)
	node2 := kshaka.NewNode(2, boltStore)
	node3 := kshaka.NewNode(3, boltStore)

	transport1 := &kshaka.InmemTransport{Node: node1}
	transport2 := &kshaka.InmemTransport{Node: node2}
	transport3 := &kshaka.InmemTransport{Node: node3}
	node1.AddTransport(transport1)
	node2.AddTransport(transport2)
	node3.AddTransport(transport3)

	kshaka.MingleNodes(node1, node2, node3)

	key := []byte("name")
	val := []byte("Masta-Ace")

	// make a proposition; consensus via CASPaxos will happen
	newstate, err := node2.Propose(key, setFunc(val))
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("\n newstate: %v \n", newstate)
}
