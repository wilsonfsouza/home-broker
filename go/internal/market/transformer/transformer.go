package transformer

import (
	"github.com/wilsonfsouza/home-broker/go/internal/market/dto"
	"github.com/wilsonfsouza/home-broker/go/internal/market/entity"
)

func TransformOrderInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOrderOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderID:       order.ID,
		InvestorID:    order.Investor.ID,
		AssetID:       order.Asset.ID,
		OrderType:     order.OrderType,
		Status:        order.Status,
		PendingShares: order.PendingShares,
		Shares:        order.Shares,
	}

	var transactionsOutput []*dto.TransactionOutput

	for _, transaction := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: transaction.ID,
			BuyerID:       transaction.BuyingOrder.Investor.ID,
			SellerID:      transaction.SellingOrder.Investor.ID,
			AssetID:       transaction.SellingOrder.Asset.ID,
			Price:         transaction.Price,
			Shares:        transaction.SellingOrder.Shares - transaction.SellingOrder.PendingShares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	output.TransactionOutput = transactionsOutput

	return output
}
