package bmfcommon

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestClearRegistrations(t *testing.T) {
	ClearRegistrations()

	if len(boxMapping) != 0 {
		t.Fatalf("Expected no registrations at top of test.")
	}

	bf := testBox1Factory{}
	RegisterBoxType(bf)

	if len(boxMapping) != 1 {
		t.Fatalf("Expected exactly one registration.")
	}

	ClearRegistrations()

	if len(boxMapping) != 0 {
		t.Fatalf("Expected no registrations at bottom of test.")
	}
}

func TestRegisterBoxType_Ok(t *testing.T) {
	ClearRegistrations()
	defer ClearRegistrations()

	bf := testBox1Factory{}
	RegisterBoxType(bf)

	if len(boxMapping) != 1 {
		t.Fatalf("Expected exactly one registration.")
	}
}

func TestRegisterBoxType_AlreadyRegistered(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if err.Error() != "box-factory already registered: [tb1 ]" {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected panic.")
		}
	}()

	ClearRegistrations()
	defer ClearRegistrations()

	bf := testBox1Factory{}
	RegisterBoxType(bf)
	RegisterBoxType(bf)
}

func TestGetFactory(t *testing.T) {
	ClearRegistrations()
	defer ClearRegistrations()

	bf1 := testBox1Factory{}
	RegisterBoxType(bf1)

	bf2 := testBox2Factory{}
	RegisterBoxType(bf2)

	bf3 := testBox3Factory{}
	RegisterBoxType(bf3)

	bf1Recovered := GetFactory("tb1 ")
	_ = bf1Recovered.(testBox1Factory)

	bf2Recovered := GetFactory("tb2 ")
	_ = bf2Recovered.(testBox2Factory)
}

func TestBoxNameIsValid_Hit(t *testing.T) {
	if BoxNameIsValid("abc") != true {
		t.Fatalf("Expected valid box name.")
	}
}

func TestBoxNameIsValid_Miss_Empty(t *testing.T) {
	if BoxNameIsValid("") != false {
		t.Fatalf("Expected invalid box name.")
	}
}

func TestBoxNameIsValid_Miss_Whitespace(t *testing.T) {
	if BoxNameIsValid(" ") != false {
		t.Fatalf("Expected invalid box name.")
	}
}

func TestBoxNameIsValid_Miss_Symbols(t *testing.T) {
	if BoxNameIsValid("abc_") != false {
		t.Fatalf("Expected invalid box name.")
	}
}

func TestGetParentBoxName_Root(t *testing.T) {
	if GetParentBoxName(nil) != "ROOT" {
		t.Fatalf("Expected root parent name.")
	}
}

func TestGetParentBoxName_Nonroot(t *testing.T) {
	tb1 := &testBox1{
		Box: Box{
			name: "parent",
		},
	}

	if GetParentBoxName(tb1) != "parent" {
		t.Fatalf("Expected node name: [%s]", GetParentBoxName(tb1))
	}
}

func TestChildBoxes(t *testing.T) {
	boxes1 := []CommonBox{&testBox1{}}
	boxes2 := []CommonBox{&testBox2{}}

	lbi := LoadedBoxIndex{
		"abc": boxes1,
		"def": boxes2,
	}

	tb4 := &testBox4{
		LoadedBoxIndex: lbi,
	}

	actual1 := ChildBoxes(tb4, "abc")

	if reflect.DeepEqual(actual1, boxes1) != true {
		t.Fatalf("Child boxes not correct (1).")
	}

	actual2 := ChildBoxes(tb4, "def")

	if reflect.DeepEqual(actual2, boxes2) != true {
		t.Fatalf("Child boxes not correct (2).")
	}
}
