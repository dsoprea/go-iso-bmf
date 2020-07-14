package atom

import (
	"bytes"
	"testing"
	// "github.com/dsoprea/go-logging"
	// "github.com/dsoprea/go-utility/filesystem"
)

// func newBoxEmptyHeader(name string, data []byte) Box {
// 	b := make([]byte, boxHeaderSize)
// 	b = append(b, data...)

// 	sb := rifs.NewSeekableBufferWithBytes(b)

// 	file := NewFile(sb, int64(len(b)))
// 	return newBox(name, 0, int64(len(data)), file)
// }

func pushBox(buffer *[]byte, name string, data []byte) {
	start := len(*buffer)

	if data == nil {
		data = make([]byte, 0)
	}

	extension := make([]byte, 8+len(data))
	*buffer = append(*buffer, extension...)

	defaultEndianness.PutUint32(
		(*buffer)[start:start+4],
		uint32(len(data))+uint32(boxHeaderSize))

	copy((*buffer)[start+4:], []byte(name))
	copy((*buffer)[start+8:], data)
}

func TestPushBox_One(t *testing.T) {
	b := make([]byte, 0)

	pushBox(&b, "abcd", []byte{1, 2, 3, 4})

	sizeBytes := make([]byte, 4)

	defaultEndianness.PutUint32(sizeBytes, 12)

	expected := []byte{
		sizeBytes[0], sizeBytes[1], sizeBytes[2], sizeBytes[3],
		'a', 'b', 'c', 'd',
		1, 2, 3, 4,
	}

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x\n", b)
	}
}

func TestPushBox_Multiple(t *testing.T) {
	b := make([]byte, 0)

	pushBox(&b, "abcd", []byte{1, 2, 3, 4, 5, 6})
	sizeBytes1 := make([]byte, 4)
	defaultEndianness.PutUint32(sizeBytes1, 14)

	pushBox(&b, "defg", nil)
	sizeBytes2 := make([]byte, 4)
	defaultEndianness.PutUint32(sizeBytes2, 8)

	pushBox(&b, "hijk", []byte{7, 8, 9, 10})
	sizeBytes3 := make([]byte, 4)
	defaultEndianness.PutUint32(sizeBytes3, 12)

	expected := []byte{
		sizeBytes1[0], sizeBytes1[1], sizeBytes1[2], sizeBytes1[3],
		'a', 'b', 'c', 'd',
		1, 2, 3, 4, 5, 6,

		sizeBytes2[0], sizeBytes2[1], sizeBytes2[2], sizeBytes2[3],
		'd', 'e', 'f', 'g',

		sizeBytes3[0], sizeBytes3[1], sizeBytes3[2], sizeBytes3[3],
		'h', 'i', 'j', 'k',
		7, 8, 9, 10,
	}

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x\n", b)
	}
}
