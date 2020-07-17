package main

import (
	"os"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-mp4"
	"github.com/dsoprea/go-mp4/mp4box"
)

type parameters struct {
	Filepath  string `short:"f" long:"filepath" required:"true" description:"File-path of image"`
	IsVerbose bool   `short:"v" long:"verbose" description:"Print logging"`
}

var (
	arguments = new(parameters)
)

func main() {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			os.Exit(-2)
		}
	}()

	_, err := flags.Parse(arguments)
	if err != nil {
		os.Exit(-1)
	}

	if arguments.IsVerbose == true {
		cla := log.NewConsoleLogAdapter()
		log.AddAdapter("console", cla)

		scp := log.NewStaticConfigurationProvider()
		scp.SetLevelName(log.LevelNameDebug)

		log.LoadConfiguration(scp)
	}

	f, err := mp4.Open(arguments.Filepath)
	log.PanicIf(err)

	mp4box.Dump(f)
}

// func getFramerate(sampleCounts []uint32, duration, timescale uint32) string {
// 	sc := 1000 * sampleCounts[0]
// 	durationMS := math.Floor(float64(duration) / float64(timescale) * 1000)
// 	return fmt.Sprintf("%.2f", float64(sc)/durationMS)
// }

// func getDurationMS(duration, timescale uint32) string {
// 	return fmt.Sprintf("%.2f", math.Floor(float64(duration)/float64(timescale)*1000))
// }

// func to16(i mp4box.Fixed32) int {
// 	return int(i / (1 << 16))
// }

// func getHandlerType(handler string) string {
// 	var t string
// 	if handler == "vide" {
// 		t = "Video"
// 	} else if handler == "soun" {
// 		t = "Sound"
// 	}
// 	return t
// }

// // Flag constants.
// const (
// 	TrackFlagEnabled   = 0x0001
// 	TrackFlagInMovie   = 0x0002
// 	TrackFlagInPreview = 0x0004
// )

// func getFlags(flags uint32) string {
// 	var f []string
// 	if flags&mp4box.TrackFlagEnabled == mp4box.TrackFlagEnabled {
// 		f = append(f, "ENABLED")
// 	}

// 	if flags&mp4box.TrackFlagInMovie == mp4box.TrackFlagInMovie {
// 		f = append(f, "IN-MOVIE")
// 	}

// 	if flags&mp4box.TrackFlagInPreview == mp4box.TrackFlagInPreview {
// 		f = append(f, "IN-PREVIEW")
// 	}
// 	str := strings.Join(f, " ")
// 	return str
// }
