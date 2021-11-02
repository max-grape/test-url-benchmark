package parser

import (
	"github.com/max-grape/test-revo/omap"
)

type Parser interface {
	Run(searchParam string) (omap.OrderedMap, error)
}
