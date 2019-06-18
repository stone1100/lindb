package meta

import (
	"fmt"
	"reflect"

	"github.com/eleme/lindb/pkg/stream"
)

// EditLog includes all metadata edit log
type EditLog struct {
	logs []Log
}

// NewEditLog new EditLog instance
func NewEditLog() *EditLog {
	return &EditLog{}
}

// Add adds edit log into log list
func (el *EditLog) Add(log Log) {
	el.logs = append(el.logs, log)
}

// bytes encodes edit log to binary data
func (el *EditLog) bytes() ([]byte, error) {
	var stream = stream.BinaryWriter()

	for _, log := range el.logs {
		logType := logTypes[reflect.TypeOf(log)]
		stream.PutInt32(logType)
		value, err := log.Encode()
		if err != nil {
			return nil, fmt.Errorf("edit logs encode error: %s", err)
		}
		stream.PutUvarint32(uint32(len(value))) // write log bytes length
		stream.PutBytes(value)                  // write log bytes data
	}
	return stream.Bytes()
}

// apply edit logs into version metadata
func (el *EditLog) apply(version *Version) {
	for _, log := range el.logs {
		log.Apply(version)
	}
}
