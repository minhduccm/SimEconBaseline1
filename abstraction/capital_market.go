package abstraction

import "github.com/ninjadotorg/SimEconBaseline1/good"

type CapitalMarket interface {
	GetAvgPrice() float64
	AddCapitalBuyOffer(*good.Capital, int)
	AddCapitalSellOffer(Agent, float64, int)
	Perform()
}
