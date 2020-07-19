package main

import (
	"fmt"
	"os"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-iso-bmf"
	"github.com/dsoprea/go-iso-bmf/common"

	_ "github.com/dsoprea/go-iso-bmf/type"
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

	f, err := bmf.Open(arguments.Filepath)
	log.PanicIf(err)

	fmt.Printf("Tree:\n")
	fmt.Printf("\n")

	bmfcommon.Dump(f)

	fmt.Printf("\n")
	fmt.Printf("Index:\n")
	fmt.Printf("\n")

	fbi := f.Index()
	fbi.Dump()
}
