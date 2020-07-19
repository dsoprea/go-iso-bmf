package bmfcommon

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"time"
)

var (
	epochTime = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
)

// GetDurationString Helper function to print a duration value in the form H:MM:SS.MS
func GetDurationString(duration uint32, timescale uint32) string {

	// TODO(dustin): Add test

	durationSec := float64(duration) / float64(timescale)

	hours := math.Floor(durationSec / 3600)
	durationSec -= hours * 3600

	minutes := math.Floor(durationSec / 60)
	durationSec -= minutes * 60

	msec := durationSec * 1000
	durationSec = math.Floor(durationSec)
	msec -= durationSec * 1000
	msec = math.Floor(msec)

	str := fmt.Sprintf("%02.0f:%02.0f:%02.0f:%.0f", hours, minutes, durationSec, msec)

	return str
}

// // TODO(dustin): !! The Fixed16 and Fixed32 types don't seem to add value. Dump?

// // Fixed16 is an 8.8 Fixed Point Decimal notation
// type Fixed16 uint16

// // String returns a descriptive string.
// func (f Fixed16) String() string {

// 	// TODO(dustin): Add test

// 	return fmt.Sprintf("%v", uint16(f)>>8)
// }

// func Fixed16(bytes []byte) Fixed16 {

// 	// TODO(dustin): Add test

// 	return Fixed16(DefaultEndianness.Uint16(bytes))
// }

// // Fixed32 is a 16.16 Fixed Point Decimal notation
// type Fixed32 uint32

// func Fixed32(bytes []byte) Fixed32 {

// 	// TODO(dustin): Add test

// 	return Fixed32(DefaultEndianness.Uint32(bytes))
// }

// bmfcommon.DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {

	// TODO(dustin): Add test

	fmt.Printf("DUMP: ")
	for _, x := range data {
		fmt.Printf("%02x ", x)
	}

	fmt.Printf("\n")
}

// EpochDelta returns the number of seconds since the MP4 epoch.
func EpochDelta(t time.Time) uint32 {

	// TODO(dustin): Add test

	d := t.Sub(epochTime)

	return uint32(math.Floor(float64(d.Seconds())))
}

// EpochToTime returns a the given MP4 epoch as a `time.Time`.
func EpochToTime(epoch uint32) time.Time {

	// TODO(dustin): Add test

	duration := time.Second * time.Duration(epoch)
	t := epochTime.Add(duration)

	return t
}

// NowTime returns a UTC time.Time that has been rounded to seconds.
func NowTime() time.Time {

	// TODO(dustin): Add test

	return time.Now().UTC().Round(time.Second)
}

func PushBox(buffer *[]byte, name string, data []byte) {
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

	DefaultEndianness.PutUint32(
		(*buffer)[start:start+4],
		uint32(len(data))+uint32(BoxHeaderSize))

	copy((*buffer)[start+4:], []byte(name))
	copy((*buffer)[start+8:], data)
}

func PushBytes(buffer *[]byte, x interface{}) {
	var encoded []byte

	if u16, ok := x.(uint16); ok == true {
		encoded = make([]byte, 2)

		DefaultEndianness.PutUint16(
			encoded,
			u16)
	} else if u32, ok := x.(uint32); ok == true {
		encoded = make([]byte, 4)

		DefaultEndianness.PutUint32(
			encoded,
			u32)
	} else if u64, ok := x.(uint64); ok == true {
		encoded = make([]byte, 8)

		DefaultEndianness.PutUint64(
			encoded,
			u64)
	} else if bs, ok := x.([]byte); ok == true {
		*buffer = append(*buffer, bs...)
	} else {
		log.Panicf("can not encode [%v] [%v]", reflect.TypeOf(x), x)
	}

	*buffer = append(*buffer, encoded...)
}
