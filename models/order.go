package models

// The request is very simple, just an int, so this is overkill,
// but in a real world app we might need to model the request with more details
// and it warrants introducing it here, rather than like a DTO on the http handler level
type OrderRequest struct {
	ItemCount int
}

// Not really needed for the task, but an example to support a more realistic UI
type OrderStatus string

const (
	OrderStatusNew     OrderStatus = "new"
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPacked  OrderStatus = "packed"
	OrderStatusShipped OrderStatus = "shipped"
)

type Order struct {
	RequestedItemCount int
	ShippedItemCount   int
	Packs              map[Pack]int
	Status             OrderStatus
}
