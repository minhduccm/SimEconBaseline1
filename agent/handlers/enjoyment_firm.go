package handlers

import (
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

func NewEnjoymentFirm(
	initWalletBal float64,
	initOutput float64,
	initWageBudget float64,
	initCapital int,
	capitalProducers []*agentModels.CapitalFirm,
) *agentModels.EnjoymentFirm {

	consumedGoodsFirm := NewConsumedGoodsFirm(
		"Enjoyment",
		initWalletBal,
		initOutput,
		initWageBudget,
		initCapital,
		capitalProducers,
	)

	consumedGoodsFirm.TechCoefficient = 2
	consumedGoodsFirm.Beta = 0.5
	consumedGoodsFirm.Phi = 0.5
	consumedGoodsFirm.Lambda = 0.2
	consumedGoodsFirm.Product = &good.Enjoyment{Quantity: 0}
	consumedGoodsFirm.EUtilThreshold = 0.9
	consumedGoodsFirm.RUtilThreshold = 0.75

	return &agentModels.EnjoymentFirm{
		ConsumedGoodsFirm: consumedGoodsFirm,
	}
}

func (enjoymentFirm *agentModels.EnjoymentFirm) Act() {
	enjoymentFirm.ConsumedGoodsFirm.Act()
}

func (enjoymentFirm *agentModels.EnjoymentFirm) GetGood(goodName string) good.Good {
	if goodName == "Enjoyment" {
		return enjoymentFirm.ConsumedGoodsFirm.Product
	}
	return nil
}
