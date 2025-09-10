package orders

import "fmt"

var InvalidOrderItemCountError = fmt.Errorf("requested count is not valid")
var OrderCalculationError = fmt.Errorf("order calculation failed")
