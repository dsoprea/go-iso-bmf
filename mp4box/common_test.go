package mp4box

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func pushBox(buffer *[]byte, name string, data []byte) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.Panic(err)
		}
	}()

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

func pushBytes(buffer *[]byte, x interface{}) {
	var encoded []byte

	if u16, ok := x.(uint16); ok == true {
		encoded = make([]byte, 2)

		defaultEndianness.PutUint16(
			encoded,
			u16)
	} else if u32, ok := x.(uint32); ok == true {
		encoded = make([]byte, 4)

		defaultEndianness.PutUint32(
			encoded,
			u32)
	} else if u64, ok := x.(uint64); ok == true {
		encoded = make([]byte, 8)

		defaultEndianness.PutUint64(
			encoded,
			u64)
	} else if bs, ok := x.([]byte); ok == true {
		*buffer = append(*buffer, bs...)
	} else {
		log.Panicf("can not encode [%v] [%v]", reflect.TypeOf(x), x)
	}

	*buffer = append(*buffer, encoded...)
}

func TestPushBytes_16_ToEmpty(t *testing.T) {
	value := uint16(1234)

	var b []byte
	pushBytes(&b, value)

	if len(b) != 2 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 2)

	defaultEndianness.PutUint16(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_16_ToNonEmpty(t *testing.T) {
	value := uint16(1234)

	b := make([]byte, 4)
	pushBytes(&b, value)

	if len(b) != 6 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 2)

	defaultEndianness.PutUint16(
		expected,
		value)

	if bytes.Equal(b[4:], expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_32_ToEmpty(t *testing.T) {
	value := uint32(1234)

	var b []byte
	pushBytes(&b, value)

	if len(b) != 4 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 4)

	defaultEndianness.PutUint32(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_32_ToNonEmpty(t *testing.T) {
	value := uint32(1234)

	b := make([]byte, 4)
	pushBytes(&b, value)

	if len(b) != 8 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 4)

	defaultEndianness.PutUint32(
		expected,
		value)

	if bytes.Equal(b[4:], expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_64_ToEmpty(t *testing.T) {
	value := uint64(12345678)

	var b []byte
	pushBytes(&b, value)

	if len(b) != 8 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 8)

	defaultEndianness.PutUint64(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_64_ToNonEmpty(t *testing.T) {
	value := uint64(12345678)

	b := make([]byte, 4)
	pushBytes(&b, value)

	if len(b) != 12 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 12)

	defaultEndianness.PutUint64(
		expected[4:],
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}
