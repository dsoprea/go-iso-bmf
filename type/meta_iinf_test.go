package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestNewIinfBox(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	if iinf.itemsById == nil {
		t.Fatalf("itemsById not initialized.")
	}

	if iinf.itemsByName == nil {
		t.Fatalf("itemsByName not initialized.")
	}
}

func TestIinfBox_loadItem(t *testing.T) {

	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   22,
		itemType: iit2,
		itemName: "def",
	}

	iinf.loadItem(infe2)

	// Validate.

	if len(iinf.itemsById) != 2 {
		t.Fatalf("itemsById is not the right length.")
	} else if len(iinf.itemsByName) != 2 {
		t.Fatalf("itemsByName is not the right length.")
	}

	if _, found := iinf.itemsById[11]; found != true {
		t.Fatalf("First item ID not found.")
	} else if _, found := iinf.itemsById[22]; found != true {
		t.Fatalf("Second item ID not found.")
	}

	if _, found := iinf.itemsByName["abc"]; found != true {
		t.Fatalf("First item name not found.")
	} else if _, found := iinf.itemsByName["def"]; found != true {
		t.Fatalf("Second item name not found.")
	}
}

func TestIinfBox_loadItem_DuplicateId(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)

			if err.Error() != "item with ID (11) occurs more than once" {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected error.")
		}
	}()

	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   11,
		itemType: iit2,
		itemName: "def",
	}

	iinf.loadItem(infe2)
}

func TestIinfBox_loadItem_DuplicateName(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)

			if err.Error() != "item with name [abc] occurs more than once" {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected error.")
		}
	}()

	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   22,
		itemType: iit2,
		itemName: "abc",
	}

	iinf.loadItem(infe2)
}

func TestIinfBox_GetItemWithId_Hit(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   22,
		itemType: iit2,
		itemName: "def",
	}

	iinf.loadItem(infe2)

	// Validate.

	recoveredInfe1, err := iinf.GetItemWithId(11)
	log.PanicIf(err)

	if recoveredInfe1 != infe1 {
		t.Fatalf("First item not correct.")
	}

	recoveredInfe2, err := iinf.GetItemWithId(22)
	log.PanicIf(err)

	if recoveredInfe2 != infe2 {
		t.Fatalf("Second item not correct.")
	}
}

func TestIinfBox_GetItemWithId_Miss(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	_, err := iinf.GetItemWithId(22)
	if err == nil {
		t.Fatalf("Expected error.")
	} else if log.Is(err, ErrNoItemsFound) != true {
		log.Panic(err)
	}
}

func TestIinfBox_GetItemWithName(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   22,
		itemType: iit2,
		itemName: "def",
	}

	iinf.loadItem(infe2)

	// Validate.

	recoveredInfe1, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe1 != infe1 {
		t.Fatalf("First item not correct.")
	}

	recoveredInfe2, err := iinf.GetItemWithName("def")
	log.PanicIf(err)

	if recoveredInfe2 != infe2 {
		t.Fatalf("Second item not correct.")
	}
}

func TestIinfBox_GetItemWithName_Miss(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	_, err := iinf.GetItemWithName("def")
	if err == nil {
		t.Fatalf("Expected error.")
	} else if log.Is(err, ErrNoItemsFound) != true {
		log.Panic(err)
	}
}

func TestIinfBox_InlineString(t *testing.T) {
	iinf := newIinfBox(bmfcommon.Box{})

	// Add first item.

	iit1 := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe1 := &InfeBox{
		itemId:   11,
		itemType: iit1,
		itemName: "abc",
	}

	iinf.loadItem(infe1)

	// Add second item.

	iit2 := InfeItemTypeFromBytes([4]byte{'e', 'f', 'g', 'h'})

	infe2 := &InfeBox{
		itemId:   22,
		itemType: iit2,
		itemName: "def",
	}

	iinf.loadItem(infe2)

	if iinf.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) ENTRY-COUNT=(0) LOADED-ITEMS=(2)" {
		t.Fatalf("InlineString() not correct: [%s]", iinf.InlineString())
	}
}

func TestIinfBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.Boxes, 0)

	iinf := new(IinfBox)
	iinf.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(iinf.LoadedBoxIndex, lbi.Index()) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestIinfBoxFactory_Name(t *testing.T) {
	factory := iinfBoxFactory{}
	if factory.Name() != "iinf" {
		t.Fatalf("Name() not correct.")
	}
}

func TestIinfBoxFactory_New_Version0(t *testing.T) {
	// Load

	var data []byte

	// Version.
	version := uint8(0)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&data, uint16(11))

	var b []byte
	bmfcommon.PushBox(&b, "iinf", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := iinfBoxFactory{}.New(box)
	log.PanicIf(err)

	iinf := cb.(*IinfBox)

	if iinf.entryCount != 11 {
		t.Fatalf("entryCount not correct.")
	}
}

func TestIinfBoxFactory_New_Version1(t *testing.T) {
	// Load

	var data []byte

	// Version.
	version := uint8(1)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&data, uint32(11))

	var b []byte
	bmfcommon.PushBox(&b, "iinf", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := iinfBoxFactory{}.New(box)
	log.PanicIf(err)

	iinf := cb.(*IinfBox)

	if iinf.entryCount != 11 {
		t.Fatalf("entryCount not correct.")
	}
}
