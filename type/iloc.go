package bmftype

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

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

// IlocBox is the "Item Location" box.
type IlocBox struct {
	bmfcommon.Box

	version byte

	offsetSize     IlocIntegerWidth
	lengthSize     IlocIntegerWidth
	baseOffsetSize IlocIntegerWidth
	indexSize      IlocIntegerWidth

	itemCount uint32

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

// InlineString returns an undecorated string of field names and values.
func (iloc *IlocBox) InlineString() string {

	// TODO(dustin): Add test

	return fmt.Sprintf(
		"%s OFFSET-SIZE=(%d) LENGTH-SIZE=(%d) BASE-OFFSET-SIZE=(%d) INDEX-SIZE=(%d) ITEM-COUNT=(%d)",
		iloc.Box.InlineString(), iloc.offsetSize, iloc.lengthSize, iloc.baseOffsetSize, iloc.indexSize, iloc.itemCount)
}

type ilocBoxFactory struct {
}

// Name returns the name of the type.
func (ilocBoxFactory) Name() string {

	// TODO(dustin): Add test

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

	// TODO(dustin): Add test

	return ie.extentIndex
}

// Offset returns the offset of the extent.
func (ie IlocExtent) Offset() uint64 {

	// TODO(dustin): Add test

	return ie.extentOffset
}

// Length returns the length of the extent.
func (ie IlocExtent) Length() uint64 {

	// TODO(dustin): Add test

	return ie.extentLength
}

// String returns a stringified description of an iloc extent.
func (ie IlocExtent) String() string {

	// TODO(dustin): Add test

	return fmt.Sprintf(
		"IlocExtent<INDEX=(0x%016x) OFFSET=(0x%016x) LENGTH=(0x%016x)>",
		ie.extentIndex, ie.extentOffset, ie.extentLength)
}

// IlocItem is one iloc location item.
type IlocItem struct {
	itemId             uint32
	dataReferenceIndex uint16
	baseOffset         []byte
	extentCount        uint16
	extents            []IlocExtent
}

// Extents returns all extents for item's data.
func (ii IlocItem) Extents() []IlocExtent {

	// TODO(dustin): Add test

	return ii.extents
}

// String returns a stringified description of an iloc item.
func (ii IlocItem) String() string {

	// TODO(dustin): Add test

	return fmt.Sprintf("IlocItem<ID=(%d) DATA-REF-INDEX=(%d) BASE-OFFSET=(0x%04x) EXTENT-COUNT=(%d)>", ii.itemId, ii.dataReferenceIndex, ii.baseOffset, ii.extentCount)
}

// IlocIntegerWidth describes how many bytes an integer will be.
type IlocIntegerWidth uint8

// IsValid returns whether the IIW describes a valid number of bytes.
func (iiw IlocIntegerWidth) IsValid() bool {

	// TODO(dustin): Add test

	return iiw == 0 || iiw == 4 || iiw == 8
}

// New returns a new value instance.
func (ilocBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	ilocLogger.Debugf(nil, "ILOC: VERSION=(%d)", version)

	if version != 1 && version != 2 {
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

	ilocLogger.Debugf(nil,
		"ILOC: OFFSET-SIZE=(%d) LEN-SIZE=(%d) BASE-OFFSET-SIZE=(%d) INDEX-SIZE=(%d)",
		offsetSize, lengthSize, baseOffsetSize, indexSize)

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
		ii := IlocItem{}

		// itemId

		if version < 2 {
			var itemId16 uint16

			err := binary.Read(br, bmfcommon.DefaultEndianness, &itemId16)
			log.PanicIf(err)

			ii.itemId = uint32(itemId16)
		} else if version == 2 {
			err := binary.Read(br, bmfcommon.DefaultEndianness, &ii.itemId)
			log.PanicIf(err)
		} else {
			log.Panicf("version (%d) not supported", version)
		}

		// ilocLogger.Debugf(nil, "[%d/%d] Loading ILOC item with ID (%d).", i+1, itemCount, ii.itemId)

		// constructionMethod

		var constructionMethod uint16

		if version == 1 || version == 2 {
			err = binary.Read(br, bmfcommon.DefaultEndianness, &constructionMethod)
			log.PanicIf(err)
		} else {
			log.Panicf("version (%d) not supported", version)
		}

		// dataReferenceIndex

		err := binary.Read(br, bmfcommon.DefaultEndianness, &ii.dataReferenceIndex)
		log.PanicIf(err)

		// baseOffset

		ii.baseOffset = make([]byte, baseOffsetSize)

		_, err = io.ReadFull(br, ii.baseOffset)
		log.PanicIf(err)

		// extentCount

		err = binary.Read(br, bmfcommon.DefaultEndianness, &ii.extentCount)
		log.PanicIf(err)

		// Load the extents.

		ii.extents = make([]IlocExtent, int(ii.extentCount))

		for j := 0; j < int(ii.extentCount); j++ {
			// ilocLogger.Debugf(nil, "[%d/%d] Loading ILOC item (ID (%d)) extent (%d)/(%d).", i+1, itemCount, ii.itemId, j+1, ii.extentCount)

			ie := IlocExtent{}

			if (version == 1 || version == 2) && indexSize > 0 {
				if indexSize == 4 {
					var extentIndex uint32

					err := binary.Read(br, bmfcommon.DefaultEndianness, &extentIndex)
					log.PanicIf(err)

					ie.extentIndex = uint64(extentIndex)
				} else if indexSize == 8 {
					err := binary.Read(br, bmfcommon.DefaultEndianness, &ie.extentIndex)
					log.PanicIf(err)
				}
			}

			if offsetSize == 4 {
				var extentOffset uint32

				err := binary.Read(br, bmfcommon.DefaultEndianness, &extentOffset)
				log.PanicIf(err)

				ie.extentOffset = uint64(extentOffset)
			} else if offsetSize == 8 {
				err := binary.Read(br, bmfcommon.DefaultEndianness, &ie.extentOffset)
				log.PanicIf(err)
			}

			if lengthSize == 4 {
				var extentLength uint32

				err := binary.Read(br, bmfcommon.DefaultEndianness, &extentLength)
				log.PanicIf(err)

				ie.extentLength = uint64(extentLength)
			} else if lengthSize == 8 {
				err := binary.Read(br, bmfcommon.DefaultEndianness, &ie.extentLength)
				log.PanicIf(err)
			}

			ii.extents[j] = ie
		}

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
		itemCount: itemCount,

		items:      items,
		itemsIndex: itemsIndex,
	}

	return iloc, nil
}

var (
	_ bmfcommon.BoxFactory = ilocBoxFactory{}
	_ bmfcommon.CommonBox  = &IlocBox{}
)

func init() {
	bmfcommon.RegisterBoxType(ilocBoxFactory{})
}
