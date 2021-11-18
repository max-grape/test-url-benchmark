package omap_test

import (
	"testing"

	"github.com/max-grape/test-url-benchmark/omap"
	"github.com/stretchr/testify/assert"
)

func TestOrderedMap_MarshalBinary(t *testing.T) {
	om := omap.OrderedMap{
		{Key: "foo", Val: 123},
		{Key: "bar", Val: struct{}{}},
	}

	actual, err := om.MarshalBinary()

	assert.NoError(t, err)
	assert.Equal(t, `[{"Key":"foo","Val":123},{"Key":"bar","Val":{}}]`, string(actual))
}

func TestOrderedMap_MarshalIndent(t *testing.T) {
	om := omap.OrderedMap{
		{Key: "foo", Val: 123},
		{Key: "bar", Val: struct{}{}},
	}

	actual, err := om.MarshalIndent("", "    ")

	assert.NoError(t, err)
	assert.Equal(t, "{\n    \"foo\": 123,\n    \"bar\": {}\n}", string(actual))
}
