package agent

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type EnjoymentFirm struct {
	ConsumedGoodsFirm *ConsumedGoodsFirm
}

func NewEnjoymentFirm(
	initWalletBal float64,
	initOutput float64,
	initWageBudget float64,
	initCapital int,
	capitalProducers []*CapitalFirm,
) *EnjoymentFirm {

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

	return &EnjoymentFirm{
		ConsumedGoodsFirm: consumedGoodsFirm,
	}
}

func (enjoymentFirm *EnjoymentFirm) Act() {
	enjoymentFirm.ConsumedGoodsFirm.Act()
}

func (enjoymentFirm *EnjoymentFirm) GetGood(goodName string) good.Good {
	if goodName == "Enjoyment" {
		return enjoymentFirm.ConsumedGoodsFirm.Product
	}
	return nil
}
