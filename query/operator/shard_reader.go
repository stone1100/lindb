package operator

type shardReader struct {
}

func NewShardReader() Operator {
	return &shardReader{}
}
func (op *shardReader) Execute() error {
	return nil
}
