package agent

import (
	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
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
	pMkt abstraction.ConsumedGoodsMarket,
	lMkt abstraction.LaborMarket,
	cMkt abstraction.CapitalMarket,
) *NecessityFirm {

	consumedGoodsFirm := NewConsumedGoodsFirm(
		"Necessity",
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

func (necessityFirm *NecessityFirm) GetGood(goodName string) abstraction.Good {
	if goodName == "Necessity" {
		return necessityFirm.ConsumedGoodsFirm.Product
	}
	return nil
}

func (necessityFirm *NecessityFirm) GetWalletAccountAddress() string {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[necessityFirm.ConsumedGoodsFirm.Firm.ID]
	return walletAcc.Address
}

func (necessityFirm *NecessityFirm) GetConsumption(goodName string) float64 {
	return 0.0
}

func (necessityFirm *NecessityFirm) GetID() string {
	return necessityFirm.ConsumedGoodsFirm.GetID()
}
