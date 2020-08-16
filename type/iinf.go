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
	ErrNoItemsFound = errors.New("item not found")
)

// IinfBox is the "Item Info" box.
type IinfBox struct {
	// Box is the base inner box.
	bmfcommon.Box

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex

	entryCount uint32

	itemsById   map[uint32]*InfeBox
	itemsByName map[string]*InfeBox
}

func newIinfBox(box bmfcommon.Box) *IinfBox {
	itemsById := make(map[uint32]*InfeBox)
	itemsByName := make(map[string]*InfeBox)

	return &IinfBox{
		Box:         box,
		itemsById:   itemsById,
		itemsByName: itemsByName,
	}
}

// loadItem indexes one item. This is called from the INFE box factory, which
// always occurs after the IINF box.

// Each ID and name must be unique. Since the specification implies and does not
// guarantee that each ID and each name occurs just once, this method has
// assertions that will fail if we were wrong rather than just cover up the
// mistake.
func (iinf *IinfBox) loadItem(infe *InfeBox) {

	// Load by-ID index.

	itemId := infe.ItemId()

	if _, found := iinf.itemsById[itemId]; found == true {
		log.Panicf("item ID (%d) occurs more than once", itemId)
	} else {
		iinf.itemsById[itemId] = infe
	}

	// Load by-name index.

	key := infe.itemName
	if key == "" {
		log.Panicf("INFE item-name is empty.")
	}

	if _, found := iinf.itemsByName[key]; found == true {
		log.Panicf("item with name [%s] occurs more than once", key)
	} else {
		iinf.itemsByName[key] = infe
	}
}

// GetItemWithId returns the item with the given ID.
func (iinf *IinfBox) GetItemWithId(itemId uint32) (infe *InfeBox, err error) {
	infe, found := iinf.itemsById[itemId]
	if found == false {
		return nil, ErrNoItemsFound
	}

	return infe, nil
}

// GetItemWithName returns the item with the given name or
// ErrNoItemsFound if none.
func (iinf *IinfBox) GetItemWithName(typeName string) (infe *InfeBox, err error) {
	infe, found := iinf.itemsByName[typeName]
	if found == false {
		return nil, ErrNoItemsFound
	}

	return infe, nil
}

// InlineString returns an undecorated string of field names and values.
func (iinf *IinfBox) InlineString() string {
	return fmt.Sprintf(
		"%s ENTRY-COUNT=(%d) LOADED-ITEMS=(%d)",
		iinf.Box.InlineString(), iinf.entryCount, len(iinf.itemsById))
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (iinf *IinfBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {
	iinf.LoadedBoxIndex = lbi
}

type iinfBoxFactory struct {
}

// Name returns the name of the type.
func (iinfBoxFactory) Name() string {
	return "iinf"
}

// New returns a new value instance.
func (iinfBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, skipBytes int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	iinf := newIinfBox(box)

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
