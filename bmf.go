package bmf

import (
	"os"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// Open opens a physical and returns a File struct. To just use a
// `io.ReadSeeker`, call `bmfcommon.NewFile()` directly.
func Open(path string) (file *bmfcommon.File, err error) {
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

	file = bmfcommon.NewFile(f, size)

	err = file.Parse()
	log.PanicIf(err)

	return file, nil
}
