package models

import (
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type CapitalMarket struct {
	// buy offers
	CapitalBuyOffers []*CapitalBuyOffer

	// sell offers
	CapitalSellOffers []*CapitalSellOffer

	// volume of capital good traded
	MktGoodVol float64

	// totalMetric supply of capital good
	Supply float64

	// sum of reciprocal of prices of all sell offers
	TotalMetric float64

	// sum of prices of all sell offers
	TotalPrice float64

	// average capital price
	AvgPrice float64
}

type CapitalBuyOffer struct {
	Capital  *good.Capital
	Quantity int
}

type CapitalSellOffer struct {
	Seller   *agentModels.CapitalFirm
	Price    float64
	Capacity int
}
