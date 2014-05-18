package meiporo

import (
	"log"

	"github.com/lestrrat/go-xslate"
)

var Xslate *xslate.Xslate

func init() {
	xt, err := xslate.New()
	if err != nil {
		log.Fatalf("Failed to create xslate: %s", err)
	}
	Xslate = xt
}
