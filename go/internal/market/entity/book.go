package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order               []*Order
	Transactions        []*Transaction
	OrdersChannelInput  chan *Order
	OrdersChannelOutput chan *Order
	WaitGroup           *sync.WaitGroup
}

func NewBook(ordersChannelInput chan *Order, ordersChannelOutput chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Order:               []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannelInput:  ordersChannelInput,
		OrdersChannelOutput: ordersChannelOutput,
		WaitGroup:           waitGroup,
	}
}

func (b *Book) AddTransaction(transaction *Transaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.SubtractSellOrderPendingShares(minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.SubtractBuyOrderPendingShares(minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	b.Transactions = append(b.Transactions, transaction)
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChannelInput {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)

				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)

					b.AddTransaction(transaction, b.WaitGroup)

					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(sellOrder.Transactions, transaction)

					b.OrdersChannelOutput <- sellOrder
					b.OrdersChannelOutput <- order

					if sellOrder.PendingShares > 0 {
						sellOrders.Push(order)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)

			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				buyOrder := buyOrders.Pop().(*Order)

				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)

					b.AddTransaction(transaction, b.WaitGroup)

					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					b.OrdersChannelOutput <- buyOrder
					b.OrdersChannelOutput <- order

					if buyOrder.PendingShares > 0 {
						buyOrders.Push(order)
					}
				}
			}
		}
	}
}
