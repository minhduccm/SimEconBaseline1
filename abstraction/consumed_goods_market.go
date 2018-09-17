package abstraction

type ConsumedGoodsMarket interface {
	GetMarketPrice() float64
	AddBuyOffer(Agent, Demand)
	AddSellOffer(Agent, float64)
	GetTotalSupply() float64
	GetTotalDemand(float64) float64
	Perform()
}
