package agent

import (
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	market "github.com/ninjadotorg/SimEconBaseline1/market"
)

type CapitalFirm struct {
	/**
	 * technology coefficient in the production function
	 */
	TechCoefficient float64

	/**
	 * sensitivity of output to labor (power on L in the production function
	 */
	Beta float64

	/**
	 * capital price (fixed for now)
	 */
	Price float64

	/**
	 * Firm prop for general props between firm types
	 */
	Firm *Firm
}

func NewCapitalFirm(
	initWalletBal float64,
	initWageBudget float64,
) *CapitalFirm {
	firm := NewFirm(initWalletBal)
	firm.WageBudget = initWageBudget

	econ := economy.GetEconInstance()
	walletAcc := econ.TransactionManager.WalletAccounts[firm.ID]
	LMkt := econ.GetMarket("Labor").(*market.LaborMarket)
	LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)

	return &CapitalFirm{
		TechCoefficient: 20000, // we assume infinite capacity here so we give TechCoefficient a very large value.
		Firm:            firm,
		Beta:            0.5,
		Price:           common.INIT_CAPITAL_PRICE,
	}
}

func (capitalFirm *CapitalFirm) Act() {
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
	LMkt := econ.GetMarket("Labor").(*market.LaborMarket)
	LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)

	CMkt := econ.GetMarket("Capital").(*market.CapitalMarket)
	CMkt.AddCapitalSellOffer(
		capitalFirm,
		capitalFirm.Price,
		int(firm.Capacity),
	)

	walletAcc.PriIC = 0
	labor.Decrease(laborQty)
}

func (capitalFirm *CapitalFirm) ConvertToProduct(laborQty float64) float64 {
	return capitalFirm.TechCoefficient * math.Pow(laborQty, capitalFirm.Beta)
}

func (capitalFirm *CapitalFirm) GetGood(goodName string) good.Good {
	if goodName == "Labor" {
		return capitalFirm.Firm.Labor
	}
	return nil
}
