package main

import (
	"os"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-mp4"
	"github.com/dsoprea/go-mp4/atom"
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

	atom.Dump(f)
}

// func getFramerate(sampleCounts []uint32, duration, timescale uint32) string {
// 	sc := 1000 * sampleCounts[0]
// 	durationMS := math.Floor(float64(duration) / float64(timescale) * 1000)
// 	return fmt.Sprintf("%.2f", float64(sc)/durationMS)
// }

// func getDurationMS(duration, timescale uint32) string {
// 	return fmt.Sprintf("%.2f", math.Floor(float64(duration)/float64(timescale)*1000))
// }

// func to16(i atom.Fixed32) int {
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

// func getFlags(flags uint32) string {
// 	var f []string
// 	if flags&atom.TrackFlagEnabled == atom.TrackFlagEnabled {
// 		f = append(f, "ENABLED")
// 	}

// 	if flags&atom.TrackFlagInMovie == atom.TrackFlagInMovie {
// 		f = append(f, "IN-MOVIE")
// 	}

// 	if flags&atom.TrackFlagInPreview == atom.TrackFlagInPreview {
// 		f = append(f, "IN-PREVIEW")
// 	}
// 	str := strings.Join(f, " ")
// 	return str
// }
