package style

import (
	"github.com/qlova/seed"
	"github.com/qlova/seed/html"
)

//New returns an HTML style element.
func New(options ...seed.Option) seed.Seed {
	return seed.New(html.SetTag("style").And(options...))
}
