package bmftype

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"encoding/binary"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

var (
	ilocLogger = log.NewLogger("bmftype.iloc")
)

var (
	// ErrLocationItemNotFound indicates that no location record could be found
	// for the given item-ID.
	ErrLocationItemNotFound = errors.New("no location record for item")
)

// IlocIntegerWidth describes how many bytes an integer will be.
type IlocIntegerWidth uint8

// IsValid returns whether the IIW describes a valid number of bytes.
func (iiw IlocIntegerWidth) IsValid() bool {
	return iiw == 0 || iiw == 4 || iiw == 8
}

// IlocBox is the "Item Location" box.
type IlocBox struct {
	bmfcommon.Box

	// TODO(dustin): Finish adding accessors

	version byte

	offsetSize     IlocIntegerWidth
	lengthSize     IlocIntegerWidth
	baseOffsetSize IlocIntegerWidth
	indexSize      IlocIntegerWidth

	items      []IlocItem
	itemsIndex map[uint32]IlocItem
}

// GetWithId returns the location record for the item with the given ID.
func (ib IlocBox) GetWithId(itemId uint32) (ii IlocItem, err error) {
	ii, found := ib.itemsIndex[itemId]
	if found == false {
		return ii, ErrLocationItemNotFound
	}

	return ii, nil
}

// sortedItemIds returns a list of sorted item-IDs. Note that, in order to
// simply do this, they are converted from []uint32 to []int. There is no risk
// of overflow.
func (iloc IlocBox) sortedItemIds() []int {
	itemIds := make([]int, len(iloc.itemsIndex))

	i := 0
	for key := range iloc.itemsIndex {
		itemIds[i] = int(key)
		i++
	}

	sort.Ints(itemIds)

	return itemIds
}

// Dump prints the item map and the extent info for each.
func (iloc IlocBox) Dump() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	itemIds := iloc.sortedItemIds()
	fbi := iloc.Index()
	for _, itemId := range itemIds {
		// Get item information provided by its reference record in this ILOC
		// box.

		ii, err := iloc.GetWithId(uint32(itemId))
		log.PanicIf(err)

		// Get IINF box so that we can get the INFE box for this item (so we can
		// get its type).

		iinfCommonBox, found := fbi[bmfcommon.IndexedBoxEntry{"meta.iinf", 0}]
		if found == false {
			log.Panicf("Could not find IINF box (ILOC/Dump)")
		}

		iinf := iinfCommonBox.(*IinfBox)

		// Get the INFE box for this item.

		infe, err := iinf.GetItemWithId(uint32(itemId))
		log.PanicIf(err)

		fmt.Printf("%s TYPE=[%s]\n", ii.InlineString(), infe.ItemType().String())

		for _, ie := range ii.Extents() {
			fmt.Printf("- %s\n", ie.InlineString())
		}

		fmt.Printf("\n")
	}

	return nil
}

func (iloc IlocBox) writeItemExtent(itemId int, infe *InfeBox, ie IlocExtent, extentNumber int, outPath string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	infePhrase := infe.ItemType().String()

	offset := int64(ie.Offset())
	length := int64(ie.Length())

	filename := fmt.Sprintf("extent.%d.%d.%s", itemId, extentNumber, infePhrase)
	filepath := path.Join(outPath, filename)

	fmt.Printf("Writing [%s] (%d bytes).\n", filepath, ie.Length())

	f, err := os.Create(filepath)
	log.PanicIf(err)

	defer f.Close()

	err = iloc.CopyBytesAt(offset, length, f)
	log.PanicIf(err)

	return nil
}

func (iloc IlocBox) writeItemExtents(itemId int, outPath string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// Get item information provided by its reference record in this ILOC
	// box.

	ii, err := iloc.GetWithId(uint32(itemId))
	log.PanicIf(err)

	// Get IINF box so that we can get the INFE box for this item (so we can
	// get its type).

	fbi := iloc.Index()

	iinfCommonBox, found := fbi[bmfcommon.IndexedBoxEntry{"meta.iinf", 0}]
	if found == false {
		log.Panicf("Could not find IINF box (ILOC/writeItemExtents)")
	}

	iinf := iinfCommonBox.(*IinfBox)

	// Get the INFE box for this item.

	infe, err := iinf.GetItemWithId(uint32(itemId))
	log.PanicIf(err)

	for i, ie := range ii.Extents() {
		err := iloc.writeItemExtent(itemId, infe, ie, i, outPath)
		log.PanicIf(err)
	}

	return nil
}

