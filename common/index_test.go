package bmfcommon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestFullyQualifiedBoxName_String(t *testing.T) {
	fqbn := FullyQualifiedBoxName{
		"aa",
		"bb",
		"cc",
	}

	if fqbn.String() != "aa.bb.cc" {
		t.Fatalf("String() not correct.")
	}
}

func TestIndexedBoxEntry_String(t *testing.T) {
	ibe := IndexedBoxEntry{
		NamePhrase:     "test_name",
		SequenceNumber: 99,
	}

	if ibe.String() != "test_name(99)" {
		t.Fatalf("String() not correct: [%s]", ibe.String())
	}
}

func TestFullBoxIndex_getBoxName(t *testing.T) {
	tb1 := &testBox1{
		Box: Box{name: "tb1"},
	}

	tb2 := &testBox1{
		Box: Box{parent: tb1, name: "tb2"},
	}

	tb3 := &testBox1{
		Box: Box{parent: tb2, name: "tb3"},
	}

	tb4 := &testBox1{
		Box: Box{parent: tb3, name: "tb4"},
	}

	var fbi FullBoxIndex
	fqbn := fbi.getBoxName(tb4)

	if fqbn.String() != "tb1.tb2.tb3.tb4" {
		t.Fatalf("Box name is not correct: [%s]", fqbn.String())
	}
}

func TestFullBoxIndex_Add(t *testing.T) {
	fbi := make(FullBoxIndex)

	tb1 := &testBox1{
		Box: Box{name: "tb1"},
	}

	fbi.Add(tb1)

	tb2a := &testBox2{
		Box: Box{name: "tb2"},
	}

	fbi.Add(tb2a)

	tb2b := &testBox2{
		Box: Box{name: "tb2"},
	}

	fbi.Add(tb2b)

	expected := FullBoxIndex{
		IndexedBoxEntry{NamePhrase: "tb1", SequenceNumber: 0}: tb1,
		IndexedBoxEntry{NamePhrase: "tb2", SequenceNumber: 0}: tb2a,
		IndexedBoxEntry{NamePhrase: "tb2", SequenceNumber: 1}: tb2b,
	}

	if reflect.DeepEqual(fbi, expected) != true {
		for ibe, box := range fbi {
			fmt.Printf("[%s] (%d) %v\n", ibe.NamePhrase, ibe.SequenceNumber, box)
		}

		t.Fatalf("FBI not correct.")
	}
}

func TestFullBoxIndex_Dump(t *testing.T) {
	fbi := make(FullBoxIndex)

	tb1 := &testBox1{
		Box: Box{name: "tb1"},
	}

	fbi.Add(tb1)

	tb2a := &testBox2{
		Box: Box{name: "tb2"},
	}

	fbi.Add(tb2a)

	tb2b := &testBox2{
		Box: Box{name: "tb2"},
	}

	fbi.Add(tb2b)

	fbi.Dump()
}

func TestLoadedBoxIndex_GetChildBoxes_Hit(t *testing.T) {
	boxes := []CommonBox{&testBox1{}}

	lbi := LoadedBoxIndex{
		"abc": boxes,
	}

	recovered, err := lbi.GetChildBoxes("abc")
	log.PanicIf(err)

	if reflect.DeepEqual(recovered, boxes) != true {
		t.Fatalf("Recovered boxes are not correct.")
	}
}

func TestLoadedBoxIndex_GetChildBoxes_Miss(t *testing.T) {
	lbi := LoadedBoxIndex{}

	_, err := lbi.GetChildBoxes("abc")
	if err == nil {
		t.Fatalf("Expected no children found.")
	} else if err.Error() != "child box not found: [abc]" {
		log.Panic(err)
	}
}

func TestLoadedBoxIndex_ChildrenTypes(t *testing.T) {
	lbi := LoadedBoxIndex{
		"abc": nil,
		"def": nil,
		"ghi": nil,
	}

	names := lbi.ChildrenTypes()

	expected := []string{
		"abc",
		"def",
		"ghi",
	}

	if reflect.DeepEqual(names, expected) != true {
		t.Fatalf("Box types are not correct: [%v]", names)
	}
}
