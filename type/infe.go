package bmftype

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"

	"encoding/binary"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// InfeItemType allows simple handling of item-types. It can be compared as a
// uint32 or as a string, adds a few boolean tests, and allows exporting as a
// simple string.
type InfeItemType uint32

// EqualsName encodes the integer back to the original byte-order, convert to a
// string, and compare. The item-type can be interpreted as both an integer and
// a string.
func (iit InfeItemType) EqualsName(s string) bool {

	// TODO(dustin): Add test

	return s == iit.string()
}

func (iit InfeItemType) string() string {

	// TODO(dustin): Add test

	b := make([]byte, 4)
	bmfcommon.DefaultEndianness.PutUint32(b, uint32(iit))

	return string(b)
}

// IsMime returns true if MIME.
func (iit InfeItemType) IsMime() bool {

	// TODO(dustin): Add test

	return iit.EqualsName("mime")
}

// IsUri returns true if URI.
func (iit InfeItemType) IsUri() bool {

	// TODO(dustin): Add test

	return iit.EqualsName("uri ")
}

// String returns the ASCII equivalent if all printable characters or else the
// hex representation.
func (iit InfeItemType) String() string {

	// TODO(dustin): Add test

	s := iit.string()

	for _, r := range s {
		if unicode.IsPrint(r) == false {
			return fmt.Sprintf("TYPE<0x%08x>", uint32(iit))
		}
	}

	return s
}

// InfeBox is the ItemInfoEntry box.
type InfeBox struct {

	// TODO(dustin): Split this into a composition of different pieces of the appropriate versions.

	// Box is the base inner box.
	bmfcommon.Box

	// TODO(dustin): Finish adding accessors

	version byte

	itemId uint32

	// TODO(dustin): This indicates that some protection might obscure metadata access.
	itemProtectionIndex uint16

	itemName        string
	contentType     string
	contentEncoding string

	extensionType uint32

	itemType InfeItemType

	itemUriType string
}

// ItemId returns the ID of the item.
func (infe *InfeBox) ItemId() uint32 {

	// TODO(dustin): Add test

	return infe.itemId
}

// ItemType returns the item-type.
func (infe *InfeBox) ItemType() InfeItemType {

	// TODO(dustin): Add test

	return infe.itemType
}

// InlineString returns an undecorated string of field names and values.
func (infe *InfeBox) InlineString() string {

	// TODO(dustin): Add test

	var extTypePhrase string

	if infe.version == 1 {
		extTypePhrase = fmt.Sprintf(" EXT-TYPE=(%d)", infe.extensionType)
	}

	var mimePhrase string
	var uriPhrase string

	if infe.version >= 2 {
		if infe.itemType.IsMime() == true {
			mimePhrase = fmt.Sprintf(
				" CONTENT-TYPE=[%s] CONTENT-ENCODING=[%s]",
				infe.contentType, infe.contentEncoding)
		}

		if infe.itemType.IsUri() == true {
			uriPhrase = fmt.Sprintf(
				" URI-TYPE=[%s]",
				infe.itemUriType)
		}
	}

	return fmt.Sprintf(
		"%s VER=(%d) ITEM-ID=(%d) PROTECTION-INDEX=(%d) NAME=[%s] ITEM-TYPE=[%s]%s%s%s",
		infe.Box.InlineString(), infe.version, infe.itemId, infe.itemProtectionIndex, infe.itemName, infe.itemType, extTypePhrase, mimePhrase, uriPhrase)
}

type infeBoxFactory struct {
}

// Name returns the name of the type.
func (infeBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "infe"
}

// New returns a new value instance.
func (infeBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	infe := &InfeBox{
		Box:     box,
		version: data[0],
	}

	if infe.version > 2 {
		// TODO(dustin): !! Circle back
		log.Panicf("versions > 2 are not yet supported: (%d)", infe.version)
	}

	b := bytes.NewBuffer(data[4:])
	br := bufio.NewReader(b)

	if infe.version == 0 || infe.version == 1 {
		// itemId

		var itemId16 uint16

		err = binary.Read(br, bmfcommon.DefaultEndianness, &itemId16)
		log.PanicIf(err)

		infe.itemId = uint32(itemId16)

		// itemProtectionIndex

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemProtectionIndex)
		log.PanicIf(err)

		// itemName

		itemNameRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.itemName = itemNameRaw[:len(itemNameRaw)-1]

		// contentType

		contentTypeRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.contentType = contentTypeRaw[:len(contentTypeRaw)-1]

		// contentEncoding

		contentEncodingRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.contentEncoding = contentEncodingRaw[:len(contentEncodingRaw)-1]

		if infe.version == 1 {
			err := binary.Read(br, bmfcommon.DefaultEndianness, &infe.extensionType)
			log.PanicIf(err)
		}
	}

	if infe.version >= 2 {
		// itemId

		if infe.version == 2 {
			var itemId16 uint16

			err = binary.Read(br, bmfcommon.DefaultEndianness, &itemId16)
			log.PanicIf(err)

			infe.itemId = uint32(itemId16)
		} else if infe.version == 3 {
			err := binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemId)
			log.PanicIf(err)
		} else {
			log.Panicf("version (%d) of INFE not supported", infe.version)
		}

		// itemProtectionIndex

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemProtectionIndex)
		log.PanicIf(err)

		// itemType

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemType)
		log.PanicIf(err)

		// itemName

		itemNameRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.itemName = itemNameRaw[:len(itemNameRaw)-1]

		if infe.itemType.IsMime() == true {
			// contentType

			contentTypeRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.contentType = contentTypeRaw[:len(contentTypeRaw)-1]

			// contentEncoding

			contentEncodingRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.contentEncoding = contentEncodingRaw[:len(contentEncodingRaw)-1]
		} else if infe.itemType.IsUri() == true {
			// itemUriType

			itemUriTypeRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.itemUriType = itemUriTypeRaw[:len(itemUriTypeRaw)-1]
		}
	}

	// Load new struct into the item index owned by the parent IINF struct.

	iinfCommon := box.Parent()

	iinf := iinfCommon.(*IinfBox)
	iinf.loadItem(infe)

	return infe, nil
}

var (
	_ bmfcommon.BoxFactory = infeBoxFactory{}
	_ bmfcommon.CommonBox  = &InfeBox{}
)

func init() {
	bmfcommon.RegisterBoxType(infeBoxFactory{})
}
