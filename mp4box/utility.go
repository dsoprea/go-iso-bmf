package mp4box

import (
	"fmt"
	"math"
	"time"
)

var (
	epochTime = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
)

// GetDurationString Helper function to print a duration value in the form H:MM:SS.MS
func GetDurationString(duration uint32, timescale uint32) string {
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

// TODO(dustin): !! The Fixed16 and Fixed32 types don't seem to add value. Dump?

// Fixed16 is an 8.8 Fixed Point Decimal notation
type Fixed16 uint16

// String returns a descriptive string.
func (f Fixed16) String() string {
	return fmt.Sprintf("%v", uint16(f)>>8)
}

func fixed16(bytes []byte) Fixed16 {
	return Fixed16(defaultEndianness.Uint16(bytes))
}

// Fixed32 is a 16.16 Fixed Point Decimal notation
type Fixed32 uint32

func fixed32(bytes []byte) Fixed32 {
	return Fixed32(defaultEndianness.Uint32(bytes))
}

// DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {
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
