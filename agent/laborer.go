package agent

import (
	"fmt"
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
	"github.com/ninjadotorg/SimEconBaseline1/util"
)

const (
	// target necessity stock
	targetNStock = 26

	// base savings to wage ratio
	baseSavingsToIncomeRatio = 10

	// quantity of necessity consumed in each step
	eatAmt = 1.0

	// sensitivity of target savings to real interest rate
	epsilon = 0.1

	// max percentage change in consumption allowed in each step
	upsilon = 0.04
)

type DemandForEnjoyment struct{}
type DemandForNecessity struct{}

type Laborer struct {
	// enjoyment market
	EMkt abstraction.ConsumedGoodsMarket

	// necessity market
	NMkt abstraction.ConsumedGoodsMarket

	// labor market
	LMkt abstraction.LaborMarket

	// enjoyment good
	Enjoyment *good.Enjoyment

	// necessity good
	Necessity *good.Necessity

	// savings rate (portion of total income+savings that is saved in the last
	// step)
	SavingsRate float64

	// consumption (in coin)
	Consumption float64

	// consumption of enjoyment (in coin)
	EConsumption float64

	// consumption of necessity (in coin)
	NConsumption float64

	// minimum necessity (in real quantity) to buy in the current step
	MinN float64

	// lowest real interest rate seen
	LowRR float64

	// highest real interest rate seen
	HighRR float64

	// demand for enjoyment
	DemandForE *DemandForEnjoyment

	// demand for necessity
	DemandForN *DemandForNecessity

	// total income
	Income float64

	// wage from employment
	Wage float64

	// each agent has an unique ID that is also used to retrieve the bank account number
	ID string

	// agent status
	IsAlive bool
}

func NewLaborer(
	initEQty float64,
	initNQty float64,
	initBalance float64,
	initSavingsRate float64,
	eMkt abstraction.ConsumedGoodsMarket,
	nMkt abstraction.ConsumedGoodsMarket,
	lMkt abstraction.LaborMarket,
) *Laborer {
	laborer := &Laborer{
		ID:          util.NewUUID(),
		IsAlive:     true,
		Enjoyment:   &good.Enjoyment{Quantity: initEQty},
		Necessity:   &good.Necessity{Quantity: initNQty},
		SavingsRate: initSavingsRate,
		DemandForE:  &DemandForEnjoyment{},
		DemandForN:  &DemandForNecessity{},
		EMkt:        eMkt,
		NMkt:        nMkt,
		LMkt:        lMkt,
	}

	transactionManager := transaction_manager.GetTransactionManagerInstance()
	transactionManager.OpenWalletAccount(laborer.ID, initBalance)
	walletAcc := transactionManager.WalletAccounts[laborer.ID]
	laborer.LMkt.AddEmployee(laborer.ID, walletAcc.Address)

	return laborer
}

func (laborer *Laborer) GetWalletAccountAddress() string {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[laborer.ID]
	return walletAcc.Address
}

func (laborer *Laborer) GetGood(goodName string) abstraction.Good {
	if goodName == "Necessity" {
		return laborer.Necessity
	}
	if goodName == "Enjoyment" {
		return laborer.Enjoyment
	}
	return nil
}

func (laborer *Laborer) GetConsumption(goodName string) float64 {
	if goodName == "Necessity" {
		return laborer.NConsumption
	}
	if goodName == "Enjoyment" {
		return laborer.EConsumption
	}
	return 0
}

func (laborer *Laborer) PrintLastState() {
	fmt.Printf("Labor ID: %s \n", laborer.ID)
	fmt.Printf("**** Enjoyment Qty: %f \n", laborer.Enjoyment.Quantity)
	fmt.Printf("**** Necessity Qty: %f \n", laborer.Necessity.Quantity)
	fmt.Printf("**** Consumption: %f \n", laborer.Consumption)
	fmt.Printf("**** Consumption for enjoyment: %f \n", laborer.EConsumption)
	fmt.Printf("**** Consumption for necessity: %f \n", laborer.NConsumption)
	fmt.Printf("**** Minimum necessity (in real quantity) to buy in the current step: %f \n", laborer.MinN)
	fmt.Printf("**** Income: %f \n", laborer.Income)
	fmt.Printf("**** Wage: %f \n\n", laborer.Wage)
}

func (laborer *Laborer) Act() {
	transactionManager := transaction_manager.GetTransactionManagerInstance()
	walletAcc := transactionManager.WalletAccounts[laborer.ID]
	laborer.Wage = walletAcc.PriIC
	laborer.Income = laborer.Wage + walletAcc.SecIC // TODO: plus interest amt from saving acc

	// not enough good to eat -> die
	if (laborer.Necessity.Quantity - eatAmt) < eatAmt {
		laborer.IsAlive = false
		fmt.Printf("Laborer %s died with balance: %f", laborer.ID, walletAcc.Balance)
		transactionManager.CloseWalletAccount(laborer.ID)
		return
	}

	depositIR := 0.08 // TODO: get deposit interest rate from bank
	if common.TimeStep > 0 {
		if depositIR < laborer.LowRR {
			laborer.LowRR = depositIR
		}
		if depositIR > laborer.HighRR {
			laborer.HighRR = depositIR
		}
	} else { // initial step
		laborer.LowRR = depositIR
		laborer.HighRR = depositIR
	}

	walletBal := walletAcc.Balance
	bankBal := 100.0 // TODO: hardcode the bal here, will get from real bank acc later

	targetSaving := laborer.Income * baseSavingsToIncomeRatio
	if laborer.HighRR > laborer.LowRR {
		targetSaving *= (depositIR-laborer.LowRR)/(laborer.HighRR-laborer.LowRR)*epsilon*2 + 1 - epsilon
	}
	targetConsumption := walletBal + bankBal - targetSaving

	// compute consumption
	if common.TimeStep == 0 {
		laborer.Consumption = laborer.Income
	} else {
		laborer.Consumption = math.Min(math.Max(laborer.Consumption*(1-upsilon), targetConsumption), laborer.Consumption*(1+upsilon))
	}

	newDeposit := walletBal - laborer.Consumption
	// TODO: deposit this amt to Bank

	// compute saving rate
	laborer.SavingsRate = (bankBal + newDeposit) / (walletBal + bankBal)

	// compute consumption of necessity
	laborer.NConsumption = laborer.Consumption * math.Max(0, 1-laborer.Necessity.Quantity/targetNStock)

	// compute consumption of enjoyment
	laborer.EConsumption = laborer.Consumption - laborer.NConsumption

	// if laborer has only 1 unit of necessity left, buy at least 1
	if laborer.Necessity.Quantity < 2 {
		laborer.MinN = 2 * eatAmt
	} else {
		laborer.MinN = 0
	}

	// post buy offer to enjoyment market
	laborer.EMkt.AddBuyOffer(laborer, laborer.DemandForE)
	// post buy offer to necessity market
	laborer.NMkt.AddBuyOffer(laborer, laborer.DemandForN)
	// post labor market
	laborer.LMkt.AddEmployee(laborer.ID, walletAcc.Address)

	walletAcc.PriIC = 0
	walletAcc.SecIC = 0
	// walletAcc.Interest = 0
}

func (de *DemandForEnjoyment) GetDemand(
	price float64,
	consumption float64,
) float64 {
	return consumption / price
}

func (de *DemandForNecessity) GetDemand(
	price float64,
	consumption float64,
) float64 {
	return consumption / price
}

func (laborer *Laborer) GetID() string {
	return laborer.ID
}
