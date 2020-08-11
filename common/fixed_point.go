package bmfcommon

import (
	"fmt"

	"github.com/dsoprea/go-logging"
)

// FixedPoint16 encapsulates a uint32 and decodes it to a float32.
type FixedPoint16 struct {
	rawValue       uint16
	integerLength  int
	mantissaLength int
}

// Rational returns the individual numerator and denominator components.
func (fp16 FixedPoint16) Rational() (numerator uint16, denominator uint16) {

	numerator = fp16.rawValue >> fp16.mantissaLength

	// Clear the bits in front of the denominator by shifting forward all of the
	// way and then back.

	denominatorShiftedLeft := fp16.rawValue << fp16.integerLength
	denominator = denominatorShiftedLeft >> fp16.integerLength

	return numerator, denominator
}

// Float returns the effective float.
func (fp16 FixedPoint16) Float() float32 {

	n, d := fp16.Rational()

	return float32(n) / float32(d)
}

// String returns a descriptive string.
func (fp16 FixedPoint16) String() string {

	numerator, denominator := fp16.Rational()

	return fmt.Sprintf(
		"FixedPoint16<RAW=(%d) BITs=[%d.%d]> NUM=(%d)(0b%016b) DEN=(%d)(0b%016b) VAL=[%.2f]>",
		fp16.rawValue, fp16.integerLength, fp16.mantissaLength, numerator, numerator,
		denominator, denominator, fp16.Float())
}

// Uint16ToFixedPoint16 returns a float produced by shifting bits in the uint16.
// Several box types have values that are encoded as integers but must be
// decoded to floats before using.
func Uint16ToFixedPoint16(x uint16, integerLength, mantissaLength int) FixedPoint16 {

	if integerLength+mantissaLength != 16 {
		log.Panicf("integer bits and mantissa bits do not equal 16")
	}

	return FixedPoint16{
		rawValue:       x,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}

// RationalUint16ToFixedPoint16 returns a FixedPoint16 given the rational and
// encoding-parameters.
func RationalUint16ToFixedPoint16(numerator, denominator uint16, integerLength, mantissaLength int) FixedPoint16 {
	if integerLength+mantissaLength != 16 {
		log.Panicf("bit lengths do not equal 16")
	}

	shiftedNumerator := numerator << integerLength
	rawValue := shiftedNumerator + denominator

	return FixedPoint16{
		rawValue:       rawValue,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}

// FixedPoint32 encapsulates a uint32 and decodes it to a float32.
type FixedPoint32 struct {
	rawValue       uint32
	integerLength  int
	mantissaLength int
}

// Rational returns the individual numerator and denominator components.
func (fp32 FixedPoint32) Rational() (numerator uint32, denominator uint32) {

	numerator = fp32.rawValue >> fp32.mantissaLength

	// Clear the bits in front of the denominator by shifting forward all of the
	// way and then back.

	denominatorShiftedLeft := fp32.rawValue << fp32.integerLength
	denominator = denominatorShiftedLeft >> fp32.integerLength

	return numerator, denominator
}

// Float returns the effective float.
func (fp32 FixedPoint32) Float() float32 {

	n, d := fp32.Rational()

	return float32(n) / float32(d)
}

// String returns a descriptive string.
func (fp32 FixedPoint32) String() string {

	numerator, denominator := fp32.Rational()

	return fmt.Sprintf(
		"FixedPoint32<RAW=(%d) BITs=[%d.%d] NUM=(%d)(0b%08b) DEN=(%d)(0b%08b) VAL=[%.2f]>",
		fp32.rawValue, fp32.integerLength, fp32.mantissaLength, numerator, numerator,
		denominator, denominator, fp32.Float())
}

// Uint32ToFixedPoint32 returns an encapsulated float produced by shifting bits
// in the uint32.
//
// Several box types have values that are encoded as integers but must be
// decoded to floats before using.
func Uint32ToFixedPoint32(x uint32, integerLength, mantissaLength int) FixedPoint32 {

	if integerLength+mantissaLength != 32 {
		log.Panicf("integer bits and mantissa bits do not equal 32")
	}

	return FixedPoint32{
		rawValue:       x,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}

// RationalUint32ToFixedPoint32 returns a FixedPoint16 given the rational and
// encoding-parameters.
func RationalUint32ToFixedPoint32(numerator, denominator uint32, integerLength, mantissaLength int) FixedPoint32 {
	if integerLength+mantissaLength != 32 {
		log.Panicf("bit lengths do not equal 32")
	}

	shiftedNumerator := numerator << integerLength
	rawValue := shiftedNumerator + denominator

	return FixedPoint32{
		rawValue:       rawValue,
		integerLength:  integerLength,
		mantissaLength: mantissaLength,
	}
}
