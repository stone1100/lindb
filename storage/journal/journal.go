package journal

import (
	"fmt"

	"github.com/eleme/lindb/pkg/io"
)

// Writer writes journals(wal records) to an underlying io.Writer
type Writer struct {
	fw *io.FileWriter
}

// Reader reads wal recores from journal file
type Reader struct {
}

// NewWriter new journals writer instance
func NewWriter(fileName string) (*Writer, error) {
	writer, err := io.NewWriter(fileName)
	if err != nil {
		return nil, fmt.Errorf("create journal writer error:%s", err)
	}
	return &Writer{
		fw: writer,
	}, nil
}

// NewReader new reader instance
func NewReader(filename string) {
	//ioutil.ReadFile(filename)
	//os.Open()
}

// Write writes data into journal file
func (w *Writer) Write(v []byte) error {
	n, err := w.fw.Write(v)
	if err != nil {
		return err
	}
	if n != len(v) {
		return fmt.Errorf("journal write wrong value, written data's length != input length")
	}

	return nil
}

func (r *Reader) Next() (bool, error) {

	return true, nil
}
