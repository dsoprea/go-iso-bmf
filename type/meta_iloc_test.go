package bmftype

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestIlocBox_GetWithId(t *testing.T) {
	ii1 := IlocItem{
		itemId: 11,
	}

	ii2 := IlocItem{
		itemId: 22,
	}

	itemsIndex := map[uint32]IlocItem{
		11: ii1,
		22: ii2,
	}

	iloc := &IlocBox{
		itemsIndex: itemsIndex,
	}

	ii1Recovered, err := iloc.GetWithId(11)
	log.PanicIf(err)

	if ii1Recovered.itemId != ii1.itemId {
		t.Fatalf("Item not found (first)")
	}

	ii2Recovered, err := iloc.GetWithId(22)
	log.PanicIf(err)

	if ii2Recovered.itemId != ii2.itemId {
		t.Fatalf("Item not found (second)")
	}

	_, err = iloc.GetWithId(33)
	if err == nil {
		t.Fatalf("Expected error.")
	} else if log.Is(err, ErrLocationItemNotFound) != true {
		log.Panic(err)
	}
}

func TestIlocBox_sortedItemIds(t *testing.T) {
	itemsIndex := map[uint32]IlocItem{
		11: {},
		22: {},
		33: {},
		44: {},
	}

	iloc := &IlocBox{
		itemsIndex: itemsIndex,
	}

	itemIds := iloc.sortedItemIds()

	expectedItemIds := []int{
		11, 22, 33, 44,
	}

	if reflect.DeepEqual(itemIds, expectedItemIds) != true {
		t.Fatalf("Item IDs not correct: %v", itemIds)
	}
}

func TestIlocBox_Dump(t *testing.T) {
	// Build IINF

	iinf := newIinfBox(bmfcommon.Box{})

	infeType11 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '1'})

	infe11 := &InfeBox{
		itemId:   11,
		itemType: infeType11,
		itemName: "abc",
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
		itemName: "def",
	}

	iinf.loadItem(infe22)

	// Build ILOC

	extents1 := []IlocExtent{{extentOffset: 11110, extentLength: 11111}}
	extents2 := []IlocExtent{{extentOffset: 22220, extentLength: 22221}}

	itemsIndex := map[uint32]IlocItem{
		11: {itemId: 11, extents: extents1},
		22: {itemId: 22, extents: extents2},
	}

	resource, err := bmfcommon.NewResource(nil, 0)
	log.PanicIf(err)

	// Load the IINF in the index so that the ILOC can find it.

	fbi := resource.Index()
	fbi[bmfcommon.IndexedBoxEntry{"meta.iinf", 0}] = iinf

	box := bmfcommon.NewBox("", 0, 0, 0, resource)

	iloc := &IlocBox{
		Box:        box,
		itemsIndex: itemsIndex,
	}

	iloc.Dump()
}

func TestIlocBox_writeItemExtent(t *testing.T) {
	b := []byte{
		// Put it in the middle of the stream so that it's easier to spot a bug.
		0, 0, 0, 0, 0, 0, 0, 0,

		// Actual data.
		1, 2, 3, 4,
	}

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	box := bmfcommon.NewBox("", 0, 0, 0, resource)

	iloc := &IlocBox{
		Box: box,
	}

	infeType := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '1'})

	infe := &InfeBox{
		itemType: infeType,
	}

	extentOffset := 8

	ie := IlocExtent{
		extentOffset: uint64(extentOffset),
		extentLength: 4,
	}

	itemId := 11
	extentNumber := 0

	outBuffer := new(bytes.Buffer)

	err = iloc.writeItemExtent(itemId, infe, ie, extentNumber, outBuffer)
	log.PanicIf(err)

	if bytes.Equal(outBuffer.Bytes(), b[extentOffset:]) != true {
		t.Fatalf("Data not correct.")
	}
}

