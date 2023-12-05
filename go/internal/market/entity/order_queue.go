package entity

type OrderQueue struct {
	Orders []*Order
}

func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

func (oq *OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order))
}

func (oq *OrderQueue) Pop() interface{} {
	oldOrders := oq.Orders
	numberOfOrders := len(oldOrders)
	lastOrder := oldOrders[numberOfOrders-1]
	oq.Orders = oldOrders[0 : numberOfOrders-1]
	return lastOrder
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
