package agent

import (
	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
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
	pMkt abstraction.ConsumedGoodsMarket,
	lMkt abstraction.LaborMarket,
	cMkt abstraction.CapitalMarket,
) *EnjoymentFirm {

	consumedGoodsFirm := NewConsumedGoodsFirm(
		"Enjoyment",
		initWalletBal,
		initOutput,
		initWageBudget,
		initCapital,
		capitalProducers,
		pMkt,
		lMkt,
		cMkt,
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

func (enjoymentFirm *EnjoymentFirm) GetGood(goodName string) abstraction.Good {
	if goodName == "Enjoyment" {
		return enjoymentFirm.ConsumedGoodsFirm.Product
	}
	return nil
}

func (enjoymentFirm *EnjoymentFirm) GetWalletAccountAddress() string {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[enjoymentFirm.ConsumedGoodsFirm.Firm.ID]
	return walletAcc.Address
}

func (enjoymentFirm *EnjoymentFirm) GetConsumption(goodName string) float64 {
	return 0.0
}

func (enjoymentFirm *EnjoymentFirm) GetID() string {
	return enjoymentFirm.ConsumedGoodsFirm.GetID()
}
