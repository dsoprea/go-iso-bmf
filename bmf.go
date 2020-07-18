package mp4

import (
	"os"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/boxtype"
)

// Open opens a file and returns a &File{}.
func Open(path string) (file *boxtype.File, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	f, err := os.Open(path)
	log.PanicIf(err)

	s, err := f.Stat()
	log.PanicIf(err)

	size := s.Size()

	file = boxtype.NewFile(f, size)

	err = file.Parse()
	log.PanicIf(err)

	return file, nil
}
