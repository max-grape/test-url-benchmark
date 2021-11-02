package omap

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type KeyVal struct {
	Key string
	Val interface{}
}

type OrderedMap []KeyVal

func (om OrderedMap) MarshalBinary() ([]byte, error) {
	return json.Marshal(om)
}

func (om OrderedMap) MarshalIndent(prefix, indent string) ([]byte, error) {
	src, err := om.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	var buf bytes.Buffer

	if err := json.Indent(&buf, src, prefix, indent); err != nil {
		return nil, fmt.Errorf("failed to indent: %w", err)
	}

	return buf.Bytes(), nil
}

func (om OrderedMap) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")

	for i, kv := range om {
		if i != 0 {
			buf.WriteString(",")
		}

		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key `%+v`: %w", kv.Key, err)
		}

		buf.Write(key)
		buf.WriteString(":")

		val, err := json.Marshal(kv.Val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value `%+v`: %w", kv.Val, err)
		}

		buf.Write(val)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
