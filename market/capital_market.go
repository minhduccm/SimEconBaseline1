package market

import (
	"math/rand"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
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
	Seller   abstraction.Agent
	Price    float64
	Capacity int
}

func NewCapitalMarket() *CapitalMarket {
	return &CapitalMarket{
		CapitalBuyOffers:  []*CapitalBuyOffer{},
		CapitalSellOffers: []*CapitalSellOffer{},
		Supply:            0,
		TotalMetric:       0,
		TotalPrice:        0,
	}
}

func (cm *CapitalMarket) AddCapitalBuyOffer(
	capital *good.Capital,
	qty int,
) {
	offer := &CapitalBuyOffer{
		Capital:  capital,
		Quantity: qty,
	}
	cm.CapitalBuyOffers = append(cm.CapitalBuyOffers, offer)
}

func (cm *CapitalMarket) AddCapitalSellOffer(
	seller abstraction.Agent,
	price float64,
	capacity int,
) {
	offer := &CapitalSellOffer{
		Seller:   seller,
		Price:    price,
		Capacity: capacity,
	}
	cm.CapitalSellOffers = append(cm.CapitalSellOffers, offer)
	cm.TotalMetric += 1 / price
	cm.TotalPrice += price * float64(capacity)
	cm.Supply += float64(capacity)
}

func (cm *CapitalMarket) Perform() {
	cm.MktGoodVol = 0
	// TODO: should shuffle BuyOffers

	for _, buyOffer := range cm.CapitalBuyOffers {
		if cm.Supply == 0 {
			return
		}
		var picked *CapitalSellOffer = nil
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
						cm.MktGoodVol += float64(buyOffer.Quantity)
					}
					break
				}
			}
		}
		buyOffer.Capital.Add(buyOffer.Quantity, picked.Price, common.CAPITAL_LIFE, picked.Seller.GetID())
		if picked.Capacity <= 0 {
			cm.CapitalSellOffers = append(cm.CapitalSellOffers[:removingIdx], cm.CapitalSellOffers[removingIdx+1:]...)
		}
	}
	cm.TotalMetric = 0
	cm.AvgPrice = cm.TotalPrice / cm.Supply
	cm.TotalPrice = 0
	cm.Supply = 0
	cm.CapitalBuyOffers = []*CapitalBuyOffer{}
	cm.CapitalSellOffers = []*CapitalSellOffer{}
}

func (cm *CapitalMarket) GetAvgPrice() float64 {
	return cm.AvgPrice
}
