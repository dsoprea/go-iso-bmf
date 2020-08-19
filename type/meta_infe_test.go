package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestInfeItemTypeFromBytes(t *testing.T) {
	b4 := [4]byte{1, 2, 3, 4}
	iit := InfeItemTypeFromBytes(b4)

	if iit != 0x01020304 {
		t.Fatalf("Type not correct.")
	}
}

func TestInfeItemType_EqualsName_Hit(t *testing.T) {
	b4 := [4]byte{'1', '2', '3', '4'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.EqualsName("1234") != true {
		t.Fatalf("EqualsName() should have returned true.")
	}
}

func TestInfeItemType_EqualsName_Miss(t *testing.T) {
	b4 := [4]byte{'1', '2', '3', '4'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.EqualsName("5678") != false {
		t.Fatalf("EqualsName() should have returned false.")
	}
}

func TestInfeItemType_string(t *testing.T) {
	b4 := [4]byte{'1', '2', '3', '4'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.string() != "1234" {
		t.Fatalf("string() was not correct.")
	}
}

func TestInfeItemType_IsMime_Hit(t *testing.T) {
	b4 := [4]byte{'m', 'i', 'm', 'e'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.IsMime() != true {
		t.Fatalf("IsMime() should be true.")
	}
}

func TestInfeItemType_IsMime_Miss(t *testing.T) {
	b4 := [4]byte{'a', 'b', 'c', 'd'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.IsMime() != false {
		t.Fatalf("IsMime() should be false.")
	}
}

func TestInfeItemType_IsUri_Hit(t *testing.T) {
	b4 := [4]byte{'u', 'r', 'i', ' '}
	iit := InfeItemTypeFromBytes(b4)

	if iit.IsUri() != true {
		t.Fatalf("IsUri() should be true.")
	}
}

func TestInfeItemType_IsUri_Miss(t *testing.T) {
	b4 := [4]byte{'a', 'b', 'c', 'd'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.IsUri() != false {
		t.Fatalf("IsUri() should be false.")
	}
}

func TestInfeItemType_String_Printable(t *testing.T) {
	b4 := [4]byte{'a', 'b', 'c', 'd'}
	iit := InfeItemTypeFromBytes(b4)

	if iit.String() != "abcd" {
		t.Fatalf("String() should be false: [%s]", iit.String())
	}
}

func TestInfeItemType_String_Nonprintable(t *testing.T) {
	b4 := [4]byte{1, 2, 3, 4}
	iit := InfeItemTypeFromBytes(b4)

	if iit.String() != "TYPE<0x01020304>" {
		t.Fatalf("String() should be false: [%s]", iit.String())
	}
}

func TestInfeBox_ItemId(t *testing.T) {
	infe := &InfeBox{
		itemId: 11,
	}

	if infe.ItemId() != 11 {
		t.Fatalf("ItemId() not correct.")
	}
}

func TestInfeBox_ItemProtectionIndex(t *testing.T) {
	infe := &InfeBox{
		itemProtectionIndex: 11,
	}

	if infe.ItemProtectionIndex() != 11 {
		t.Fatalf("ItemProtectionIndex() not correct.")
	}
}

func TestInfeBox_ItemName(t *testing.T) {
	infe := &InfeBox{
		itemName: "abc",
	}

	if infe.ItemName() != "abc" {
		t.Fatalf("ItemName() not correct.")
	}
}

func TestInfeBox_ContentType(t *testing.T) {
	infe := &InfeBox{
		contentType: "abc",
	}

	if infe.ContentType() != "abc" {
		t.Fatalf("ContentType() not correct.")
	}
}

func TestInfeBox_ContentEncoding(t *testing.T) {
	infe := &InfeBox{
		contentEncoding: "abc",
	}

	if infe.ContentEncoding() != "abc" {
		t.Fatalf("ContentEncoding() not correct.")
	}
}

func TestInfeBox_ItemType(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe := &InfeBox{
		itemType: itemType,
	}

	if infe.ItemType() != itemType {
		t.Fatalf("ItemType() not correct.")
	}
}

func TestInfeBox_ExtensionType(t *testing.T) {
	infe := &InfeBox{
		extensionType: 11,
	}

	if infe.ExtensionType() != 11 {
		t.Fatalf("ExtensionType() not correct.")
	}
}

func TestInfeBox_ItemUriType(t *testing.T) {
	infe := &InfeBox{
		itemUriType: "abc",
	}

	if infe.ItemUriType() != "abc" {
		t.Fatalf("ItemUriType() not correct.")
	}
}

func TestInfeBox_InlineString_Version0(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe := &InfeBox{
		version:             0,
		itemType:            itemType,
		itemId:              11,
		itemProtectionIndex: 0,
		itemName:            "test-name",
	}

	if infe.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(0) ITEM-ID=(11) PROTECTION-INDEX=(0) NAME=[test-name] ITEM-TYPE=[abcd]" {
		t.Fatalf("InlineString() not correct: [%s]", infe.InlineString())
	}
}

func TestInfeBox_InlineString_Version1(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe := &InfeBox{
		version:             1,
		itemType:            itemType,
		itemId:              11,
		itemProtectionIndex: 0,
		itemName:            "test-name",
		extensionType:       11,
	}

	if infe.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(1) ITEM-ID=(11) PROTECTION-INDEX=(0) NAME=[test-name] ITEM-TYPE=[abcd] EXT-TYPE=(11)" {
		t.Fatalf("InlineString() not correct: [%s]", infe.InlineString())
	}
}

func TestInfeBox_InlineString_Version2_NoMimeNoUri(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})

	infe := &InfeBox{
		version:             2,
		itemType:            itemType,
		itemId:              11,
		itemProtectionIndex: 0,
		itemName:            "test-name",
		extensionType:       11,
	}

	if infe.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(2) ITEM-ID=(11) PROTECTION-INDEX=(0) NAME=[test-name] ITEM-TYPE=[abcd]" {
		t.Fatalf("InlineString() not correct: [%s]", infe.InlineString())
	}
}

func TestInfeBox_InlineString_Version2_Mime(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'m', 'i', 'm', 'e'})

	infe := &InfeBox{
		version:             2,
		itemType:            itemType,
		itemId:              11,
		itemProtectionIndex: 0,
		itemName:            "test-name",
		extensionType:       11,
	}

	if infe.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(2) ITEM-ID=(11) PROTECTION-INDEX=(0) NAME=[test-name] ITEM-TYPE=[mime] CONTENT-TYPE=[] CONTENT-ENCODING=[]" {
		t.Fatalf("InlineString() not correct: [%s]", infe.InlineString())
	}
}

