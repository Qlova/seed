package shadow

import (
	"fmt"
	"image/color"

	"github.com/qlova/seed"
	"github.com/qlova/seed/css"
	"github.com/qlova/seed/style"
)

//Shadow Defines a shadow that should be applied to the Element, with offset X and Y, Blur and of the specified color.
type Shadow struct {
	X, Y, Blur, Spread style.Unit
	Color              color.Color

	Inset bool
}

//New is an alias for Gradient.
type New = Shadow

func (shadow Shadow) AddTo(c seed.Seed) {
	var inset = ""
	if shadow.Inset {
		inset = "inset "
	}

	if shadow.Color == nil {
		shadow.Color = seed.Black
	}

	css.Set("box-shadow",
		fmt.Sprint(
			inset,
			shadow.X.Unit().Rule(), " ",
			shadow.Y.Unit().Rule(), " ",
			shadow.Blur.Unit().Rule(), " ",
			shadow.Spread.Unit().Rule(), " ",
			css.RGB{Color: shadow.Color}.Rule(),
		),
	).AddTo(c)
}
