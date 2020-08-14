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
		11: IlocItem{},
		22: IlocItem{},
		33: IlocItem{},
		44: IlocItem{},
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
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
	}

	iinf.loadItem(infe22)

	// Build ILOC

	extents1 := []IlocExtent{IlocExtent{extentOffset: 11110, extentLength: 11111}}
	extents2 := []IlocExtent{IlocExtent{extentOffset: 22220, extentLength: 22221}}

	itemsIndex := map[uint32]IlocItem{
		11: IlocItem{itemId: 11, extents: extents1},
		22: IlocItem{itemId: 22, extents: extents2},
	}

	resource := bmfcommon.NewBmfResource(nil, 0)

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
	resource := bmfcommon.NewBmfResource(sb, 0)
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

	tempPath, err := ioutil.TempDir("", "")
	log.PanicIf(err)

	defer os.RemoveAll(tempPath)

	itemId := 11
	extentNumber := 0

	err = iloc.writeItemExtent(itemId, infe, ie, extentNumber, tempPath)
	log.PanicIf(err)

	// Confirm that the file exists.

	infePhrase := infe.ItemType().String()

	filename := fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath := path.Join(tempPath, filename)

	recoveredData, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData, b[extentOffset:]) != true {
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
	resource := bmfcommon.NewBmfResource(sb, 0)
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
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
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
	extentNumber := 0
	infePhrase := infe11.ItemType().String()

	filename := fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath := path.Join(tempPath, filename)

	recoveredData1, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData1, b[extentOffset1:extentOffset1+4]) != true {
		t.Fatalf("Extent 1 data not correct.")
	}

	// Test second extent of the first write.

	itemId = 11
	extentNumber = 1
	infePhrase = infe11.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData2, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData2, b[extentOffset2:extentOffset2+4]) != true {
		t.Fatalf("Extent 2 data not correct.")
	}

	err = iloc.writeItemExtents(22, tempPath)
	log.PanicIf(err)

	// Test first extent of the second write.

	itemId = 22
	extentNumber = 0
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData3, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData3, b[extentOffset3:extentOffset3+4]) != true {
		t.Fatalf("Extent 3 data not correct:.")
	}

	// Test second extent of the second write.

	itemId = 22
	extentNumber = 1
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData4, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData4, b[extentOffset4:extentOffset4+4]) != true {
		t.Fatalf("Extent 4 data not correct.")
	}
}

func TestIlocBox_WriteExtents(t *testing.T) {
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
	resource := bmfcommon.NewBmfResource(sb, 0)
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
	}

	iinf.loadItem(infe11)

	infeType22 := InfeItemTypeFromBytes([4]byte{'0', '0', '0', '2'})

	infe22 := &InfeBox{
		itemId:   22,
		itemType: infeType22,
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

	err = iloc.WriteExtents(tempPath)
	log.PanicIf(err)

	// Test first extent of the first write.

	itemId := 11
	extentNumber := 0
	infePhrase := infe11.ItemType().String()

	filename := fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath := path.Join(tempPath, filename)

	recoveredData1, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData1, b[extentOffset1:extentOffset1+4]) != true {
		t.Fatalf("Extent 1 data not correct.")
	}

	// Test second extent of the first write.

	itemId = 11
	extentNumber = 1
	infePhrase = infe11.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData2, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData2, b[extentOffset2:extentOffset2+4]) != true {
		t.Fatalf("Extent 2 data not correct.")
	}

	// Test first extent of the second write.

	itemId = 22
	extentNumber = 0
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData3, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData3, b[extentOffset3:extentOffset3+4]) != true {
		t.Fatalf("Extent 3 data not correct:.")
	}

	// Test second extent of the second write.

	itemId = 22
	extentNumber = 1
	infePhrase = infe22.ItemType().String()

	filename = fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath = path.Join(tempPath, filename)

	recoveredData4, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	if bytes.Equal(recoveredData4, b[extentOffset4:extentOffset4+4]) != true {
		t.Fatalf("Extent 4 data not correct.")
	}
}

func TestIlocBox_InlineString(t *testing.T) {
	iloc := &IlocBox{
		offsetSize:     11,
		lengthSize:     22,
		baseOffsetSize: 33,
		indexSize:      44,
		items:          []IlocItem{IlocItem{}, IlocItem{}},
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

		extents: []IlocExtent{IlocExtent{}, IlocExtent{}},
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

		extents: []IlocExtent{IlocExtent{}, IlocExtent{}},
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
