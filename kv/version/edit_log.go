package version

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

// marshal encodes edit log to binary data
func (el *EditLog) marshal() ([]byte, error) {
	stream := stream.BinaryWriter()

	// write num of logs
	stream.PutUvarint64(uint64(len(el.logs)))

	// write detail log data
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

// unmarshal create an edit log from its seriealized in buf
func (el *EditLog) unmarshal(buf []byte) error {
	stream := stream.BinaryReader(buf)
	// read num of logs
	count := stream.ReadUvarint64()
	// read detail log data
	for ; count > 0; count-- {
		logType := stream.ReadInt32()
		fn, ok := newLogFuncMap[logType]
		if !ok {
			return fmt.Errorf("cannot get log type new func, type is:[%d]", logType)
		}
		l := fn()
		len := int(stream.ReadUvarint32())
		logData := stream.ReadBytes(len)
		if err := l.Decode(logData); err != nil {
			return fmt.Errorf("unmarshal log data error, type is:[%d],error:%s", logType, err)
		}
		el.Add(l)
	}
	return stream.Error()
}

// apply edit logs into version metadata
func (el *EditLog) apply(version *Version) {
	for _, log := range el.logs {
		log.Apply(version)
	}
}
