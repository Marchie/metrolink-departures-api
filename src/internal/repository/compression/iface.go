package compression

import "io"

type Extractor interface {
	ExtractFile(zipData []byte, filename string) (io.ReadCloser, error)
}
