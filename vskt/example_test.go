package vskt_test

import (
	"bytes"
	"fmt"

	"github.com/imajinyun/go-knifer/vskt"
)

func ExampleGetRemoteAddress() {
	fmt.Println(vskt.GetRemoteAddress(nil) == nil)
	// Output: true
}

func ExampleFuncEncoder() {
	encoder := vskt.FuncEncoder[string](func(_ *vskt.AioSession, b *bytes.Buffer, data string) {
		b.WriteString("encoded:" + data)
	})
	var out bytes.Buffer

	encoder.Encode(nil, &out, "payload")
	fmt.Println(out.String())
	// Output: encoded:payload
}

func ExampleNewSocketErrorf() {
	err := vskt.NewSocketErrorf("socket %s", "closed")

	fmt.Println(err.Error())
	// Output: socket closed
}

func ExampleNewSocketErrorMsg() {
	err := vskt.NewSocketErrorMsg("connection reset")

	fmt.Println(err.Error())
	// Output: connection reset
}
