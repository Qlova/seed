package forms

import (
	"qlova.org/seed"
	"qlova.org/seed/client"
	"qlova.org/seed/new/html/form"
	"qlova.org/seed/use/css"
	"qlova.org/seed/use/html"
	"qlova.org/seed/use/html/attr"
	"qlova.org/seed/use/js"
)

//New returns a new HTML form element.
func New(options ...seed.Option) seed.Seed {
	return form.New(
		css.SetDisplay(css.Flex),
		css.SetFlexDirection(css.Column),
		seed.Options(options),
	)
}

//ReportValidity reports the validity of the form.
func ReportValidity(form seed.Seed) client.Bool {
	return js.Bool{Value: html.Element(form).Call("reportValidity")}
}

//SetRequired sets the input element to be required for the form.
func SetRequired() seed.Option {
	return attr.Set("required", "")
}
