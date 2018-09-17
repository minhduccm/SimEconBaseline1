package agent

import (
	"fmt"
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
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

	// labor market
	LMkt abstraction.LaborMarket

	// capital market
	CMkt abstraction.CapitalMarket

	/**
	 * Firm prop for general props between firm types
	 */
	Firm *Firm
}

func NewCapitalFirm(
	initWalletBal float64,
	initWageBudget float64,
	lMkt abstraction.LaborMarket,
	cMkt abstraction.CapitalMarket,
) *CapitalFirm {
	firm := NewFirm(initWalletBal)
	firm.WageBudget = initWageBudget

	cFirm := &CapitalFirm{
		TechCoefficient: 20000, // we assume infinite capacity here so we give TechCoefficient a very large value.
		Firm:            firm,
		Beta:            0.5,
		Price:           common.INIT_CAPITAL_PRICE,
		LMkt:            lMkt,
		CMkt:            cMkt,
	}
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[firm.ID]
	cFirm.LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)

	return cFirm
}

func (capitalFirm *CapitalFirm) PrintLastState() {
	fmt.Printf("Capital firm ID: %s \n", capitalFirm.Firm.ID)
	fmt.Printf("**** Number of laborers: %f \n", capitalFirm.Firm.Labor.GetQuantity())
	fmt.Printf("**** Capacity (max output the firm could produce with the current capital and labor): %f \n", capitalFirm.Firm.Capacity)
	fmt.Printf("**** Output: %f \n", capitalFirm.Firm.Output)
	fmt.Printf("**** Wage budget: %f \n", capitalFirm.Firm.WageBudget)
	fmt.Printf("**** Wage per worker: %f \n", capitalFirm.Firm.Wage)
	fmt.Printf("**** Revenue: %f \n", capitalFirm.Firm.Revenue)
	fmt.Printf("**** Profit: %f \n", capitalFirm.Firm.Profit)
	fmt.Printf("**** Marginal profit: %f \n", capitalFirm.Firm.MarginalProfit)
	fmt.Printf("**** Cost ot capital: %f \n", capitalFirm.Firm.CapitalCost)
	fmt.Printf("**** Total cost: %f \n", capitalFirm.Firm.TotalCost)
	fmt.Printf("**** Technology coefficient in the production function: %f \n", capitalFirm.TechCoefficient)
	fmt.Printf("**** Capital price (fixed for now): %f \n\n", capitalFirm.Price)

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

	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[firm.ID]

	revenue := walletAcc.PriIC
	firm.Output = revenue / capitalFirm.Price
	firm.WageBudget = revenue // - loan

	// post to markets
	capitalFirm.LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)
	capitalFirm.CMkt.AddCapitalSellOffer(
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

func (capitalFirm *CapitalFirm) GetGood(goodName string) abstraction.Good {
	if goodName == "Labor" {
		return capitalFirm.Firm.Labor
	}
	return nil
}

func (capitalFirm *CapitalFirm) GetWalletAccountAddress() string {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[capitalFirm.Firm.ID]
	return walletAcc.Address
}

func (capitalFirm *CapitalFirm) GetConsumption(goodName string) float64 {
	return 0.0
}

func (capitalFirm *CapitalFirm) GetID() string {
	return capitalFirm.Firm.ID
}
