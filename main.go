package main

import (
	"fmt"

	"github.com/goforj/godump"
	"github.com/irreal/order-packs/models"
	"github.com/irreal/order-packs/orders"
)

func main() {
	fmt.Println("Order Packs calculator")

	packs := []models.Pack{250, 500, 1000, 2000, 5000}

	solution, err := orders.CalculatePack(packs, 12001)
	godump.Dump(err)
	godump.Dump(solution)
}