func TestIlocBox_writeItemExtents(t *testing.T) {
	// Establish base box.

	b := []byte{
		// Put it in the middle of the stream so that it's easier to spot a bug.
		0, 0, 0, 0, 0, 0, 0, 0,

		// Actual data.
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	box := bmfcommon.NewBox("", 0, 0, 0, resource)

	// Establish first ILOC item and its extents.

	extentOffset1 := 8

	ie1 := IlocExtent{
		extentOffset: uint64(extentOffset1),
		extentLength: 4,
	}

	extentOffset2 := 12

	ie2 := IlocExtent{
		extentOffset: uint64(extentOffset2),
		extentLength: 4,
	}

	ii11 := IlocItem{
		itemId:  11,
		extents: []IlocExtent{ie1, ie2},
	}

	// Establish second ILOC item and its extents.

	extentOffset3 := 16

	ie3 := IlocExtent{
		extentOffset: uint64(extentOffset3),
		extentLength: 4,
	}

	extentOffset4 := 20

	ie4 := IlocExtent{
		extentOffset: uint64(extentOffset4),
		extentLength: 4,
	}

	ii22 := IlocItem{
		itemId:  22,
		extents: []IlocExtent{ie3, ie4},
	}

	// Establish and register some INFE boxes with an IINF box.

	iinf := newIinfBox(bmfcommon.Box{})

	infeType11 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '1'})

	infe11 := &InfeBox{
		itemId:   11,
		itemType: infeType11,
		itemName: "abc",
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
		itemName: "def",
	}

	iinf.loadItem(infe22)

	// Register the IINF box so it's findable.

	fbi := resource.Index()
	fbi[bmfcommon.IndexedBoxEntry{"meta.iinf", 0}] = iinf

	// Finally, create the ILOC box so that we can do some testing.

	// Setup the items index and the ILOC box.

	itemsIndex := map[uint32]IlocItem{
		11: ii11,
		22: ii22,
	}

	iloc := &IlocBox{
		Box:        box,
		itemsIndex: itemsIndex,
	}

	// Write extents.

	tempPath, err := ioutil.TempDir("", "")
	log.PanicIf(err)

	defer os.RemoveAll(tempPath)

	err = iloc.writeItemExtents(11, tempPath)
	log.PanicIf(err)

	// Test first extent of the first write.

	itemId := 11
	infePhrase := infe11.ItemType().String()

	filename := fmt.Sprintf("data.%d.%s", itemId, infePhrase)
	filepath := path.Join(tempPath, filename)

	recoveredData1, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData1, b[extentOffset1:extentOffset1+8]) != true {
		t.Fatalf("Extent 1 data not correct.")
	}

	err = iloc.writeItemExtents(22, tempPath)
	log.PanicIf(err)

	// Test first extent of the second write.

	itemId = 22
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("data.%d.%s", itemId, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData3, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData3, b[extentOffset3:extentOffset3+8]) != true {
		t.Fatalf("Extent 3 data not correct.")
	}
}

func TestIlocBox_Write(t *testing.T) {
	// Establish base box.

	b := []byte{
		// Put it in the middle of the stream so that it's easier to spot a bug.
		0, 0, 0, 0, 0, 0, 0, 0,

		// Actual data.
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	box := bmfcommon.NewBox("", 0, 0, 0, resource)

	// Establish first ILOC item and its extents.

	extentOffset1 := 8

	ie1 := IlocExtent{
		extentOffset: uint64(extentOffset1),
		extentLength: 4,
	}

	extentOffset2 := 12

	ie2 := IlocExtent{
		extentOffset: uint64(extentOffset2),
		extentLength: 4,
	}

	ii11 := IlocItem{
		itemId:  11,
		extents: []IlocExtent{ie1, ie2},
	}

	// Establish second ILOC item and its extents.

	extentOffset3 := 16

	ie3 := IlocExtent{
		extentOffset: uint64(extentOffset3),
		extentLength: 4,
	}

	extentOffset4 := 20

	ie4 := IlocExtent{
		extentOffset: uint64(extentOffset4),
		extentLength: 4,
	}

	ii22 := IlocItem{
		itemId:  22,
		extents: []IlocExtent{ie3, ie4},
	}

	// Establish and register some INFE boxes with an IINF box.

	iinf := newIinfBox(bmfcommon.Box{})

	infeType11 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '1'})

	infe11 := &InfeBox{
		itemId:   11,
		itemType: infeType11,
		itemName: "abc",
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
		itemName: "def",
	}

	iinf.loadItem(infe22)

	// Register the IINF box so it's findable.

	fbi := resource.Index()
	fbi[bmfcommon.IndexedBoxEntry{"meta.iinf", 0}] = iinf

	// Finally, create the ILOC box so that we can do some testing.

	// Setup the items index and the ILOC box.

	itemsIndex := map[uint32]IlocItem{
		11: ii11,
		22: ii22,
	}

	iloc := &IlocBox{
		Box:        box,
		itemsIndex: itemsIndex,
	}

	// Write extents.

	tempPath, err := ioutil.TempDir("", "")
	log.PanicIf(err)

	defer os.RemoveAll(tempPath)

	err = iloc.Write(tempPath)
	log.PanicIf(err)

	// Test first extent of the first write.

	itemId := 11
	infePhrase := infe11.ItemType().String()

	filename := fmt.Sprintf("data.%d.%s", itemId, infePhrase)
	filepath := path.Join(tempPath, filename)

	recoveredData1, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData1, b[extentOffset1:extentOffset1+8]) != true {
		t.Fatalf("Extent 1 data not correct.")
	}

	// Test first extent of the second write.

	itemId = 22
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("data.%d.%s", itemId, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData3, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData3, b[extentOffset3:extentOffset3+8]) != true {
		t.Fatalf("Extent 3 data not correct.")
	}
}

