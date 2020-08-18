package main

import (
	"fmt"
	"os"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-iso-bmf/common"
	"github.com/dsoprea/go-iso-bmf/type"
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

	f, err := os.Open(arguments.Filepath)
	log.PanicIf(err)

	s, err := f.Stat()
	log.PanicIf(err)

	size := s.Size()

	file, err := bmfcommon.NewBmfResource(f, size)
	log.PanicIf(err)

	fmt.Printf("Tree:\n")
	fmt.Printf("\n")

	bmfcommon.Dump(file)

	fmt.Printf("\n")
	fmt.Printf("Item extents:\n")
	fmt.Printf("\n")

	fbi := file.Index()

	ilocCommonBox, found := fbi[bmfcommon.IndexedBoxEntry{"meta.iloc", 0}]
	if found == false {
		fmt.Printf("No ILOC box found. No extents will be written.\n")
		fmt.Printf("\n")
	} else {
		iloc := ilocCommonBox.(*bmftype.IlocBox)

		err = iloc.Dump()
		log.PanicIf(err)
	}

	fmt.Printf("Index:\n")
	fmt.Printf("\n")

	fbi.Dump()

	fmt.Printf("\n")
}
