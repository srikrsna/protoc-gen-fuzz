package wkt_test

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/srikrsna/protoc-gen-fuzz/wkt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFuzzAny(t *testing.T) {
	fz := fuzz.New().Funcs(wkt.FuzzAny, wkt.FuzzStruct)
	var exp anypb.Any
	fz.Fuzz(&exp)

	buf, err := protojson.Marshal(&exp)
	if err != nil {
		t.Fatal(err)
	}

	var act anypb.Any
	if err := protojson.Unmarshal(buf, &act); err != nil {
		t.Fatal(err)
	}
}

func TestFuzzStruct(t *testing.T) {
	fz := fuzz.New().Funcs(wkt.FuzzStruct)
	var exp structpb.Struct
	fz.Fuzz(&exp)
	if !exp.ProtoReflect().IsValid() {
		t.Fatal("Invalid")
	}
}

func TestFuzzTimestamp(t *testing.T) {
	fz := fuzz.New().Funcs(wkt.FuzzTimestamp)
	var exp timestamppb.Timestamp
	fz.Fuzz(&exp)
	if err := exp.CheckValid(); err != nil {
		t.Fatal(err)
	}
}

func TestFuzzDuration(t *testing.T) {
	fz := fuzz.New().Funcs(wkt.FuzzDuration)
	var exp durationpb.Duration
	fz.Fuzz(&exp)
	if err := exp.CheckValid(); err != nil {
		t.Fatal(err)
	}
}
