package visible

import (
	"fmt"

	"qlova.org/seed"
	"qlova.org/seed/client"
	"qlova.org/seed/client/clientside"
	"qlova.org/seed/web/js"
)

//On executes the provided script whenever this seed is rendered.
func On(do client.Script) seed.Option {
	return client.On("render", do.GetScript())
}

//When renders the give seeds only when the condition is true.
func When(condition client.Bool, seeds ...seed.Seed) seed.Option {
	return seed.NewOption(func(c seed.Seed) {
		clientside.Hook(condition, c)

		c.With(On(client.If(condition, js.Script(func(q js.Ctx) {
			for _, child := range seeds {
				if child.ID() == 0 {
					continue
				}
				fmt.Fprintf(q, `%v.style.display = ""; if (%[1]v.onvisible) %[1]v.onvisible();`, client.Element(child))
				child.AddTo(c)
			}
		})).Else(js.Script(func(q js.Ctx) {
			for _, child := range seeds {
				if child.ID() == 0 {
					continue
				}
				fmt.Fprintf(q, `%v.style.display = "none";  if (%[1]v.onhidden) %[1]v.onhidden();`, client.Element(child))
				child.AddTo(c)
			}
		}))))
	})
}
