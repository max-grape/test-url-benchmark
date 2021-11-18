package parser

import (
	"github.com/max-grape/test-url-benchmark/omap"
)

type Parser interface {
	Run(searchParam string) (omap.OrderedMap, error)
}
