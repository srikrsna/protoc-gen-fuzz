# protoc-gen-fuzz (PGF)

PGF is a protoc plugin used to generate fuzzing functions suitable for use with [gofuzz](https://github.com/google/gofuzz).
Get it using `go get github.com/srikrsna/protoc-gen-fuzz `

## Usage

For a protobuf definition like this,

```proto
syntax = "proto3";

package pb;

option go_package = "github.com/srikrsna/protoc-gen-fuzz/example;pb";

message SomeMessage {
  string string_field = 1;
  int32 int32_field = 2;
  bool bool_field = 3;
  InnerMessage inner_message = 4;
  oneof some_oneof {
    string oneof_string = 5;
    bool oneof_bool = 6;
  }
  repeated string slice_string = 7;
  repeated InnerMessage slice_message = 8;
}

message InnerMessage { string inner_field = 1; }
```

Fuzz functions can be generated using

```bash
protoc -I . --fuzz_out=:. example/example.proto
```

and can be used in testing as follow,

```go
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
```
