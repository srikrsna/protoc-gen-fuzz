package pb_test

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	pb "github.com/srikrsna/protoc-gen-fuzz/example"
	pbfuzz "github.com/srikrsna/protoc-gen-fuzz/example/fuzz"
)

func TestFuzz(t *testing.T) {
	fz := fuzz.New().Funcs(pbfuzz.FuzzFuncs()...)

	var msg pb.SomeMessage
	fz.Fuzz(&msg)

	// Test using random msg
}
