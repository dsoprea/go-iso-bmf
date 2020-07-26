package bmftype

import (
	"errors"
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

var (
	// ErrNoItemsFound indicates that no items were found in the metadata with
	// the given name/ID.
	ErrNoItemsFound = errors.New("no items found")
)

// IinfBox is the "Item Info" box.
type IinfBox struct {
	// Box is the base inner box.
	bmfcommon.Box

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex

	entryCount  uint32
	loadedCount int

	itemsById   map[uint32]*InfeBox
	itemsByName map[string][]*InfeBox
}

func (iinf *IinfBox) loadItem(infe *InfeBox) {

	// TODO(dustin): Add test

	// Load by-ID index.

	itemId := infe.ItemId()
	iinf.itemsById[itemId] = infe

	// Load by-name index.

	key := infe.ItemType().String()

	if existing, found := iinf.itemsByName[key]; found == true {
		iinf.itemsByName[key] = append(existing, infe)
	} else {
		iinf.itemsByName[key] = []*InfeBox{infe}
	}

	iinf.loadedCount++
}

// GetItemWithId returns the item with the given ID.
func (iinf *IinfBox) GetItemWithId(itemId uint32) (infe *InfeBox, err error) {

	// TODO(dustin): Add test

	infe, found := iinf.itemsById[itemId]
	if found == false {
		return nil, ErrNoItemsFound
	}

	return infe, nil
}

// GetItemsWithName returns all metadata items with the given name or
// ErrNoItemsFound if none.
func (iinf *IinfBox) GetItemsWithName(typeName string) (collected []*InfeBox, err error) {

	// TODO(dustin): Add test

	collected, found := iinf.itemsByName[typeName]
	if found == false {
		return nil, ErrNoItemsFound
	}

	return collected, nil
}

// InlineString returns an undecorated string of field names and values.
func (iinf *IinfBox) InlineString() string {

	// TODO(dustin): Add test

	return fmt.Sprintf(
		"%s ENTRY-COUNT=(%d) LOADED-TYPES=(%d) LOADED-ITEMS=(%d) INDEXED-TYPES=(%d)",
		iinf.Box.InlineString(), iinf.entryCount, iinf.loadedCount, len(iinf.itemsById), len(iinf.itemsByName))
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (iinf *IinfBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	iinf.LoadedBoxIndex = lbi
}

type iinfBoxFactory struct {
}

// Name returns the name of the type.
func (iinfBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "iinf"
}

// New returns a new value instance.
func (iinfBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, skipBytes int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	itemsById := make(map[uint32]*InfeBox)
	itemsByName := make(map[string][]*InfeBox)

	iinf := &IinfBox{
		Box:         box,
		itemsById:   itemsById,
		itemsByName: itemsByName,
	}

	skipBytes = 4

	if version == 0 {
		size := 2

		entryCount16 := bmfcommon.DefaultEndianness.Uint16(data[skipBytes : skipBytes+size])

		iinf.entryCount = uint32(entryCount16)

		skipBytes += size
	} else {
		size := 4

		iinf.entryCount = bmfcommon.DefaultEndianness.Uint32(data[skipBytes : skipBytes+size])

		skipBytes += size
	}

	return iinf, skipBytes, nil
}

var (
	_ bmfcommon.BoxFactory = iinfBoxFactory{}
	_ bmfcommon.CommonBox  = &IinfBox{}
)

func init() {
	bmfcommon.RegisterBoxType(iinfBoxFactory{})
}
