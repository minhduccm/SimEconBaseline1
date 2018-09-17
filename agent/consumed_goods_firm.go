package agent

import (
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
	"github.com/ninjadotorg/SimEconBaseline1/util"
)

type ConsumedGoodsFirm struct {
	/**
	 * product name: Necessity or Enjoyment
	 */
	ProductName string

	/**
	 * technology coefficient in the production function
	 */
	TechCoefficient float64

	/**
	 * sensitivity of output to labor (power on L in the production function
	 */
	Beta float64

	/**
	 * sensitivity of output to marginal profit
	 */
	Phi float64

	/**
	 * sensitivity of wage to money flow gap
	 */
	Lambda float64

	/**
	 * minimal capacity utilization to allow capital expansion
	 */
	EUtilThreshold float64

	/**
	 * minimal capacity utilization to allow capital replacement
	 */
	RUtilThreshold float64

	/**
	 * product the firm is producing/selling (enjoyment or necessity)
	 */
	Product abstraction.Good

	// Product market
	PMkt abstraction.ConsumedGoodsMarket

	// labor market
	LMkt abstraction.LaborMarket

	// capital market
	CMkt abstraction.CapitalMarket

	/**
	 * capital owned by the firm
	 */
	Capital *good.Capital

	/**
	 * quantity of capital
	 */
	CapitalQty float64

	/**
	 * present value of capital
	 */
	CapitalVal float64

	/**
	 * used to calculate average profit
	 */
	// TODO: private Averager pfAvger;

	/**
	 * Firm prop for general props between firm types
	 */
	Firm *Firm
}

type Averager struct {
	Sum  float64 // sum of data
	Size int     // buffer size
	Data []float64
}

func NewAverager(size int) *Averager {
	return &Averager{
		Data: []float64{},
		Size: size,
		Sum:  0,
	}
}

func (averager *Averager) Update(val float64) float64 {
	averager.Data = append(averager.Data, val)
	averager.Sum += val
	if len(averager.Data) > averager.Size {
		averager.Sum -= 0 // TODO: averager.Data.RemoveFirst()
	}
	return averager.Sum / float64(len(averager.Data))
}

func NewConsumedGoodsFirm(
	productName string,
	initWalletBal float64,
	initOutput float64,
	initWageBudget float64,
	initCapital int,
	capitalProducers []*CapitalFirm,
	pMkt abstraction.ConsumedGoodsMarket,
	lMkt abstraction.LaborMarket,
	cMkt abstraction.CapitalMarket,
) *ConsumedGoodsFirm {
	firm := NewFirm(initWalletBal)
	producerIDs := []string{}
	for _, cp := range capitalProducers {
		producerIDs = append(producerIDs, cp.Firm.ID)
	}
	capital := good.NewCapital(initCapital, firm.ID, producerIDs)
	firm.Output = initOutput
	firm.WageBudget = initWageBudget
	firm.Loan = 0
	firm.CapitalCost = 0

	// TODO: init pfAvger = new Averager(AVG_PROFIT_WIN)

	// post wage to the labor market so that the firm
	// gets employees before the first round begins
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[firm.ID]

	consumedGoodsFirm := &ConsumedGoodsFirm{
		Firm:        firm,
		Capital:     capital,
		ProductName: productName,
		PMkt:        pMkt,
		LMkt:        lMkt,
		CMkt:        cMkt,
	}
	consumedGoodsFirm.LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, firm.WageBudget)
	return consumedGoodsFirm
}

func (cgf *ConsumedGoodsFirm) UseCapital() float64 {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	var cost float64 = 0
	capital := cgf.Capital
	machines := capital.Machines
	removingMachineIdxs := []int{}
	for i, m := range machines {
		m.RemainingLife -= 1
		cost += m.Price
		transactionManager.Pay(
			capital.OwnerID,
			m.ProducerID,
			m.Price,
			common.PRIIC,
		)
		if m.RemainingLife == 0 {
			removingMachineIdxs = append(removingMachineIdxs, i)
			capital.Quantity -= 1
		}
	}

	remainingMachines := []*good.Machine{}
	for i, m := range machines {
		if !util.IsExisted(removingMachineIdxs, i) {
			remainingMachines = append(remainingMachines, m)
		}
	}
	capital.Machines = remainingMachines

	return cost
}

/**
 * Return output produced by labor amount of labor and c
 * amount of capital
 *
 * @param labor
 *            amount of labor
 * @param c
 *            amount of capital
 * @return output produced by labor amount of labor and c
 *         amount of capital
 */
func (cgf *ConsumedGoodsFirm) ConvertToProduct(labor, c float64) float64 {
	return cgf.TechCoefficient * math.Pow(labor, cgf.Beta) * math.Pow(c, 1-cgf.Beta)
}

