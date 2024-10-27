package replication

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Request struct {
	PrevSegmentName string
}

type Response struct {
	Ok          bool
	SegmentName string
	Data        []byte
}

func respReqEncode(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(v); err != nil {
		return nil, fmt.Errorf("failed to encode value: %w", err)
	}
	return buffer.Bytes(), nil
}

func respReqDecode[T any](val *T, data []byte) error {
	buffer := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buffer).Decode(val); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	return nil
}
