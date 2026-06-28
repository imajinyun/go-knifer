package vhash_test

import (
	"fmt"
	"hash/fnv"

	"github.com/imajinyun/knifer-go/vhash"
)

func ExampleAdditiveHash() {
	fmt.Println(vhash.AdditiveHash("abc", 31))
	// Output: 18
}

func ExampleJavaDefaultHash() {
	// Equivalent to Java String.hashCode.
	fmt.Println(vhash.JavaDefaultHash("a"))
	// Output: 97
}

func ExampleFnvHash() {
	fmt.Println(vhash.FnvHash("abc"))
	// Output: 1134309195
}

func ExampleHash32() {
	fmt.Println(vhash.Hash32("abc", fnv.New32a))
	fmt.Println(vhash.Hash32("abc", nil))
	// Output:
	// 440920331
	// 1134309195
}

func ExampleFnvHashString() {
	fmt.Println(vhash.FnvHashString("abc"))
	// Output: 33957123
}

func ExampleRsHash() {
	fmt.Println(vhash.RsHash("abc"))
	// Output: 822160044
}

func ExampleJsHash() {
	fmt.Println(vhash.JsHash("abc"))
	// Output: 895805535
}

func ExamplePjwHash() {
	fmt.Println(vhash.PjwHash("abc"))
	// Output: 26499
}

func ExampleElfHash() {
	fmt.Println(vhash.ElfHash("abc"))
	// Output: 26499
}

func ExampleBkdrHash() {
	fmt.Println(vhash.BkdrHash("a"))
	// Output: 97
}

func ExampleSdbmHash() {
	fmt.Println(vhash.SdbmHash("abc"))
	// Output: 807794786
}

func ExampleDjbHash() {
	fmt.Println(vhash.DjbHash("a"))
	// Output: 177670
}

func ExampleApHash() {
	fmt.Println(vhash.ApHash("abc"))
	// Output: -25651485
}

func ExampleHfHash() {
	fmt.Println(vhash.HfHash("abc"))
	// Output: 888
}

func ExampleHfIpHash() {
	fmt.Println(vhash.HfIpHash("10.0.0.1"))
	// Output: 32
}

func ExampleTianlHash() {
	fmt.Println(vhash.TianlHash("abc"))
	// Output: 33734718
}

func ExampleNewConsistentHash() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(8))
	ring.Add("cache-a")
	ring.Add("cache-b")
	ring.Add("cache-c")

	node, err := ring.Get("user:42")
	if err != nil {
		panic(err)
	}
	replicas, err := ring.GetN("user:42", 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(node != "")
	fmt.Println(len(replicas))
	// Output:
	// true
	// 2
}

func ExampleConsistentHash_Add() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(4))
	ring.Add("cache-a")

	node, err := ring.Get("user:42")
	if err != nil {
		panic(err)
	}
	fmt.Println(node)
	// Output: cache-a
}

func ExampleConsistentHash_Remove() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(4))
	ring.Add("cache-a")
	ring.Add("cache-b")
	before, err := ring.Get("user:42")
	if err != nil {
		panic(err)
	}

	ring.Remove(before)
	after, err := ring.Get("user:42")
	if err != nil {
		panic(err)
	}
	fmt.Println(before != after)
	// Output: true
}

func ExampleConsistentHash_Get() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(4))
	ring.Add("cache-a")
	node, err := ring.Get("asset:logo")
	if err != nil {
		panic(err)
	}
	fmt.Println(node)
	// Output: cache-a
}

func ExampleConsistentHash_GetN() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(4))
	ring.Add("cache-a")
	ring.Add("cache-b")
	nodes, err := ring.GetN("asset:logo", 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(nodes), nodes[0] != nodes[1])
	// Output: 2 true
}

func ExampleWithVirtualNodes() {
	ring := vhash.NewConsistentHash(vhash.WithVirtualNodes(4))
	ring.Add("cache-a")
	ring.Add("cache-b")

	node, err := ring.Get("asset:logo")
	if err != nil {
		panic(err)
	}
	fmt.Println(node != "")
	// Output: true
}

func ExampleWithReplicaCount() {
	ring := vhash.NewConsistentHash(vhash.WithReplicaCount(4))
	ring.Add("cache-a")
	ring.Add("cache-b")

	nodes, err := ring.GetN("asset:logo", 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(nodes))
	// Output: 2
}

func ExampleWithHashFunc() {
	hashFunc := func(data []byte) uint64 {
		var sum uint64
		for _, b := range data {
			sum = sum*131 + uint64(b)
		}
		return sum
	}
	ring := vhash.NewConsistentHash(
		vhash.WithVirtualNodes(2),
		vhash.WithHashFunc(hashFunc),
	)
	ring.Add("cache-a")
	ring.Add("cache-b")

	node, err := ring.Get("asset:logo")
	if err != nil {
		panic(err)
	}
	fmt.Println(node != "")
	// Output: true
}
