package bmfcommon

import (
	"fmt"

	"github.com/dsoprea/go-logging"
)

// FixedPoint16 encapsulates a uint32 and decodes it to a float32.
type FixedPoint16 struct {
	float32

	rawValue       uint16
	integerLength  int
	mantissaLength int
}

// Float returns the embedded float (the actual value).
func (fp16 FixedPoint16) Float() float32 {

	// TODO(dustin): Add test

	return fp16.float32
}

// Rational returns the individual numerator and denominator components.
func (fp16 FixedPoint16) Rational() (numerator uint16, denominator uint16) {

	// TODO(dustin): Add test

	numerator = fp16.rawValue >> fp16.mantissaLength
	denominator = fp16.rawValue << fp16.integerLength

	return numerator, denominator
}

// String returns a descriptive string.
func (fp16 FixedPoint16) String() string {

	// TODO(dustin): Add test

	numerator, denominator := fp16.Rational()

	return fmt.Sprintf(
		"FixedPoint16<VAL=[%.2f] RAW=(%d) BITs=[%d.%d]> NUM=(%d) DEN=(%d)",
		fp16.float32, fp16.rawValue, fp16.integerLength, fp16.mantissaLength,
		numerator, denominator)
}

// FixedPoint16 returns a float produced by shifting bits in the uint16. Several
// box types have values that are encoded as integers but must be decoded to
// floats before using.
func Uint16ToFixedPoint16(x uint16, integerLength, mantissaLength int) FixedPoint16 {

	// TODO(dustin): Add test

	if integerLength+mantissaLength != 16 {
		log.Panicf("integer bits and mantissa bits do not equal 32")
	}

	n := float32(x >> mantissaLength)
	m := float32(x << integerLength)

	return FixedPoint16{
		float32:        n / m,
		rawValue:       x,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}

// FixedPoint32 encapsulates a uint32 and decodes it to a float32.
type FixedPoint32 struct {
	float32

	rawValue       uint32
	integerLength  int
	mantissaLength int
}

// Float returns the embedded float (the actual value).
func (fp32 FixedPoint32) Float() float32 {

	// TODO(dustin): Add test

	return fp32.float32
}

// Rational returns the individual numerator and denominator components.
func (fp32 FixedPoint32) Rational() (numerator uint32, denominator uint32) {

	// TODO(dustin): Add test

	numerator = fp32.rawValue >> fp32.mantissaLength
	denominator = fp32.rawValue << fp32.integerLength

	return numerator, denominator
}

// String returns a descriptive string.
func (fp32 FixedPoint32) String() string {

	// TODO(dustin): Add test

	numerator, denominator := fp32.Rational()

	return fmt.Sprintf(
		"FixedPoint32<VAL=[%.2f] RAW=(%d) BITs=[%d.%d] NUM=(%d) DEN=(%d)>",
		fp32.float32, fp32.rawValue, fp32.integerLength, fp32.mantissaLength,
		numerator, denominator)
}

// FixedPoint32 returns an encapsulated float produced by shifting bits in the
// uint32.
//
// Several box types have values that are encoded as integers but must be
// decoded to floats before using.
func Uint32ToFixedPoint32(x uint32, integerLength, mantissaLength int) FixedPoint32 {

	// TODO(dustin): Add test

	if integerLength+mantissaLength != 32 {
		log.Panicf("integer bits and mantissa bits do not equal 32")
	}

	n := float32(x >> mantissaLength)
	m := float32(x << integerLength)

	return FixedPoint32{
		float32:        n / m,
		rawValue:       x,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}
