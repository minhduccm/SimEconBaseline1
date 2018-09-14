package handlers

import (
	"math"

	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	marketModels "github.com/ninjadotorg/SimEconBaseline1/market/models"
)

func NewCapitalFirm(
	initWalletBal float64,
	initWageBudget float64,
) *agentModels.CapitalFirm {
	firm := NewFirm(initWalletBal)
	firm.WageBudget = initWageBudget

	econ := economy.GetEconInstance()
	walletAcc := econ.TransactionManager.WalletAccounts[firm.ID]
	LMkt := econ.GetMarket("Labor").(*marketModels.LaborMarket)
	LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)

	return &agentModels.CapitalFirm{
		TechCoefficient: 20000, // we assume infinite capacity here so we give TechCoefficient a very large value.
		Firm:            firm,
		Beta:            0.5,
		Price:           common.INIT_CAPITAL_PRICE,
	}
}

func (capitalFirm *agentModels.CapitalFirm) Act() {
	// Capital firms are not supposed to have loans in this
	// phase. But if for some reason a firm has a positive
	// loan, pay back that loan.
	// loan = -Bank.getBalance(getID(), Bank.SAVINGS);
	// if (loan > 0)
	// 	Bank.deposit(getID(), loan);

	firm := capitalFirm.Firm
	labor := firm.Labor
	laborQty := labor.Quantity
	firm.Capacity = capitalFirm.ConvertToProduct(laborQty)
	if laborQty > 0 {
		firm.Wage = firm.WageBudget / laborQty
	} else {
		firm.Wage = 0
	}
	econ := economy.GetEconInstance()
	walletAcc := econ.TransactionManager.WalletAccounts[firm.ID]
	revenue := walletAcc.PriIC
	firm.Output = revenue / capitalFirm.Price
	firm.WageBudget = revenue // - loan

	// post to markets
	LMkt := econ.GetMarket("Labor").(*marketModels.LaborMarket)
	LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)

	CMkt := econ.GetMarket("Capital").(*marketModels.CapitalMarket)
	CMkt.AddCapitalSellOffer(
		capitalFirm,
		capitalFirm.Price,
		int(firm.Capacity),
	)

	walletAcc.PriIC = 0
	labor.Decrease(laborQty)
}

func (capitalFirm *agentModels.CapitalFirm) ConvertToProduct(laborQty float64) float64 {
	return capitalFirm.TechCoefficient * math.Pow(laborQty, capitalFirm.Beta)
}

func (capitalFirm *agentModels.CapitalFirm) GetGood(goodName string) good.Good {
	if goodName == "Labor" {
		return capitalFirm.Firm.Labor
	}
	return nil
}
