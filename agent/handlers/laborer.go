package handlers

import (
	"fmt"
	"math"

	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	marketModels "github.com/ninjadotorg/SimEconBaseline1/market/models"
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

func NewLaborer(
	initEQty float64,
	initNQty float64,
	initBalance float64,
	initSavingsRate float64,
) *agentModels.Laborer {
	econ := economy.GetEconInstance()
	laborer := &agentModels.Laborer{
		ID:          util.NewUUID(),
		IsAlive:     true,
		Enjoyment:   &good.Enjoyment{Quantity: initEQty},
		Necessity:   &good.Necessity{Quantity: initNQty},
		SavingsRate: initSavingsRate,
		DemandForE:  &DemandForEnjoyment{},
		DemandForN:  &DemandForNecessity{},
	}

	econ.TransactionManager.OpenWalletAccount(
		laborer.ID,
		initBalance,
	)

	walletAcc := econ.TransactionManager.WalletAccounts[laborer.ID]
	LMkt := econ.GetMarket("Labor").(*marketModels.LaborMarket)
	LMkt.AddEmployee(laborer.ID, walletAcc.Address)
	return laborer
}

func (laborer *agentModels.Laborer) GetWalletAccountAddress() string {
	econ := economy.GetEconInstance()
	walletAcc := econ.TransactionManager.WalletAccounts[laborer.ID]
	return walletAcc.Address
}

func (laborer *agentModels.Laborer) GetGood(goodName string) good.Good {
	if goodName == "Necessity" {
		return laborer.Necessity
	}
	if goodName == "Enjoyment" {
		return laborer.Enjoyment
	}
	return nil
}

func (laborer *agentModels.Laborer) GetConsumption(goodName string) float64 {
	if goodName == "Necessity" {
		return laborer.NConsumption
	}
	if goodName == "Enjoyment" {
		return laborer.EConsumption
	}
	return 0
}

func (laborer *agentModels.Laborer) Act() {
	econ := economy.GetEconInstance()
	walletAcc := econ.TransactionManager.WalletAccounts[laborer.ID]
	laborer.Wage = walletAcc.PriIC
	laborer.Income = laborer.Wage + walletAcc.SecIC // TODO: plus interest amt from saving acc

	// not enough good to eat -> die
	if (laborer.Necessity - eatAmt) < eatAmt {
		laborer.IsAlive = false
		fmt.Printf("Laborer %d died with balance: %f", laborer.ID, walletAcc.Balance)
		econ.TransactionManager.CloseWalletAccount(laborer.ID)
		return
	}

	depositIR := 0.08 // TODO: get deposit interest rate from bank
	if econ.TimeStep > 0 {
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
	bankBal := 100 // TODO: hardcode the bal here, will get from real bank acc later

	targetSaving := laborer.Income + baseSavingsToIncomeRatio
	if laborer.HighRR > laborer.LowRR {
		targetSaving *= (depositIR-laborer.LowRR)/(laborer.LowRR-laborer.LowRR)*epsilon*2 + 1 - epsilon
	}
	targetConsumption := walletBal + bankBal - targetSaving

	// compute consumption
	if econ.TimeStep == 0 {
		laborer.Consumption = laborer.Income
	} else {
		laborer.Consumption = math.Min(math.Max(laborer.Consumption*(1-upsilon), targetConsumption), laborer.Consumption*(1+upsilon))
	}

	newDeposit := walletBal - laborer.Consumption
	// TODO: deposit this amt to Bank

	// compute saving rate
	laborer.SavingRate = (bankBal + newDeposit) / (walletBal + bankBal)

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

	EMkt = econ.GetMarket("Enjoyment").(*marketModels.ConsumedGoodsMarket)
	NMkt = econ.GetMarket("Necessity").(*marketModels.ConsumedGoodsMarket)
	LMkt = econ.GetMarket("Labor").(*marmarketModelsket.LaborMarket)
	// post buy offer to enjoyment market
	EMkt.AddBuyOffer(laborer, laborer.DemandForE)
	// post buy offer to necessity market
	NMkt.AddBuyOffer(laborer, laborer.DemandForN)
	// post labor market
	LMkt.AddEmployee(laborer.ID, walletAcc.Address)

	walletAcc.PriIC = 0
	walletAcc.SecIC = 0
	// walletAcc.Interest = 0
}

func (de *agentModels.DemandForEnjoyment) GetDemand(
	price float64,
	consumption float64,
) float64 {
	return consumption / price
}

func (de *agentModels.DemandForNecessity) GetDemand(
	price float64,
	consumption float64,
) float64 {
	return consumption / price
}