func TestIlocBox_InlineString(t *testing.T) {
	iloc := &IlocBox{
		offsetSize:     11,
		lengthSize:     22,
		baseOffsetSize: 33,
		indexSize:      44,
		items:          []IlocItem{{}, {}},
	}

	if iloc.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) OFFSET-SIZE=(11) LENGTH-SIZE=(22) BASE-OFFSET-SIZE=(33) INDEX-SIZE=(44) ITEMS=(2)" {
		t.Fatalf("InlineString() not correct: [%s]", iloc.InlineString())
	}
}

func TestIlocBoxFactory_Name(t *testing.T) {
	factory := ilocBoxFactory{}
	if factory.Name() != "iloc" {
		t.Fatalf("Name() not correct.")
	}
}

func TestIlocExtent_Index(t *testing.T) {
	extent := IlocExtent{
		extentIndex: 11,
	}

	if extent.Index() != 11 {
		t.Fatalf("Index() not correct.")
	}
}

func TestIlocExtent_Offset(t *testing.T) {
	extent := IlocExtent{
		extentOffset: 11,
	}

	if extent.Offset() != 11 {
		t.Fatalf("Offset() not correct.")
	}
}

func TestIlocExtent_Length(t *testing.T) {
	extent := IlocExtent{
		extentLength: 11,
	}

	if extent.Length() != 11 {
		t.Fatalf("Length() not correct.")
	}
}

func TestIlocExtent_InlineString(t *testing.T) {
	extent := IlocExtent{
		extentIndex:  11,
		extentOffset: 22,
		extentLength: 33,
	}

	if extent.InlineString() != "OFFSET=(0x0000000000000016) LENGTH=(33) INDEX=(0x000000000000000b)" {
		t.Fatalf("InlineString() not correct: [%s]", extent.InlineString())
	}
}

func TestIlocExtent_String(t *testing.T) {
	extent := IlocExtent{
		extentIndex:  11,
		extentOffset: 22,
		extentLength: 33,
	}

	if extent.String() != "IlocExtent<OFFSET=(0x0000000000000016) LENGTH=(33) INDEX=(0x000000000000000b)>" {
		t.Fatalf("String() not correct: [%s]", extent.String())
	}
}

func TestIlocItem_Extents(t *testing.T) {
	extent1 := IlocExtent{
		extentOffset: 11,
	}

	extent2 := IlocExtent{
		extentOffset: 22,
	}

	ii := IlocItem{
		extents: []IlocExtent{extent1, extent2},
	}

	extents := ii.Extents()

	if len(extents) != 2 {
		t.Fatalf("Exactly two extents expected.")
	} else if extents[0] != extent1 {
		t.Fatalf("First extent not correct.")
	} else if extents[1] != extent2 {
		t.Fatalf("Second extent not correct.")
	}
}

func TestIlocItem_InlineString(t *testing.T) {
	ii := IlocItem{
		itemId:             11,
		constructionMethod: 22,
		dataReferenceIndex: 33,
		baseOffset:         []byte{1, 2, 3, 4},

		extents: []IlocExtent{{}, {}},
	}

	if ii.InlineString() != "ID=(11) DATA-REF-INDEX=(33) BASE-OFFSET=(0x01020304) EXTENT-COUNT=(2)" {
		t.Fatalf("InlineString() not correct: [%s]", ii.InlineString())
	}
}

func TestIlocItem_String(t *testing.T) {
	ii := IlocItem{
		itemId:             11,
		constructionMethod: 22,
		dataReferenceIndex: 33,
		baseOffset:         []byte{1, 2, 3, 4},

		extents: []IlocExtent{{}, {}},
	}

	if ii.String() != "IlocItem<ID=(11) DATA-REF-INDEX=(33) BASE-OFFSET=(0x01020304) EXTENT-COUNT=(2)>" {
		t.Fatalf("String() not correct: [%s]", ii.String())
	}
}

func TestIlocIntegerWidth_IsValid_0(t *testing.T) {
	if IlocIntegerWidth(0).IsValid() != true {
		t.Fatalf("IsValid() for (0) not true.")
	}
}

