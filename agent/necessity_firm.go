package agent

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type NecessityFirm struct {
	ConsumedGoodsFirm *ConsumedGoodsFirm
}

func NewNecessityFirm(
	initWalletBal float64,
	initOutput float64,
	initWageBudget float64,
	initCapital int,
	capitalProducers []*CapitalFirm,
) *NecessityFirm {

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

	return &NecessityFirm{
		ConsumedGoodsFirm: consumedGoodsFirm,
	}
}

func (necessityFirm *NecessityFirm) Act() {
	necessityFirm.ConsumedGoodsFirm.Act()
}

func (necessityFirm *NecessityFirm) GetGood(goodName string) good.Good {
	if goodName == "Necessity" {
		return necessityFirm.ConsumedGoodsFirm.Product
	}
	return nil
}
