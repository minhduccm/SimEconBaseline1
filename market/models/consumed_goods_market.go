package models

import (
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
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
	Buyer agentModels.Agent
	Demd  Demand
}

type SellOffer struct {
	Seller agentModels.Agent
	Qty    float64
}
