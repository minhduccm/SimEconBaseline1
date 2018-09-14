package handlers

import (
	"math/rand"

	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	marketModels "github.com/ninjadotorg/SimEconBaseline1/market/models"
)

func NewCapitalMarket() *marketModels.CapitalMarket {
	return &marketModels.CapitalMarket{
		CapitalBuyOffers:  []*CapitalBuyOffer{},
		CapitalSellOffers: []*CapitalSellOffer{},
		Supply:            0,
		TotalMetric:       0,
		TotalPrice:        0,
	}
}

func (cm *marketModels.CapitalMarket) AddCapitalBuyOffer(
	capital *good.Capital,
	qty int,
) {
	offer := &marketModels.CapitalBuyOffer{
		Capital:  capital,
		Quantity: qty,
	}
	cm.CapitalBuyOffers = append(cm.CapitalBuyOffers, offer)
}

func (cm *marketModels.CapitalMarket) AddCapitalSellOffer(
	seller *agentModels.CapitalFirm,
	price float64,
	capacity int,
) {
	offer := &marketModels.CapitalSellOffer{
		Seller:   seller,
		Price:    price,
		Capacity: capacity,
	}
	cm.CapitalSellOffers = append(cm.CapitalSellOffers, offer)
	cm.TotalMetric += 1 / price
	cm.TotalPrice += price * capacity
	cm.Supply += capacity
}

func (cm *marketModels.CapitalMarket) Perform() {
	cm.MktGoodVol = 0
	// TODO: should shuffle BuyOffers

	for _, buyOffer := range cm.CapitalBuyOffers {
		if cm.Supply == 0 {
			return
		}
		var picked *marketModels.CapitalSellOffer = nil
		var removingIdx int = -1
		for {
			if picked != nil {
				break
			}
			var winner float64 = rand.Float64() // random num in [0.0, 1)
			var val float64 = 0
			for idx, sellOffer := range cm.CapitalSellOffers {
				val += 1 / sellOffer.Price / cm.TotalMetric
				if val > winner {
					if buyOffer.Quantity <= sellOffer.Capacity {
						picked = sellOffer
						removingIdx = idx
						picked.Capacity -= buyOffer.Quantity
						cm.MktGoodVol += buyOffer.Quantity
					}
					break
				}
			}
		}
		buyOffer.Capital.Add(buyOffer.Quantity, picked.Price, common.CAPITAL_LIFE, picked.Seller.ID)
		if picked.Capacity <= 0 {
			cm.CapitalSellOffers = append(cm.CapitalSellOffers[:removingIdx], cm.CapitalSellOffers[removingIdx+1:]...)
		}
	}
	cm.TotalMetric = 0
	cm.AvgPrice = cm.TotalPrice / cm.Supply
	cm.TotalPrice = 0
	cm.Supply = 0
	cm.CapitalBuyOffers = []*marketModels.CapitalBuyOffer{}
	cm.CapitalSellOffers = []*marketModels.CapitalSellOffer{}
}
