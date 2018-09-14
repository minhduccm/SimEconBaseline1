package handlers

import (
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

func NewNecessityFirm(
	initWalletBal float64,
	initOutput float64,
	initWageBudget float64,
	initCapital int,
	capitalProducers []*agentModels.CapitalFirm,
) *agentModels.NecessityFirm {

	consumedGoodsFirm := NewConsumedGoodsFirm(
		"Necessity",
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
	consumedGoodsFirm.Product = &good.Necessity{Quantity: 0}
	consumedGoodsFirm.EUtilThreshold = 0.9
	consumedGoodsFirm.RUtilThreshold = 0.75

	return &agentModels.NecessityFirm{
		ConsumedGoodsFirm: consumedGoodsFirm,
	}
}

func (necessityFirm *agentModels.NecessityFirm) Act() {
	necessityFirm.ConsumedGoodsFirm.Act()
}

func (necessityFirm *agentModels.NecessityFirm) GetGood(goodName string) good.Good {
	if goodName == "Necessity" {
		return necessityFirm.ConsumedGoodsFirm.Product
	}
	return nil
}