func TestIlocIntegerWidth_IsValid_4(t *testing.T) {
	if IlocIntegerWidth(4).IsValid() != true {
		t.Fatalf("IsValid() for (4) not true.")
	}
}

func TestIlocIntegerWidth_IsValid_8(t *testing.T) {
	if IlocIntegerWidth(8).IsValid() != true {
		t.Fatalf("IsValid() for (8) not true.")
	}
}

func TestIlocIntegerWidth_IsValid_3(t *testing.T) {
	if IlocIntegerWidth(3).IsValid() != false {
		t.Fatalf("IsValid() for (3) not false.")
	}
}

func TestIlocBoxFactory_readExtent_version0_uint32(t *testing.T) {
	var data []byte

	extentOffset := uint32(11)
	extentLength := uint32(22)
	writeIlocExtentVersion032bitBytes(&data, extentOffset, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(0)

	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  0,
		extentOffset: uint64(extentOffset),
		extentLength: uint64(extentLength),
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readExtent_version0_uint64(t *testing.T) {
	var data []byte

	extentOffset := uint64(11)
	bmfcommon.PushBytes(&data, extentOffset)

	extentLength := uint64(22)
	bmfcommon.PushBytes(&data, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(0)

	offsetSize := IlocIntegerWidth(8)
	lengthSize := IlocIntegerWidth(8)
	indexSize := IlocIntegerWidth(8)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  0,
		extentOffset: extentOffset,
		extentLength: extentLength,
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readExtent_version1_uint32(t *testing.T) {
	var data []byte

	extentIndex := uint32(33)
	extentOffset := uint32(11)
	extentLength := uint32(22)

	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex, extentOffset, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(1)

	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  uint64(extentIndex),
		extentOffset: uint64(extentOffset),
		extentLength: uint64(extentLength),
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readExtent_version1_uint64(t *testing.T) {
	var data []byte

	extentIndex := uint64(33)
	bmfcommon.PushBytes(&data, extentIndex)

	extentOffset := uint64(11)
	bmfcommon.PushBytes(&data, extentOffset)

	extentLength := uint64(22)
	bmfcommon.PushBytes(&data, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(1)

	offsetSize := IlocIntegerWidth(8)
	lengthSize := IlocIntegerWidth(8)
	indexSize := IlocIntegerWidth(8)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  extentIndex,
		extentOffset: extentOffset,
		extentLength: extentLength,
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readExtent_version2_uint32(t *testing.T) {
	var data []byte

	extentIndex := uint32(33)
	extentOffset := uint32(11)
	extentLength := uint32(22)

	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex, extentOffset, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(1)

	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  uint64(extentIndex),
		extentOffset: uint64(extentOffset),
		extentLength: uint64(extentLength),
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readExtent_version2_uint64(t *testing.T) {
	var data []byte

	extentIndex := uint64(33)
	bmfcommon.PushBytes(&data, extentIndex)

	extentOffset := uint64(11)
	bmfcommon.PushBytes(&data, extentOffset)

	extentLength := uint64(22)
	bmfcommon.PushBytes(&data, extentLength)

	b := bytes.NewBuffer(data)

	version := byte(2)

	offsetSize := IlocIntegerWidth(8)
	lengthSize := IlocIntegerWidth(8)
	indexSize := IlocIntegerWidth(8)

	factory := ilocBoxFactory{}

	ie, err := factory.readExtent(b, version, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	expected := IlocExtent{
		extentIndex:  extentIndex,
		extentOffset: extentOffset,
		extentLength: extentLength,
	}

	if ie != expected {
		t.Fatalf("IlocExtent not correct.")
	}
}

func TestIlocBoxFactory_readItem_Version0_NoExtents(t *testing.T) {
	// Build stream.

	var data []byte

	itemId := uint16(11)
	bmfcommon.PushBytes(&data, itemId)

	dataReferenceIndex := uint16(22)
	bmfcommon.PushBytes(&data, dataReferenceIndex)

	baseOffset := uint32(0x1234)
	bmfcommon.PushBytes(&data, baseOffset)

	extentCount := uint16(0)
	bmfcommon.PushBytes(&data, extentCount)

	b := bytes.NewBuffer(data)

	// Parse.

	factory := ilocBoxFactory{}

	version := byte(0)

	baseOffsetSize := IlocIntegerWidth(4)
	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	ii, err := factory.readItem(b, version, baseOffsetSize, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	if ii.itemId != 11 {
		t.Fatalf("itemId not correct.")
	} else if ii.constructionMethod != 0 {
		t.Fatalf("constructionMethod not correct.")
	} else if ii.dataReferenceIndex != 22 {
		t.Fatalf("dataReferenceIndex not correct.")
	} else if bmfcommon.DefaultEndianness.Uint32(ii.baseOffset) != 0x1234 {
		t.Fatalf("baseOffset not correct.")
	} else if len(ii.extents) != 0 {
		t.Fatalf("There should be no extents.")
	}
}

func TestIlocBoxFactory_readItem_Version0_WithExtents(t *testing.T) {
	// Build stream.

	var data []byte

	itemId := uint16(11)
	dataReferenceIndex := uint16(22)
	baseOffset := uint32(0x1234)

	extents := []IlocExtent{
		{extentOffset: 11, extentLength: 22},
		{extentOffset: 33, extentLength: 44},
	}

	writeIlocItemVersion032bitBytes(&data, itemId, dataReferenceIndex, baseOffset, extents)

	b := bytes.NewBuffer(data)

	// Parse.

	factory := ilocBoxFactory{}

	version := byte(0)

	baseOffsetSize := IlocIntegerWidth(4)
	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	ii, err := factory.readItem(b, version, baseOffsetSize, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	if ii.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if ii.constructionMethod != 0 {
		t.Fatalf("constructionMethod not correct.")
	} else if ii.dataReferenceIndex != dataReferenceIndex {
		t.Fatalf("dataReferenceIndex not correct.")
	} else if bmfcommon.DefaultEndianness.Uint32(ii.baseOffset) != baseOffset {
		t.Fatalf("baseOffset not correct.")
	} else if len(ii.extents) != 2 {
		t.Fatalf("There should be no extents.")
	}

	ie1 := ii.extents[0]
	ie2 := ii.extents[1]

	if ie1.Offset() != extents[0].extentOffset {
		t.Fatalf("First extent offset not correct.")
	} else if ie1.Length() != extents[0].extentLength {
		t.Fatalf("First extent length not correct.")
	} else if ie2.Offset() != extents[1].extentOffset {
		t.Fatalf("Second extent offset not correct.")
	} else if ie2.Length() != extents[1].extentLength {
		t.Fatalf("Second extent length not correct.")
	}
}

func TestIlocBoxFactory_readItem_Version1_WithExtents(t *testing.T) {
	// Build stream.

	var data []byte

	itemId := uint16(11)
	bmfcommon.PushBytes(&data, itemId)

	constructionMethod := uint16(55)
	bmfcommon.PushBytes(&data, constructionMethod)

	dataReferenceIndex := uint16(22)
	bmfcommon.PushBytes(&data, dataReferenceIndex)

	baseOffset := uint32(0x1234)
	bmfcommon.PushBytes(&data, baseOffset)

	extentCount := uint16(2)
	bmfcommon.PushBytes(&data, extentCount)

	extentIndex0 := uint32(55)
	extentOffset0 := uint32(11)
	extentLength0 := uint32(22)
	extentIndex1 := uint32(66)
	extentOffset1 := uint32(33)
	extentLength1 := uint32(44)

	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex0, extentOffset0, extentLength0)
	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex1, extentOffset1, extentLength1)

	b := bytes.NewBuffer(data)

	// Parse.

	factory := ilocBoxFactory{}

	version := byte(1)

	baseOffsetSize := IlocIntegerWidth(4)
	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	ii, err := factory.readItem(b, version, baseOffsetSize, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	if ii.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if ii.constructionMethod != constructionMethod {
		t.Fatalf("constructionMethod not correct.")
	} else if ii.dataReferenceIndex != dataReferenceIndex {
		t.Fatalf("dataReferenceIndex not correct.")
	} else if bmfcommon.DefaultEndianness.Uint32(ii.baseOffset) != baseOffset {
		t.Fatalf("baseOffset not correct.")
	} else if len(ii.extents) != 2 {
		t.Fatalf("There should be no extents.")
	}

	ie1 := ii.extents[0]
	ie2 := ii.extents[1]

	if ie1.Index() != uint64(extentIndex0) {
		t.Fatalf("First extent index not correct.")
	} else if ie1.Offset() != uint64(extentOffset0) {
		t.Fatalf("First extent offset not correct.")
	} else if ie1.Length() != uint64(extentLength0) {
		t.Fatalf("First extent length not correct.")
	} else if ie2.Index() != uint64(extentIndex1) {
		t.Fatalf("Second extent index not correct.")
	} else if ie2.Offset() != uint64(extentOffset1) {
		t.Fatalf("Second extent offset not correct.")
	} else if ie2.Length() != uint64(extentLength1) {
		t.Fatalf("Second extent length not correct.")
	}
}

func TestIlocBoxFactory_readItem_Version2_WithExtents(t *testing.T) {
	// Build stream.

	var data []byte

	itemId := uint32(11)
	bmfcommon.PushBytes(&data, itemId)

	constructionMethod := uint16(55)
	bmfcommon.PushBytes(&data, constructionMethod)

	dataReferenceIndex := uint16(22)
	bmfcommon.PushBytes(&data, dataReferenceIndex)

	baseOffset := uint32(0x1234)
	bmfcommon.PushBytes(&data, baseOffset)

	extentCount := uint16(2)
	bmfcommon.PushBytes(&data, extentCount)

	extentIndex0 := uint32(55)
	extentOffset0 := uint32(11)
	extentLength0 := uint32(22)
	extentIndex1 := uint32(66)
	extentOffset1 := uint32(33)
	extentLength1 := uint32(44)

	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex0, extentOffset0, extentLength0)
	writeIlocExtentVersion1OrVersion232bitBytes(&data, extentIndex1, extentOffset1, extentLength1)

	b := bytes.NewBuffer(data)

	// Parse.

	factory := ilocBoxFactory{}

	version := byte(2)

	baseOffsetSize := IlocIntegerWidth(4)
	offsetSize := IlocIntegerWidth(4)
	lengthSize := IlocIntegerWidth(4)
	indexSize := IlocIntegerWidth(4)

	ii, err := factory.readItem(b, version, baseOffsetSize, offsetSize, lengthSize, indexSize)
	log.PanicIf(err)

	if ii.itemId != uint32(itemId) {
		t.Fatalf("itemId not correct.")
	} else if ii.constructionMethod != constructionMethod {
		t.Fatalf("constructionMethod not correct.")
	} else if ii.dataReferenceIndex != dataReferenceIndex {
		t.Fatalf("dataReferenceIndex not correct.")
	} else if bmfcommon.DefaultEndianness.Uint32(ii.baseOffset) != baseOffset {
		t.Fatalf("baseOffset not correct.")
	} else if len(ii.extents) != 2 {
		t.Fatalf("There should be no extents.")
	}

	ie1 := ii.extents[0]
	ie2 := ii.extents[1]

	if ie1.Index() != uint64(extentIndex0) {
		t.Fatalf("First extent index not correct.")
	} else if ie1.Offset() != uint64(extentOffset0) {
		t.Fatalf("First extent offset not correct.")
	} else if ie1.Length() != uint64(extentLength0) {
		t.Fatalf("First extent length not correct.")
	} else if ie2.Index() != uint64(extentIndex1) {
		t.Fatalf("Second extent index not correct.")
	} else if ie2.Offset() != uint64(extentOffset1) {
		t.Fatalf("Second extent offset not correct.")
	} else if ie2.Length() != uint64(extentLength1) {
		t.Fatalf("Second extent length not correct.")
	}
}

func TestIlocBoxFactory_New_Version0_NoItems(t *testing.T) {
	// Build stream.

	var data []byte

	version := byte(0)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	offsetSize := 4
	lengthSize := 4
	baseOffsetSize := 4

	packedSize1 := (uint8(offsetSize) << 4) + uint8(lengthSize)
	bmfcommon.PushBytes(&data, packedSize1)

	packedSize2 := (uint8(baseOffsetSize) << 4)
	bmfcommon.PushBytes(&data, packedSize2)

	itemCount := uint16(0)
	bmfcommon.PushBytes(&data, itemCount)

	var b []byte
	bmfcommon.PushBox(&b, "iloc", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := ilocBoxFactory{}.New(box)
	log.PanicIf(err)

	iloc := cb.(*IlocBox)

	if iloc.offsetSize != 4 {
		t.Fatalf("offsetSize not correct.")
	} else if iloc.lengthSize != 4 {
		t.Fatalf("lengthSize not correct.")
	} else if iloc.baseOffsetSize != 4 {
		t.Fatalf("baseOffsetSize not correct.")
	} else if iloc.indexSize != 0 {
		t.Fatalf("indexSize not correct.")
	} else if len(iloc.items) != 0 {
		t.Fatalf("items should be empty.")
	} else if len(iloc.itemsIndex) != 0 {
		t.Fatalf("itemsIndex should be empty.")
	}
}

func TestIlocBoxFactory_New_Version0_WithItems(t *testing.T) {
	// Build stream.

	var data []byte

	version := byte(0)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	offsetSize := 4
	lengthSize := 4
	baseOffsetSize := 4

	packedSize1 := (uint8(offsetSize) << 4) + uint8(lengthSize)
	bmfcommon.PushBytes(&data, packedSize1)

	packedSize2 := (uint8(baseOffsetSize) << 4)
	bmfcommon.PushBytes(&data, packedSize2)

	itemCount := uint16(2)
	bmfcommon.PushBytes(&data, itemCount)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId := uint16(11)
	dataReferenceIndex := uint16(22)
	baseOffset := uint32(0x1234)

	extents11 := []IlocExtent{
		{extentOffset: 11, extentLength: 22},
		{extentOffset: 33, extentLength: 44},
	}

	writeIlocItemVersion032bitBytes(&data, itemId, dataReferenceIndex, baseOffset, extents11)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId = uint16(22)
	dataReferenceIndex = uint16(33)
	baseOffset = uint32(0x5678)

	extents22 := []IlocExtent{
		{extentOffset: 55, extentLength: 66},
		{extentOffset: 77, extentLength: 88},
	}

	writeIlocItemVersion032bitBytes(&data, itemId, dataReferenceIndex, baseOffset, extents22)

	// Push to stream.

	var b []byte
	bmfcommon.PushBox(&b, "iloc", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := ilocBoxFactory{}.New(box)
	log.PanicIf(err)

	iloc := cb.(*IlocBox)

	if iloc.offsetSize != 4 {
		t.Fatalf("offsetSize not correct.")
	} else if iloc.lengthSize != 4 {
		t.Fatalf("lengthSize not correct.")
	} else if iloc.baseOffsetSize != 4 {
		t.Fatalf("baseOffsetSize not correct.")
	} else if iloc.indexSize != 0 {
		t.Fatalf("indexSize not correct.")
	} else if len(iloc.items) != 2 {
		t.Fatalf("items should be empty.")
	} else if len(iloc.itemsIndex) != 2 {
		t.Fatalf("itemsIndex should be empty.")
	}

	ii11, err := iloc.GetWithId(11)
	log.PanicIf(err)

	recoveredExtents11 := ii11.Extents()

	if len(recoveredExtents11) != 2 {
		t.Fatalf("Expected two extents in first item.")
	} else if reflect.DeepEqual(recoveredExtents11, extents11) != true {
		t.Fatalf("Two extents in first item are not correct.")
	}

	ii22, err := iloc.GetWithId(22)
	log.PanicIf(err)

	recoveredExtents22 := ii22.Extents()

	if len(recoveredExtents22) != 2 {
		t.Fatalf("Expected two extents in second item.")
	} else if reflect.DeepEqual(recoveredExtents22, extents22) != true {
		t.Fatalf("Two extents in second item are not correct.")
	}
}

func TestIlocBoxFactory_New_Version1_WithItems(t *testing.T) {
	// Build stream.

	var data []byte

	version := byte(1)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	offsetSize := 4
	lengthSize := 4
	baseOffsetSize := 4
	indexSize := 4

	packedSize1 := (uint8(offsetSize) << 4) + uint8(lengthSize)
	bmfcommon.PushBytes(&data, packedSize1)

	packedSize2 := (uint8(baseOffsetSize) << 4) + uint8(indexSize)
	bmfcommon.PushBytes(&data, packedSize2)

	itemCount := uint16(2)
	bmfcommon.PushBytes(&data, itemCount)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId := uint16(11)
	constructionMethod := uint16(33)
	dataReferenceIndex := uint16(22)
	baseOffset := uint32(0x1234)

	extents11 := []IlocExtent{
		{extentIndex: 55, extentOffset: 11, extentLength: 22},
		{extentIndex: 66, extentOffset: 33, extentLength: 44},
	}

	writeIlocItemVersion132bitBytes(&data, itemId, constructionMethod, dataReferenceIndex, baseOffset, extents11)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId = uint16(22)
	constructionMethod = uint16(33)
	dataReferenceIndex = uint16(33)
	baseOffset = uint32(0x5678)

	extents22 := []IlocExtent{
		{extentOffset: 55, extentLength: 66},
		{extentOffset: 77, extentLength: 88},
	}

	writeIlocItemVersion132bitBytes(&data, itemId, constructionMethod, dataReferenceIndex, baseOffset, extents22)

	// Push to stream.

	var b []byte
	bmfcommon.PushBox(&b, "iloc", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := ilocBoxFactory{}.New(box)
	log.PanicIf(err)

	iloc := cb.(*IlocBox)

	if iloc.offsetSize != 4 {
		t.Fatalf("offsetSize not correct.")
	} else if iloc.lengthSize != 4 {
		t.Fatalf("lengthSize not correct.")
	} else if iloc.baseOffsetSize != 4 {
		t.Fatalf("baseOffsetSize not correct.")
	} else if iloc.indexSize != 4 {
		t.Fatalf("indexSize not correct.")
	} else if len(iloc.items) != 2 {
		t.Fatalf("items should be empty.")
	} else if len(iloc.itemsIndex) != 2 {
		t.Fatalf("itemsIndex should be empty.")
	}

	ii11, err := iloc.GetWithId(11)
	log.PanicIf(err)

	recoveredExtents11 := ii11.Extents()

	if len(recoveredExtents11) != 2 {
		t.Fatalf("Expected two extents in first item.")
	} else if reflect.DeepEqual(recoveredExtents11, extents11) != true {
		t.Fatalf("Two extents in first item are not correct.")
	}

	ii22, err := iloc.GetWithId(22)
	log.PanicIf(err)

	recoveredExtents22 := ii22.Extents()

	if len(recoveredExtents22) != 2 {
		t.Fatalf("Expected two extents in second item.")
	} else if reflect.DeepEqual(recoveredExtents22, extents22) != true {
		t.Fatalf("Two extents in second item are not correct.")
	}
}

func TestIlocBoxFactory_New_Version2_WithItems(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Build stream.

	var data []byte

	version := byte(2)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	offsetSize := 4
	lengthSize := 4
	baseOffsetSize := 4
	indexSize := 4

	packedSize1 := (uint8(offsetSize) << 4) + uint8(lengthSize)
	bmfcommon.PushBytes(&data, packedSize1)

	packedSize2 := (uint8(baseOffsetSize) << 4) + uint8(indexSize)
	bmfcommon.PushBytes(&data, packedSize2)

	itemCount := uint32(2)
	bmfcommon.PushBytes(&data, itemCount)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId := uint32(11)
	constructionMethod := uint16(33)
	dataReferenceIndex := uint16(22)
	baseOffset := uint32(0x1234)

	extents11 := []IlocExtent{
		{extentIndex: 55, extentOffset: 11, extentLength: 22},
		{extentIndex: 66, extentOffset: 33, extentLength: 44},
	}

	writeIlocItemVersion232bitBytes(&data, itemId, constructionMethod, dataReferenceIndex, baseOffset, extents11)

	// Add item 1.

	// TODO(dustin): !! Test these values.
	itemId = uint32(22)
	constructionMethod = uint16(33)
	dataReferenceIndex = uint16(33)
	baseOffset = uint32(0x5678)

	extents22 := []IlocExtent{
		{extentOffset: 55, extentLength: 66},
		{extentOffset: 77, extentLength: 88},
	}

	writeIlocItemVersion232bitBytes(&data, itemId, constructionMethod, dataReferenceIndex, baseOffset, extents22)

	// Push to stream.

	var b []byte
	bmfcommon.PushBox(&b, "iloc", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := ilocBoxFactory{}.New(box)
	log.PanicIf(err)

	iloc := cb.(*IlocBox)

	if iloc.offsetSize != 4 {
		t.Fatalf("offsetSize not correct.")
	} else if iloc.lengthSize != 4 {
		t.Fatalf("lengthSize not correct.")
	} else if iloc.baseOffsetSize != 4 {
		t.Fatalf("baseOffsetSize not correct.")
	} else if iloc.indexSize != 4 {
		t.Fatalf("indexSize not correct.")
	} else if len(iloc.items) != 2 {
		t.Fatalf("items should be empty.")
	} else if len(iloc.itemsIndex) != 2 {
		t.Fatalf("itemsIndex should be empty.")
	}

	ii11, err := iloc.GetWithId(11)
	log.PanicIf(err)

	recoveredExtents11 := ii11.Extents()

	if len(recoveredExtents11) != 2 {
		t.Fatalf("Expected two extents in first item.")
	} else if reflect.DeepEqual(recoveredExtents11, extents11) != true {
		t.Fatalf("Two extents in first item are not correct.")
	}

	ii22, err := iloc.GetWithId(22)
	log.PanicIf(err)

	recoveredExtents22 := ii22.Extents()

	if len(recoveredExtents22) != 2 {
		t.Fatalf("Expected two extents in second item.")
	} else if reflect.DeepEqual(recoveredExtents22, extents22) != true {
		t.Fatalf("Two extents in second item are not correct.")
	}
}
