package journal

import (
	"fmt"

	"github.com/eleme/lindb/pkg/bufioutil"
)

// Writer writes journals(wal records)
type Writer struct {
	writer bufioutil.BufioWriter
}

// Reader reads wal recores from journal file
type Reader struct {
	reader bufioutil.BufioReader

	value []byte
	eof   bool
	err   error
}

// NewWriter new journals writer instance
func NewWriter(fileName string) (*Writer, error) {
	writer, err := bufioutil.NewBufioWriter(fileName)
	if err != nil {
		return nil, fmt.Errorf("create journal writer error:%s", err)
	}
	return &Writer{
		writer: writer,
	}, nil
}

// NewReader new reader instance
func NewReader(filename string) (*Reader, error) {
	br, err := bufioutil.NewBufioReader(filename)
	if err != nil {
		return nil, fmt.Errorf("create journal reader error:%s", err)
	}
	return &Reader{
		reader: br,
	}, nil
}

// Write writes data into journal file
func (w *Writer) Write(v []byte) error {
	n, err := w.writer.Write(v)
	if err != nil {
		return err
	}
	if n != len(v) {
		return fmt.Errorf("journal write wrong value, written data's length != input length")
	}
	return nil
}

// Sync focus flush data for peresist
func (w *Writer) Sync() error {
	return w.writer.Sync()
}

// Close closes file writer
func (w *Writer) Close() error {
	return w.writer.Close()
}

// Next returns if has data
func (r *Reader) Next() (bool, error) {
	if r.eof || r.err != nil {
		return false, r.err
	}
	r.eof, r.value, r.err = r.reader.Read()
	if r.eof {
		return false, nil
	}
	if r.err != nil {
		return false, r.err
	}
	return len(r.value) > 0, r.err
}

// Record returns wal log record
func (r *Reader) Record() []byte {
	return r.value
}

// Close closes file reader
func (r *Reader) Close() error {
	return r.reader.Close()
}
