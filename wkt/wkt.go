// Package wkt provides Fuzz functions for protobuf well known types (google.protobuf.*)
package wkt

import (
	"time"

	fuzz "github.com/google/gofuzz"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var FuzzWKT = [...]interface{}{FuzzAny, FuzzDuration, FuzzStruct, FuzzTimestamp}

// FuzzAny can be used to Fuzz google.protobuf.Any messages.
// It fills a fuzzed google.protobuf.Struct message.
func FuzzAny(msg *anypb.Any, c fuzz.Continue) {
	var st structpb.Struct
	c.Fuzz(&st)
	msg.MarshalFrom(&st)
}

// FuzzStruct can be used to Fuzz google.protobuf.Struct messages.
// # of fields will be in [0, 64). Each field will be one of float, bool, or string
func FuzzStruct(msg *structpb.Struct, c fuzz.Continue) {
	fc := c.Rand.Int() % 64
	msg.Fields = make(map[string]*structpb.Value, fc)
	for i := 0; i < fc; i++ {
		var v *structpb.Value
		switch c.Int() % 3 {
		case 0:
			v = structpb.NewNumberValue(c.Float64())
		case 1:
			v = structpb.NewBoolValue(c.RandBool())
		case 2:
			v = structpb.NewStringValue(c.RandString())
		}

		msg.Fields[c.RandString()] = v
	}
}

// FuzzTimestamp can be used to Fuzz google.protobuf.Timestamp messages.
// It's range is same as default time.Time fuzzer of gofuzz
func FuzzTimestamp(msg *timestamppb.Timestamp, c fuzz.Continue) {
	var t time.Time
	c.Fuzz(&t)
	msg.Seconds = t.Unix()
	msg.Nanos = 0
}

// FuzzDuration can be used to Fuzz google.protobuf.Duration messages.
// It's range is [0, 1 year)
func FuzzDuration(msg *durationpb.Duration, c fuzz.Continue) {
	msg.Seconds = c.Int63n(365 * 24 * 60 * 60)
}
