package px

import "fmt"

//Unit is a pixel-quantity unit.
type Unit float64

//One px.
const One = Unit(1)

//New returns the given quantity as a unit.
func New(quantity float64) Unit {
	return Unit(quantity)
}

func (u Unit) String() string {
	return fmt.Sprintf("%fpx", u)
}

//Measure implements unit.Unit
func (u Unit) Measure() (quantity float64, reference string) {
	return float64(u), "px"
}