func TestInfeBox_InlineString_Version2_Uri(t *testing.T) {
	itemType := InfeItemTypeFromBytes([4]byte{'u', 'r', 'i', ' '})

	infe := &InfeBox{
		version:             3,
		itemType:            itemType,
		itemId:              11,
		itemProtectionIndex: 0,
		itemName:            "test-name",
		extensionType:       11,
	}

	if infe.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(3) ITEM-ID=(11) PROTECTION-INDEX=(0) NAME=[test-name] ITEM-TYPE=[uri ] URI-TYPE=[]" {
		t.Fatalf("InlineString() not correct: [%s]", infe.InlineString())
	}
}

func TestInfeBoxFactory_Name(t *testing.T) {
	factory := infeBoxFactory{}
	if factory.Name() != "infe" {
		t.Fatalf("Name() not correct.")
	}
}

func TestInfeBox_New_Version0(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	// Construct INFE data (INFE boxes are children of IINF boxes).

	var infeData []byte

	version := uint8(0)
	bmfcommon.PushBytes(&infeData, []byte{version, 0, 0, 0})

	itemId := uint16(11)
	bmfcommon.PushBytes(&infeData, itemId)

	itemProtectionIndex := uint16(0)
	bmfcommon.PushBytes(&infeData, itemProtectionIndex)

	itemNameRaw := []byte("abc\000")
	bmfcommon.PushBytes(&infeData, itemNameRaw)

	contentTypeRaw := []byte("def\000")
	bmfcommon.PushBytes(&infeData, contentTypeRaw)

	contentEncodingRaw := []byte("ghi\000")
	bmfcommon.PushBytes(&infeData, contentEncodingRaw)

	// Construct an IINF.

	var iinfData []byte

	// Version.
	iinfVersion := uint8(0)
	bmfcommon.PushBytes(&iinfData, []byte{iinfVersion, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&iinfData, uint16(11))

	// Embed child INFE box.
	bmfcommon.PushBox(&iinfData, "infe", infeData)

	var b []byte
	bmfcommon.PushBox(&b, "iinf", iinfData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	fbi := resource.Index()

	if len(fbi) != 2 {
		t.Fatalf("Incorrect number of indexed boxes.")
	}

	// Validate INFE.

	infeCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf.infe", 0}]
	if found != true {
		t.Fatalf("INFE not indexed.")
	}

	infe := infeCb.(*InfeBox)

	if infe.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if infe.itemProtectionIndex != 0 {
		t.Fatalf("itemProtectionIndex not correct.")
	} else if infe.itemName != "abc" {
		t.Fatalf("itemName not correct.")
	} else if infe.contentType != "def" {
		t.Fatalf("contentType not correct.")
	} else if infe.contentEncoding != "ghi" {
		t.Fatalf("contentEncoding not correct.")
	}

	// Validate IINF (that the INFE is registered in it).

	iinfCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf", 0}]
	if found != true {
		t.Fatalf("IINF not indexed.")
	}

	iinf := iinfCb.(*IinfBox)

	if len(iinf.itemsById) != 1 {
		t.Fatalf("itemsById should have exactly one item.")
	} else if len(iinf.itemsByName) != 1 {
		t.Fatalf("itemsByName should have exactly one item.")
	}

	recoveredInfe, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe != infe {
		t.Fatalf("GetItemWithName() does not return the correct INFE record.")
	}
}

func TestInfeBox_New_Version1(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	// Construct INFE data (INFE boxes are children of IINF boxes).

	var infeData []byte

	version := uint8(1)
	bmfcommon.PushBytes(&infeData, []byte{version, 0, 0, 0})

	itemId := uint16(11)
	bmfcommon.PushBytes(&infeData, itemId)

	itemProtectionIndex := uint16(0)
	bmfcommon.PushBytes(&infeData, itemProtectionIndex)

	itemNameRaw := []byte("abc\000")
	bmfcommon.PushBytes(&infeData, itemNameRaw)

	contentTypeRaw := []byte("def\000")
	bmfcommon.PushBytes(&infeData, contentTypeRaw)

	contentEncodingRaw := []byte("ghi\000")
	bmfcommon.PushBytes(&infeData, contentEncodingRaw)

	extensionType := uint32(22)
	bmfcommon.PushBytes(&infeData, extensionType)

	// Construct an IINF.

	var iinfData []byte

	// Version.
	iinfVersion := uint8(0)
	bmfcommon.PushBytes(&iinfData, []byte{iinfVersion, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&iinfData, uint16(11))

	// Embed child INFE box.
	bmfcommon.PushBox(&iinfData, "infe", infeData)

	var b []byte
	bmfcommon.PushBox(&b, "iinf", iinfData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	fbi := resource.Index()

	if len(fbi) != 2 {
		t.Fatalf("Incorrect number of indexed boxes.")
	}

	// Validate INFE.

	infeCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf.infe", 0}]
	if found != true {
		t.Fatalf("INFE not indexed.")
	}

	infe := infeCb.(*InfeBox)

	if infe.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if infe.itemProtectionIndex != 0 {
		t.Fatalf("itemProtectionIndex not correct.")
	} else if infe.itemName != "abc" {
		t.Fatalf("itemName not correct.")
	} else if infe.contentType != "def" {
		t.Fatalf("contentType not correct.")
	} else if infe.contentEncoding != "ghi" {
		t.Fatalf("contentEncoding not correct.")
	} else if infe.extensionType != 22 {
		t.Fatalf("extensionType not correct.")
	}

	// Validate IINF (that the INFE is registered in it).

	iinfCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf", 0}]
	if found != true {
		t.Fatalf("IINF not indexed.")
	}

	iinf := iinfCb.(*IinfBox)

	if len(iinf.itemsById) != 1 {
		t.Fatalf("itemsById should have exactly one item.")
	} else if len(iinf.itemsByName) != 1 {
		t.Fatalf("itemsByName should have exactly one item.")
	}

	recoveredInfe, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe != infe {
		t.Fatalf("GetItemWithName() does not return the correct INFE record.")
	}
}

func TestInfeBox_New_Version3(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	// Construct INFE data (INFE boxes are children of IINF boxes).

	var infeData []byte

	infeVersion := uint8(3)
	bmfcommon.PushBytes(&infeData, []byte{infeVersion, 0, 0, 0})

	itemId := uint32(11)
	bmfcommon.PushBytes(&infeData, itemId)

	itemProtectionIndex := uint16(0)
	bmfcommon.PushBytes(&infeData, itemProtectionIndex)

	itemType := InfeItemTypeFromBytes([4]byte{'a', 'b', 'c', 'd'})
	bmfcommon.PushBytes(&infeData, uint32(itemType))

	itemNameRaw := []byte("abc\000")
	bmfcommon.PushBytes(&infeData, itemNameRaw)

	// Construct an IINF.

	var iinfData []byte

	// Version.
	iinfVersion := uint8(0)
	bmfcommon.PushBytes(&iinfData, []byte{iinfVersion, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&iinfData, uint16(11))

	// Embed child INFE box.
	bmfcommon.PushBox(&iinfData, "infe", infeData)

	var b []byte
	bmfcommon.PushBox(&b, "iinf", iinfData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	fbi := resource.Index()

	if len(fbi) != 2 {
		t.Fatalf("Incorrect number of indexed boxes.")
	}

	// Validate INFE.

	infeCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf.infe", 0}]
	if found != true {
		t.Fatalf("INFE not indexed.")
	}

	infe := infeCb.(*InfeBox)

	if infe.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if infe.itemProtectionIndex != 0 {
		t.Fatalf("itemProtectionIndex not correct.")
	} else if infe.itemName != "abc" {
		t.Fatalf("itemName not correct.")
	}

	// Validate IINF (that the INFE is registered in it).

	iinfCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf", 0}]
	if found != true {
		t.Fatalf("IINF not indexed.")
	}

	iinf := iinfCb.(*IinfBox)

	if len(iinf.itemsById) != 1 {
		t.Fatalf("itemsById should have exactly one item.")
	} else if len(iinf.itemsByName) != 1 {
		t.Fatalf("itemsByName should have exactly one item.")
	}

	recoveredInfe, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe != infe {
		t.Fatalf("GetItemWithName() does not return the correct INFE record.")
	}
}

func TestInfeBox_New_Version2_Mime(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	// Construct INFE data (INFE boxes are children of IINF boxes).

	var infeData []byte

	infeVersion := uint8(2)
	bmfcommon.PushBytes(&infeData, []byte{infeVersion, 0, 0, 0})

	itemId := uint16(11)
	bmfcommon.PushBytes(&infeData, itemId)

	itemProtectionIndex := uint16(0)
	bmfcommon.PushBytes(&infeData, itemProtectionIndex)

	itemType := InfeItemTypeFromBytes([4]byte{'m', 'i', 'm', 'e'})
	bmfcommon.PushBytes(&infeData, uint32(itemType))

	itemNameRaw := []byte("abc\000")
	bmfcommon.PushBytes(&infeData, itemNameRaw)

	contentTypeRaw := []byte("def\000")
	bmfcommon.PushBytes(&infeData, contentTypeRaw)

	contentEncodingRaw := []byte("ghi\000")
	bmfcommon.PushBytes(&infeData, contentEncodingRaw)

	// Construct an IINF.

	var iinfData []byte

	// Version.
	iinfVersion := uint8(0)
	bmfcommon.PushBytes(&iinfData, []byte{iinfVersion, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&iinfData, uint16(11))

	// Embed child INFE box.
	bmfcommon.PushBox(&iinfData, "infe", infeData)

	var b []byte
	bmfcommon.PushBox(&b, "iinf", iinfData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	fbi := resource.Index()

	if len(fbi) != 2 {
		t.Fatalf("Incorrect number of indexed boxes.")
	}

	// Validate INFE.

	infeCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf.infe", 0}]
	if found != true {
		t.Fatalf("INFE not indexed.")
	}

	infe := infeCb.(*InfeBox)

	if infe.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if infe.itemProtectionIndex != 0 {
		t.Fatalf("itemProtectionIndex not correct.")
	} else if infe.itemName != "abc" {
		t.Fatalf("itemName not correct.")
	} else if infe.contentType != "def" {
		t.Fatalf("contentType not correct.")
	} else if infe.contentEncoding != "ghi" {
		t.Fatalf("contentEncoding not correct.")
	}

	// Validate IINF (that the INFE is registered in it).

	iinfCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf", 0}]
	if found != true {
		t.Fatalf("IINF not indexed.")
	}

	iinf := iinfCb.(*IinfBox)

	if len(iinf.itemsById) != 1 {
		t.Fatalf("itemsById should have exactly one item.")
	} else if len(iinf.itemsByName) != 1 {
		t.Fatalf("itemsByName should have exactly one item.")
	}

	recoveredInfe, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe != infe {
		t.Fatalf("GetItemWithName() does not return the correct INFE record.")
	}
}

func TestInfeBox_New_Version2_Uri(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	// Construct INFE data (INFE boxes are children of IINF boxes).

	var infeData []byte

	infeVersion := uint8(2)
	bmfcommon.PushBytes(&infeData, []byte{infeVersion, 0, 0, 0})

	itemId := uint16(11)
	bmfcommon.PushBytes(&infeData, itemId)

	itemProtectionIndex := uint16(0)
	bmfcommon.PushBytes(&infeData, itemProtectionIndex)

	itemType := InfeItemTypeFromBytes([4]byte{'u', 'r', 'i', ' '})
	bmfcommon.PushBytes(&infeData, uint32(itemType))

	itemNameRaw := []byte("abc\000")
	bmfcommon.PushBytes(&infeData, itemNameRaw)

	itemUriTypeRaw := []byte("def\000")
	bmfcommon.PushBytes(&infeData, itemUriTypeRaw)

	// Construct an IINF.

	var iinfData []byte

	// Version.
	iinfVersion := uint8(0)
	bmfcommon.PushBytes(&iinfData, []byte{iinfVersion, 0, 0, 0})

	// Entry count.
	bmfcommon.PushBytes(&iinfData, uint16(11))

	// Embed child INFE box.
	bmfcommon.PushBox(&iinfData, "infe", infeData)

	var b []byte
	bmfcommon.PushBox(&b, "iinf", iinfData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	fbi := resource.Index()

	if len(fbi) != 2 {
		t.Fatalf("Incorrect number of indexed boxes.")
	}

	// Validate INFE.

	infeCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf.infe", 0}]
	if found != true {
		t.Fatalf("INFE not indexed.")
	}

	infe := infeCb.(*InfeBox)

	if infe.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if infe.itemProtectionIndex != 0 {
		t.Fatalf("itemProtectionIndex not correct.")
	} else if infe.itemName != "abc" {
		t.Fatalf("itemName not correct.")
	} else if infe.itemUriType != "def" {
		t.Fatalf("itemUriType not correct.")
	}

	// Validate IINF (that the INFE is registered in it).

	iinfCb, found := fbi[bmfcommon.IndexedBoxEntry{"iinf", 0}]
	if found != true {
		t.Fatalf("IINF not indexed.")
	}

	iinf := iinfCb.(*IinfBox)

	if len(iinf.itemsById) != 1 {
		t.Fatalf("itemsById should have exactly one item.")
	} else if len(iinf.itemsByName) != 1 {
		t.Fatalf("itemsByName should have exactly one item.")
	}

	recoveredInfe, err := iinf.GetItemWithName("abc")
	log.PanicIf(err)

	if recoveredInfe != infe {
		t.Fatalf("GetItemWithName() does not return the correct INFE record.")
	}
}
