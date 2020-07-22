package bmfcommon

import (
	"bytes"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestFile_readBytesAt_Front(t *testing.T) {
	data := []byte{
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	file := NewFile(sb, int64(len(data)))

	recovered, err := file.readBytesAt(0, 5)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[:5]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestFile_readBytesAt_Middle(t *testing.T) {
	data := []byte{
		0, 0, 0, 0, 0,
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	file := NewFile(sb, int64(len(data)))

	recovered, err := file.readBytesAt(5, 5)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[5:10]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestFile_readBytesAt_MiddleToEnd(t *testing.T) {
	data := []byte{
		0, 0, 0, 0, 0,
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	file := NewFile(sb, int64(len(data)))

	recovered, err := file.readBytesAt(5, 10)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[5:15]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestFile_readBoxAt_Front(t *testing.T) {
	data := []byte{
		0x1, 0x2, 0x3, 0x4,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	file := NewFile(sb, int64(len(data)))

	box, err := file.readBoxAt(0)
	log.PanicIf(err)

	if box.Size() != int64(0x01020304) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestFile_readBoxAt_Middle(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		0x1, 0x2, 0x3, 0x4,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	file := NewFile(sb, int64(len(data)))

	box, err := file.readBoxAt(4)
	log.PanicIf(err)

	if box.Size() != int64(0x01020304) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

// TODO(dustin): !! Fix this to test without depending on registered types.
//
// func TestReadBoxes(t *testing.T) {
// 	defer func() {
// 		if errRaw := recover(); errRaw != nil {
// 			err := errRaw.(error)
// 			log.PrintError(err)

// 			t.Fatalf("Test failed.")
// 		}
// 	}()

// 	ftypBoxData := []byte{
// 		// 0: majorBrand (4)
// 		'a', 'b', 'c', 'd',

// 		// 4: minorVersion (4)
// 		0x1, 0x2, 0x3, 0x4,

// 		// 8: compatibleBrands (8; we chose to add two brands here).
// 		'e', 'f', 'g', 'h',
// 		'i', 'j', 'k', 'l',
// 	}

// 	flagsBytes := make([]byte, 4)
// 	DefaultEndianness.PutUint32(flagsBytes, 0x01020304)

// 	var hdlrBoxData []byte

// 	// Version and flags.
// 	PushBytes(&hdlrBoxData, uint32(0x11223344))

// 	// Reserved spacing.
// 	PushBytes(&hdlrBoxData, uint32(0))

// 	// Handler name
// 	PushBytes(&hdlrBoxData, []byte{'a', 'b', 'c', 'd'})

// 	// Reserved spacing.
// 	// TODO(dustin): This is probably data that we need to add support for.
// 	PushBytes(&hdlrBoxData, uint32(0))
// 	PushBytes(&hdlrBoxData, uint32(0))
// 	PushBytes(&hdlrBoxData, uint32(0))

// 	// handler name (all remaining)
// 	// TODO(dustin): Update this comment to not be a duplicate.
// 	PushBytes(&hdlrBoxData, []byte{'t', 'e', 's', 't', 'n', 'a', 'm', 'e'})

// 	data := []byte{}
// 	PushBox(&data, "ftyp", ftypBoxData)
// 	PushBox(&data, "hdlr", hdlrBoxData)

// 	sb := rifs.NewSeekableBufferWithBytes(data)

// 	size := int64(len(data))

// 	file := NewFile(sb, size)

// 	boxes, err := readBoxes(file, 0, size)
// 	log.PanicIf(err)

// 	if len(boxes) != 2 {
// 		t.Fatalf("Expected two boxes: (%d)", len(boxes))
// 	}

// 	ftypBox := Box{
// 		name:  "ftyp",
// 		start: 0,
// 		size:  24,
// 		file:  file,
// 	}

// 	ftypData, err := ftypBox.ReadBoxData()
// 	log.PanicIf(err)

// 	DumpBytes(ftypData)

// 	hdlrBox := Box{
// 		name:  "hdlr",
// 		start: 24,
// 		size:  40,
// 		file:  file,
// 	}

// 	hdlrData, err := hdlrBox.ReadBoxData()
// 	log.PanicIf(err)

// 	DumpBytes(hdlrData)
// }