// WriteExtents writes each of the extents out to files for surgical debugging.
func (iloc IlocBox) WriteExtents(outPath string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	itemIds := iloc.sortedItemIds()
	for _, itemId := range itemIds {
		err := iloc.writeItemExtents(itemId, outPath)
		log.PanicIf(err)
	}

	return nil
}

// InlineString returns an undecorated string of field names and values.
func (iloc *IlocBox) InlineString() string {
	return fmt.Sprintf(
		"%s OFFSET-SIZE=(%d) LENGTH-SIZE=(%d) BASE-OFFSET-SIZE=(%d) INDEX-SIZE=(%d) ITEMS=(%d)",
		iloc.Box.InlineString(), iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize, len(iloc.items))
}

type ilocBoxFactory struct {
}

// Name returns the name of the type.
func (ilocBoxFactory) Name() string {
	return "iloc"
}

// IlocExtent describes a single ILOC extent.
type IlocExtent struct {
	extentIndex  uint64
	extentOffset uint64
	extentLength uint64
}

// Index returns the extent index.
func (ie IlocExtent) Index() uint64 {
	return ie.extentIndex
}

// Offset returns the offset of the extent.
func (ie IlocExtent) Offset() uint64 {
	return ie.extentOffset
}

// Length returns the length of the extent.
func (ie IlocExtent) Length() uint64 {
	return ie.extentLength
}

// InlineString returns an undecorated string of field names and values.
func (ie IlocExtent) InlineString() string {
	return fmt.Sprintf(
		"OFFSET=(0x%016x) LENGTH=(%d) INDEX=(0x%016x)",
		ie.extentOffset, ie.extentLength, ie.extentIndex)
}

// String returns a stringified description of an iloc extent.
func (ie IlocExtent) String() string {
	return fmt.Sprintf("IlocExtent<%s>", ie.InlineString())
}

// IlocItem is one iloc location item.
type IlocItem struct {
	itemId             uint32
	constructionMethod uint16
	dataReferenceIndex uint16

	// NOTE(dustin): It's not clear what the baseOffset is used for since the extent-offsets seem to already be absolute and sufficient.
	baseOffset []byte

	extents []IlocExtent
}

// Extents returns all extents for item's data.
func (ii IlocItem) Extents() []IlocExtent {
	return ii.extents
}

// InlineString returns an undecorated string of field names and values.
func (ii IlocItem) InlineString() string {
	return fmt.Sprintf("ID=(%d) DATA-REF-INDEX=(%d) BASE-OFFSET=(0x%04x) EXTENT-COUNT=(%d)", ii.itemId, ii.dataReferenceIndex, ii.baseOffset, len(ii.extents))
}

// String returns a stringified description of an iloc item.
func (ii IlocItem) String() string {
	return fmt.Sprintf("IlocItem<%s>", ii.InlineString())
}

// readExtent reads one entry for the current item. One item potentially has many
// entries.
func (factory ilocBoxFactory) readExtent(r io.Reader, version byte, offsetSize, lengthSize, indexSize IlocIntegerWidth) (ie IlocExtent, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if (version == 1 || version == 2) && indexSize > 0 {
		if indexSize == 4 {
			var extentIndex uint32

			err := binary.Read(r, bmfcommon.DefaultEndianness, &extentIndex)
			log.PanicIf(err)

			ie.extentIndex = uint64(extentIndex)
		} else if indexSize == 8 {
			err := binary.Read(r, bmfcommon.DefaultEndianness, &ie.extentIndex)
			log.PanicIf(err)
		}
	}

	if offsetSize == 4 {
		var extentOffset uint32

		err := binary.Read(r, bmfcommon.DefaultEndianness, &extentOffset)
		log.PanicIf(err)

		ie.extentOffset = uint64(extentOffset)
	} else if offsetSize == 8 {
		err := binary.Read(r, bmfcommon.DefaultEndianness, &ie.extentOffset)
		log.PanicIf(err)
	}

	if lengthSize == 4 {
		var extentLength uint32

		err := binary.Read(r, bmfcommon.DefaultEndianness, &extentLength)
		log.PanicIf(err)

		ie.extentLength = uint64(extentLength)
	} else if lengthSize == 8 {
		err := binary.Read(r, bmfcommon.DefaultEndianness, &ie.extentLength)
		log.PanicIf(err)
	}

	return ie, nil
}

