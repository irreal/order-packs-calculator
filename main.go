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

	service := orders.NewService(1000000000)
	order, err := service.CreateOrder(models.OrderRequest{ItemCount: 12001}, packs)

	godump.Dump(err)
	godump.Dump(order)
}
