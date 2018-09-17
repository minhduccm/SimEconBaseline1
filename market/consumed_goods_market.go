package market

import (
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
)

const (
	zeta = 0.1
)

type ConsumedGoodsMarket struct {
	GoodName          string
	InitLow           float64
	InitHigh          float64
	BuyOffers         []*BuyOffer
	SellOffers        []*SellOffer
	MarketPrice       float64
	MarketGoodVolume  float64
	MarketMoneyVolume float64
	MarketSupply      float64
}

type BuyOffer struct {
	Buyer abstraction.Agent
	Demd  abstraction.Demand
}

type SellOffer struct {
	Seller abstraction.Agent
	Qty    float64
}

func NewConsumedGoodsMarket(
	goodName string,
	initLow float64,
	initHigh float64,
) *ConsumedGoodsMarket {
	return &ConsumedGoodsMarket{
		GoodName:   goodName,
		InitLow:    initLow,
		InitHigh:   initHigh,
		BuyOffers:  []*BuyOffer{},
		SellOffers: []*SellOffer{},
	}
}

func (consumedGoodsMarket *ConsumedGoodsMarket) AddBuyOffer(
	buyer abstraction.Agent,
	demand abstraction.Demand,
) {
	offer := &BuyOffer{
		Buyer: buyer,
		Demd:  demand,
	}
	consumedGoodsMarket.BuyOffers = append(consumedGoodsMarket.BuyOffers, offer)
}

func (consumedGoodsMarket *ConsumedGoodsMarket) AddSellOffer(
	seller abstraction.Agent,
	qty float64,
) {
	offer := &SellOffer{
		Seller: seller,
		Qty:    qty,
	}
	consumedGoodsMarket.SellOffers = append(consumedGoodsMarket.SellOffers, offer)
}

func (consumedGoodsMarket *ConsumedGoodsMarket) GetTotalSupply() float64 {
	var supply float64 = 0
	for _, offer := range consumedGoodsMarket.SellOffers {
		supply += offer.Qty
	}
	return supply
}

func (consumedGoodsMarket *ConsumedGoodsMarket) GetTotalDemand(price float64) float64 {
	var demand float64 = 0
	for _, offer := range consumedGoodsMarket.BuyOffers {
		consumption := offer.Buyer.GetConsumption(consumedGoodsMarket.GoodName)
		demand += offer.Demd.GetDemand(price, consumption)
	}
	return demand
}

func (consumedGoodsMarket *ConsumedGoodsMarket) Perform() {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	var low, high, price float64
	if common.TimeStep == 0 {
		low = consumedGoodsMarket.InitLow
		high = consumedGoodsMarket.InitHigh
	} else {
		low = consumedGoodsMarket.MarketPrice * (1 - zeta)
		high = consumedGoodsMarket.MarketPrice * (1 + zeta)
	}

	totalSuppy := consumedGoodsMarket.GetTotalSupply()
	var totalDemand float64 = 0
	for {
		price = (low + high) / 2
		totalDemand = consumedGoodsMarket.GetTotalDemand(price)
		if math.Abs(totalDemand-totalSuppy) < 0.1 || math.Abs(high-low) < 0.01 {
			break
		}
		if totalDemand > totalSuppy {
			low = price
		} else {
			high = price
		}
	}

	vol := math.Min(totalDemand, totalSuppy)
	if vol > 0.1 {
		for _, offer := range consumedGoodsMarket.BuyOffers {
			consumption := offer.Buyer.GetConsumption(consumedGoodsMarket.GoodName)
			qty := offer.Demd.GetDemand(price, consumption) / totalDemand * vol
			payAmt := qty * price
			transactionManager.PayFrom(
				// offer.Buyer.GetWalletAccountAddress(),
				offer.Buyer.GetID(),
				payAmt,
			)
			good := offer.Buyer.GetGood(consumedGoodsMarket.GoodName)
			good.Increase(qty)
		}
		for _, offer := range consumedGoodsMarket.SellOffers {
			qty := offer.Qty / totalSuppy * vol
			payAmt := qty * price
			transactionManager.PayTo(
				// offer.Seller.GetWalletAccountAddress(),
				offer.Seller.GetID(),
				payAmt,
				common.PRIIC,
			)
			good := offer.Seller.GetGood(consumedGoodsMarket.GoodName)
			good.Decrease(qty)
		}
	}

	consumedGoodsMarket.MarketPrice = price
	consumedGoodsMarket.MarketGoodVolume = vol
	consumedGoodsMarket.MarketMoneyVolume = price * vol
	consumedGoodsMarket.MarketSupply = totalSuppy

	// reset
	consumedGoodsMarket.BuyOffers = []*BuyOffer{}
	consumedGoodsMarket.SellOffers = []*SellOffer{}
}

func (consumedGoodsMarket *ConsumedGoodsMarket) GetMarketPrice() float64 {
	return consumedGoodsMarket.MarketPrice
}