func (cgf *ConsumedGoodsFirm) Act() {
	var newOutput, newWageBudget, pPrice float64
	firm := cgf.Firm
	laborQty := firm.Labor.Quantity
	capitalQty := cgf.Capital.Quantity

	// get firm finance information
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[firm.ID]

	firm.Revenue = walletAcc.PriIC
	firm.Loan = 0                                       // TODO: because Bank has not existed yet
	firm.TotalCost = firm.WageBudget + firm.CapitalCost // - acct.interest (bank loan interest amount)
	firm.Profit = firm.Revenue - firm.TotalCost

	firm.Capacity = cgf.ConvertToProduct(laborQty, capitalQty)
	if laborQty > 0 {
		firm.Wage = firm.WageBudget / laborQty
	} else {
		firm.Wage = 0
	}

	if laborQty > 0 {
		if common.TimeStep == 0 {
			newOutput = firm.Output
			newWageBudget = firm.WageBudget
		} else {
			var moneyFlowGap float64 = walletAcc.Balance - firm.TotalCost

			// set new wage budget
			newWageBudget = firm.WageBudget + cgf.Lambda*moneyFlowGap
			newWageBudget = math.Max(0, newWageBudget)

			// TODO: pay interest on loans (if any)

			// compute marginal cost
			var MC float64 = firm.Wage / cgf.Beta * math.Pow(cgf.TechCoefficient, -1/cgf.Beta) * math.Pow(firm.Output, 1/cgf.Beta-1) * math.Pow(capitalQty, 1-1/cgf.Beta)

			pPrice = cgf.PMkt.GetMarketPrice() // product price
			firm.MarginalProfit = pPrice - MC  // marginal profit

			// set new output
			newOutput = firm.Output * (1 + cgf.Phi*firm.MarginalProfit/pPrice)
		}

		newOutput = math.Min(firm.Capacity, newOutput)
		if newOutput > 0 {
			cgf.Product.Increase(newOutput)
		}
	} else {
		newOutput = firm.Output
		newWageBudget = firm.WageBudget
	}

	//////////////////////////
	if cgf.Product.GetQuantity() > 0 {
		cgf.PMkt.AddSellOffer(cgf, cgf.Product.GetQuantity())
	}
	cgf.LMkt.AddEmployer(firm.ID, walletAcc.Address, firm.Labor, newWageBudget)

	// TODO: pay loan (if any)

	// firgure out buying capital decision
	buyCapital(cgf, newOutput, newWageBudget, walletAcc)
}

func buyCapital(
	cgf *ConsumedGoodsFirm,
	newOutput float64,
	newWageBudget float64,
	walletAcc *transaction_manager.WalletAccount,
) {
	cMkt := cgf.CMkt
	firm := cgf.Firm
	// laborQty := firm.Labor.Quantity
	// capitalQty := cgf.Capital.Quantity
	var oldCapitalVal float64 = cgf.CapitalVal
	cgf.CapitalQty = cgf.Capital.Quantity          // quantity of machines
	cgf.CapitalVal = cgf.Capital.GetPresentValue() // total present value of capital
	firm.CapitalCost = cgf.UseCapital()

	if common.TimeStep > 0 {
		var capitalToBuy int = 0 // number of machines to purchase
		var capitalPrice float64 = cMkt.GetAvgPrice()
		var IR float64 = 0                                                 // TODO: Bank.getLoanIR() // interest rate
		var IK float64 = firm.Profit / oldCapitalVal                       // rate of return on capital
		var utilization float64 = newOutput / firm.Capacity                // capacity utilization
		var MR float64 = walletAcc.PriIC / cgf.CapitalQty * (1 - cgf.Beta) // marginal
		// revenue on
		// capital

		// buy capital if rate of return on capital >= interest rate,
		// capacity utilization >= eUtilThreshold,
		// marginal revenue >= capital price
		if IK >= IR && utilization >= cgf.EUtilThreshold && MR >= capitalPrice {
			capitalToBuy += 1
		}

		// TODO: fix it
		var avgProfit float64 = 0 // TODO: pfAvger.update(math.Abs(firm.Profit))

		// hacking
		if IK >= IR && utilization > 0.8 && firm.Profit > 5*avgProfit && common.TimeStep > 2000 {
			capitalToBuy += 1
		}

		/*
		 * Buy less capital when it is making a lot of losses (if the firm
		 * still decides to expand in this case, which does not really make
		 * sense)
		 */
		if firm.Profit < -5*avgProfit && common.TimeStep > 2000 {
			capitalToBuy -= 1
		}
		/***************************************************************/

		// number of machines that are written off
		var scrapped float64 = cgf.CapitalQty - cgf.Capital.Quantity
		if scrapped > 0 {
			// replace scrapped machines

			var x float64 = cgf.CapitalQty - scrapped
			for {
				if firm.Output/cgf.ConvertToProduct(firm.Labor.Quantity, x) > cgf.RUtilThreshold && x < cgf.CapitalQty && walletAcc.PriIC/x > capitalPrice {
					x += 1
					continue
				}
				break
			}
			capitalToBuy += int(scrapped + x - float64(cgf.CapitalQty))
		}

		// buy capital if there's none left
		if cgf.Capital.Quantity < 1 {
			capitalToBuy = int(math.Max(1, float64(capitalToBuy)))
		}

		// post buy offer to capital market
		if capitalToBuy > 0 {
			cMkt.AddCapitalBuyOffer(cgf.Capital, capitalToBuy)
		}
	}

	if newOutput > 0 {
		firm.Output = newOutput
	}
	firm.WageBudget = newWageBudget
	walletAcc.PriIC = 0
	firm.Labor.Decrease(firm.Labor.Quantity) // clear unused labor
	// TODO: loan = -acct.getBalance(Bank.SAVINGS)
}

func (cgf *ConsumedGoodsFirm) GetGood(goodName string) abstraction.Good {
	return cgf.Product
}

func (cgf *ConsumedGoodsFirm) GetWalletAccountAddress() string {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[cgf.Firm.ID]
	return walletAcc.Address
}

func (cgf *ConsumedGoodsFirm) GetConsumption(goodName string) float64 {
	return 0.0
}

func (cgf *ConsumedGoodsFirm) GetID() string {
	return cgf.Firm.ID
}
