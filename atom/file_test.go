package atom

import (
	"bytes"
	"fmt"
	"reflect"
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

	boxSize, boxType, err := file.readBoxAt(0)
	log.PanicIf(err)

	if boxSize != uint32(0x01020304) {
		t.Fatalf("Size not correct: (%d)", boxSize)
	} else if boxType != "abcd" {
		t.Fatalf("Type not correct: [%s]", boxType)
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

	boxSize, boxType, err := file.readBoxAt(4)
	log.PanicIf(err)

	if boxSize != uint32(0x01020304) {
		t.Fatalf("Size not correct: (%d)", boxSize)
	} else if boxType != "abcd" {
		t.Fatalf("Type not correct: [%s]", boxType)
	}
}

func TestReadBoxes(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	ftypBoxData := []byte{
		// 0: majorBrand (4)
		'a', 'b', 'c', 'd',

		// 4: minorVersion (4)
		0x1, 0x2, 0x3, 0x4,

		// 8: compatibleBrands (8; we chose to add two brands here).
		'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l',
	}

	flagsBytes := make([]byte, 4)
	defaultEndianness.PutUint32(flagsBytes, 0x01020304)

	hdlrBoxData := []byte{
		// 0: version (1)
		0x11,

		// 1: flags (3)
		flagsBytes[0], flagsBytes[1], flagsBytes[2],

		// 4: (reserved spacing) (4)
		0, 0, 0, 0,

		// 8: handler (4)
		'a', 'b', 'c', 'd',

		// 12: (reserved spacing) (12)
		// TODO(dustin): This is probably data that we need to add support for.
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,

		// 24: hdlrName (all remaining)
		't', 'e', 's', 't', 'n', 'a', 'm', 'e',
	}

	data := []byte{}
	pushBox(&data, "ftyp", ftypBoxData)
	pushBox(&data, "hdlr", hdlrBoxData)

	sb := rifs.NewSeekableBufferWithBytes(data)

	size := int64(len(data))

	file := NewFile(sb, size)

	boxes, err := readBoxes(file, 0, size)
	log.PanicIf(err)

	if len(boxes) != 2 {
		t.Fatalf("Expected two boxes: (%d)", len(boxes))
	}

	expectedFtypBox := Box{
		name:  "ftyp",
		start: 0,
		size:  24,
		file:  file,
	}

	expectedHdlrBox := Box{
		name:  "hdlr",
		start: 24,
		size:  40,
		file:  file,
	}

	expectedFtyp, err := ftypBoxFactory{}.New(expectedFtypBox)
	log.PanicIf(err)

	expectedHdlr, err := hdlrBoxFactory{}.New(expectedHdlrBox)
	log.PanicIf(err)

	expectedBoxes := Boxes{
		expectedFtyp,
		expectedHdlr,
	}

	if reflect.DeepEqual(boxes, expectedBoxes) != true {
		t.Fatalf("Boxes not correct.")
	}

	// Use the string functions to validate the actual contents.

	actualPhrases := make([]string, len(boxes))
	for i, box := range boxes {
		actualPhrases[i] = fmt.Sprintf("%s", box)
	}

	expectedPhrases := []string{
		"ftyp<NAME=[ftyp] START=(0) SIZE=(24) MAJOR-BRAND=[abcd] MINOR-VER=(0x01020304) COMPAT-BRANDS=[efgh,ijkl]>",
		"hdlr<NAME=[hdlr] START=(24) SIZE=(40) VER=(0x11) FLAGS=(0x01020300) HANDLER=[abcd] HDLR-NAME=[testname]>",
	}

	if reflect.DeepEqual(actualPhrases, expectedPhrases) != true {
		for i, s := range actualPhrases {
			fmt.Printf("ACTUAL(%d): %s\n", i, s)
		}

		for i, s := range expectedPhrases {
			fmt.Printf("EXPECTED(%d): %s\n", i, s)
		}

		t.Fatalf("String phrases not correct.")
	}
}
