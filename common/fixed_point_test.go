package bmfcommon

import (
	"testing"
)

func TestFixedPoint16_Rational_8_8(t *testing.T) {
	fp16 := FixedPoint16{
		rawValue:       0b0000111111110000,
		integerLength:  8,
		mantissaLength: 8,
	}

	n, d := fp16.Rational()

	if n != 0b00001111 {
		t.Fatalf("Numerator not correct: (0b%016b)", n)
	}

	if d != 0b11110000 {
		t.Fatalf("Denominator not correct: (0b%016b)", d)
	}
}

func TestFixedPoint16_Rational_12_4(t *testing.T) {
	fp16 := FixedPoint16{
		rawValue:       0b0000111111110000,
		integerLength:  12,
		mantissaLength: 4,
	}

	n, d := fp16.Rational()

	if n != 0b11111111 {
		t.Fatalf("Numerator not correct: (0b%016b)", n)
	}

	if d != 0b0 {
		t.Fatalf("Denominator not correct: (0b%016b)", d)
	}
}

func TestFixedPoint16_Float(t *testing.T) {
	fp16 := FixedPoint16{
		rawValue:       0b0000111111110000,
		integerLength:  8,
		mantissaLength: 8,
	}

	n, d := fp16.Rational()

	f := fp16.Float()

	nRestored := uint16(f * float32(d))
	if nRestored != n {
		t.Fatalf("Float() not correct.")
	}
}

func TestFixedPoint16_String(t *testing.T) {
	fp16 := FixedPoint16{
		rawValue:       0b0000111111110000,
		integerLength:  8,
		mantissaLength: 8,
	}

	if fp16.String() != "FixedPoint16<RAW=(4080) BITs=[8.8]> NUM=(15)(0b0000000000001111) DEN=(240)(0b0000000011110000) VAL=[0.06]>" {
		t.Fatalf("String() not correct: [%s]", fp16.String())
	}
}

func TestUint16ToFixedPoint16(t *testing.T) {
	fp16 := Uint16ToFixedPoint16(0b0000111111110000, 8, 8)

	if fp16.String() != "FixedPoint16<RAW=(4080) BITs=[8.8]> NUM=(15)(0b0000000000001111) DEN=(240)(0b0000000011110000) VAL=[0.06]>" {
		t.Fatalf("String() not correct: [%s]", fp16.String())
	}
}

func TestFixedPoint32_Rational_16_16(t *testing.T) {
	fp32 := FixedPoint32{
		rawValue:       0b00001111000000000000000011110000,
		integerLength:  16,
		mantissaLength: 16,
	}

	n, d := fp32.Rational()

	if n != 0b0000111100000000 {
		t.Fatalf("Numerator not correct: (0b%032b)", n)
	}

	if d != 0b0000000011110000 {
		t.Fatalf("Denominator not correct: (0b%032b)", d)
	}
}

func TestFixedPoint32_Rational_20_12(t *testing.T) {
	fp32 := FixedPoint32{
		rawValue:       0b00001111000000000000000011110000,
		integerLength:  20,
		mantissaLength: 12,
	}

	n, d := fp32.Rational()

	if n != 0b00001111000000000000 {
		t.Fatalf("Numerator not correct: (0b%032b)", n)
	}

	if d != 0b000011110000 {
		t.Fatalf("Denominator not correct: (0b%032b)", d)
	}
}

func TestFixedPoint32_Float(t *testing.T) {
	fp32 := FixedPoint32{
		rawValue:       0b00001111000000000000000011110000,
		integerLength:  16,
		mantissaLength: 16,
	}

	n, d := fp32.Rational()

	f := fp32.Float()

	nRestored := uint32(f * float32(d))
	if nRestored != n {
		t.Fatalf("Float() not correct.")
	}
}

func TestFixedPoint32_String(t *testing.T) {
	fp32 := FixedPoint32{
		rawValue:       0b00001111000000000000000011110000,
		integerLength:  16,
		mantissaLength: 16,
	}

	if fp32.String() != "FixedPoint32<RAW=(251658480) BITs=[16.16] NUM=(3840)(0b111100000000) DEN=(240)(0b11110000) VAL=[16.00]>" {
		t.Fatalf("String() not correct: [%s]", fp32.String())
	}
}

func TestUint32ToFixedPoint32(t *testing.T) {
	fp32 := Uint32ToFixedPoint32(0b00001111000000000000000011110000, 16, 16)

	if fp32.String() != "FixedPoint32<RAW=(251658480) BITs=[16.16] NUM=(3840)(0b111100000000) DEN=(240)(0b11110000) VAL=[16.00]>" {
		t.Fatalf("String() not correct: [%s]", fp32.String())
	}
}
