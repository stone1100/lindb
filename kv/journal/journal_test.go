package journal

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJournal(t *testing.T) {
	file := "jouranl.test"
	defer os.Remove(file)
	writer, _ := NewWriter(file)
	writer.Write([]byte("test0"))
	writer.Write([]byte("test1"))
	writer.Write([]byte("test2"))
	writer.Close()

	reader, _ := NewReader(file)
	var count int
	// read log
	for {
		next, err := reader.Next()
		assert.Nil(t, err)
		if !next {
			break
		}
		record := reader.Record()
		assert.Equal(t, []byte(fmt.Sprintf("test%d", count)), record)
		count++
	}
	reader.Close()

	assert.Equal(t, 3, count)
}