// readItem parses one item out of the stream. These belong to the ILOC record.
func (factory ilocBoxFactory) readItem(r io.Reader, version byte, baseOffsetSize IlocIntegerWidth, offsetSize, lengthSize, indexSize IlocIntegerWidth) (ii IlocItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// itemId

	if version < 2 {
		var itemId16 uint16

		err := binary.Read(r, bmfcommon.DefaultEndianness, &itemId16)
		log.PanicIf(err)

		ii.itemId = uint32(itemId16)
	} else if version == 2 {
		err := binary.Read(r, bmfcommon.DefaultEndianness, &ii.itemId)
		log.PanicIf(err)
	} else {
		log.Panicf("version (%d) not supported", version)
	}

	// constructionMethod

	if version == 0 {
		// It should *already be* (0), but we're choosing to do something here
		// rather than nothing.
		ii.constructionMethod = 0
	} else if version == 1 || version == 2 {
		err = binary.Read(r, bmfcommon.DefaultEndianness, &ii.constructionMethod)
		log.PanicIf(err)
	} else {
		log.Panicf("version (%d) not supported", version)
	}

	// dataReferenceIndex

	err = binary.Read(r, bmfcommon.DefaultEndianness, &ii.dataReferenceIndex)
	log.PanicIf(err)

	// baseOffset

	ii.baseOffset = make([]byte, baseOffsetSize)

	_, err = io.ReadFull(r, ii.baseOffset)
	log.PanicIf(err)

	// extentCount

	var extentCount uint16

	err = binary.Read(r, bmfcommon.DefaultEndianness, &extentCount)
	log.PanicIf(err)

	// Load the extents.

	ii.extents = make([]IlocExtent, int(extentCount))

	for j := 0; j < int(extentCount); j++ {
		ie, err := factory.readExtent(r, version, offsetSize, lengthSize, indexSize)
		log.PanicIf(err)

		ii.extents[j] = ie
	}

	return ii, nil
}

// New returns a new value instance.
func (factory ilocBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	if version > 2 {
		log.Panicf("version of ILOC not supported: (%d)", version)
	}

	b := bytes.NewBuffer(data[4:])
	br := bufio.NewReader(b)

	var packedSize1 IlocIntegerWidth
	err = binary.Read(br, bmfcommon.DefaultEndianness, &packedSize1)
	log.PanicIf(err)

	offsetSize := packedSize1 >> 4

	if offsetSize.IsValid() == false {
		log.Panicf("offset-size is not valid: (%d)", offsetSize)
	}

	lengthSize := (packedSize1 & 0x0f)

	if lengthSize.IsValid() == false {
		log.Panicf("length-size is not valid: (%d)", lengthSize)
	}

	var packedSize2 IlocIntegerWidth
	err = binary.Read(br, bmfcommon.DefaultEndianness, &packedSize2)
	log.PanicIf(err)

	baseOffsetSize := packedSize2 >> 4

	if baseOffsetSize.IsValid() == false {
		log.Panicf("base-offset-size is not valid: (%d)", baseOffsetSize)
	}

	var indexSize IlocIntegerWidth
	if version == 1 || version == 2 {
		indexSize = (packedSize2 & 0x0f)

		if indexSize.IsValid() == false {
			log.Panicf("index-size is not valid: (%d)", indexSize)
		}
	}

	var itemCount uint32

	if version < 2 {
		var itemCount16 uint16
		err := binary.Read(br, bmfcommon.DefaultEndianness, &itemCount16)
		log.PanicIf(err)

		itemCount = uint32(itemCount16)
	} else {
		err := binary.Read(br, bmfcommon.DefaultEndianness, &itemCount)
		log.PanicIf(err)
	}

	itemsIndex := make(map[uint32]IlocItem)

	items := make([]IlocItem, int(itemCount))
	for i := 0; i < int(itemCount); i++ {
		ii, err := factory.readItem(br, version, baseOffsetSize, offsetSize, lengthSize, indexSize)
		log.PanicIf(err)

		items[i] = ii
		itemsIndex[ii.itemId] = ii
	}

	iloc := &IlocBox{
		Box:     box,
		version: version,

		offsetSize:     offsetSize,
		lengthSize:     lengthSize,
		baseOffsetSize: baseOffsetSize,

		indexSize: indexSize,

		items:      items,
		itemsIndex: itemsIndex,
	}

	return iloc, -1, nil
}

var (
	_ bmfcommon.BoxFactory = ilocBoxFactory{}
	_ bmfcommon.CommonBox  = &IlocBox{}
)

func init() {
	bmfcommon.RegisterBoxType(ilocBoxFactory{})
}
