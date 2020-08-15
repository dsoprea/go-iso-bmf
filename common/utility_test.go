package bmfcommon

import (
	"bytes"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestPushBox_32(t *testing.T) {
	b := make([]byte, 0)

	PushBox(&b, "abcd", []byte{1, 2, 3, 4})

	sizeBytes := make([]byte, 4)

	DefaultEndianness.PutUint32(sizeBytes, 12)

	expected := []byte{
		sizeBytes[0], sizeBytes[1], sizeBytes[2], sizeBytes[3],
		'a', 'b', 'c', 'd',
		1, 2, 3, 4,
	}

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x\n", b)
	}
}

func TestPushBox_64(t *testing.T) {
	b := make([]byte, 0)

	PushBox(&b, "abcd", Data64BitDescribed{1, 2, 3, 4})

	sizeBytes := make([]byte, 8)

	DefaultEndianness.PutUint64(sizeBytes, 20)

	expected := []byte{
		0, 0, 0, 1,
		'a', 'b', 'c', 'd',
		sizeBytes[0], sizeBytes[1], sizeBytes[2], sizeBytes[3], sizeBytes[4], sizeBytes[5], sizeBytes[6], sizeBytes[7],
		1, 2, 3, 4,
	}

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x\n", b)
	}
}

func TestPushBox_Multiple(t *testing.T) {
	b := make([]byte, 0)

	PushBox(&b, "abcd", []byte{1, 2, 3, 4, 5, 6})
	sizeBytes1 := make([]byte, 4)
	DefaultEndianness.PutUint32(sizeBytes1, 14)

	PushBox(&b, "defg", nil)
	sizeBytes2 := make([]byte, 4)
	DefaultEndianness.PutUint32(sizeBytes2, 8)

	PushBox(&b, "hijk", []byte{7, 8, 9, 10})
	sizeBytes3 := make([]byte, 4)
	DefaultEndianness.PutUint32(sizeBytes3, 12)

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

func TestPushBytes_8_ToEmpty(t *testing.T) {
	value := uint8(0x12)

	var b []byte
	PushBytes(&b, value)

	if len(b) != 1 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	if b[0] != value {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_8_ToNonEmpty(t *testing.T) {
	value := uint8(0x12)

	b := make([]byte, 4)
	PushBytes(&b, value)

	if len(b) != 5 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	if b[4] != value {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_16_ToEmpty(t *testing.T) {
	value := uint16(1234)

	var b []byte
	PushBytes(&b, value)

	if len(b) != 2 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 2)

	DefaultEndianness.PutUint16(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_16_ToNonEmpty(t *testing.T) {
	value := uint16(1234)

	b := make([]byte, 4)
	PushBytes(&b, value)

	if len(b) != 6 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 2)

	DefaultEndianness.PutUint16(
		expected,
		value)

	if bytes.Equal(b[4:], expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_32_ToEmpty(t *testing.T) {
	value := uint32(1234)

	var b []byte
	PushBytes(&b, value)

	if len(b) != 4 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 4)

	DefaultEndianness.PutUint32(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_32_ToNonEmpty(t *testing.T) {
	value := uint32(1234)

	b := make([]byte, 4)
	PushBytes(&b, value)

	if len(b) != 8 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 4)

	DefaultEndianness.PutUint32(
		expected,
		value)

	if bytes.Equal(b[4:], expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_64_ToEmpty(t *testing.T) {
	value := uint64(12345678)

	var b []byte
	PushBytes(&b, value)

	if len(b) != 8 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 8)

	DefaultEndianness.PutUint64(
		expected,
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_64_ToNonEmpty(t *testing.T) {
	value := uint64(12345678)

	b := make([]byte, 4)
	PushBytes(&b, value)

	if len(b) != 12 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	expected := make([]byte, 12)

	DefaultEndianness.PutUint64(
		expected[4:],
		value)

	if bytes.Equal(b, expected) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_Bytes_ToEmpty(t *testing.T) {
	value := []byte{1, 2, 3, 4, 5}

	var b []byte
	PushBytes(&b, value)

	if len(b) != 5 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	if bytes.Equal(b, value) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_Bytes_ToNonEmpty(t *testing.T) {
	value := []byte{1, 2, 3, 4, 5}

	b := make([]byte, 4)
	PushBytes(&b, value)

	if len(b) != 9 {
		t.Fatalf("Length not correct: (%d)", len(b))
	}

	if bytes.Equal(b[4:], value) != true {
		t.Fatalf("Bytes not correct: %x", b)
	}
}

func TestPushBytes_UnsupportedType(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if err.Error() != "can not encode [int64] [11]" {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected panic.")
		}
	}()

	b := make([]byte, 0)
	PushBytes(&b, int64(11))
}

func TestDumpBytes(t *testing.T) {
	b := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	}

	DumpBytes(b)
}
