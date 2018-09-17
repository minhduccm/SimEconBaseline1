package main

import (
	"fmt"

	agent "github.com/ninjadotorg/SimEconBaseline1/agent"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	market "github.com/ninjadotorg/SimEconBaseline1/market"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
)

const (
	// number of steps to run
	NUM_STEP = 2000

	NUM_LABORERS = 450
	NUM_EFIRMS   = 10
	NUM_NFIRMS   = 10

	MIN_INIT_E_PRICE = 0.1
	MAX_INIT_E_PRICE = 5
	MIN_INIT_N_PRICE = 0.1
	MAX_INIT_N_PRICE = 5

	EFIRM_INIT_CHECKING = 1000
	// EFIRM_INIT_SAVINGS    = -1000
	EFIRM_INIT_OUTPUT     = 40
	EFIRM_INIT_WAGEBUDGET = 100
	EFIRM_INIT_CAPITAL    = 30

	NFIRM_INIT_CHECKING = 100
	// NFIRM_INIT_SAVINGS    = -1000
	NFIRM_INIT_OUTPUT     = 50
	NFIRM_INIT_WAGEBUDGET = 100
	NFIRM_INIT_CAPITAL    = 30

	CFIRM_INIT_WAGEBUDGET = 500
	CFIRM_INIT_CHECKING   = CFIRM_INIT_WAGEBUDGET
	// CFIRM_INIT_SAVINGS    = 0

	LABORER_INIT_E        = 0
	LABORER_INIT_CHECKING = 0
	// LABORER_INIT_SAVINGS      = 100
	LABORER_INIT_SAVINGS_RATE = 0.9
)

func main() {
	fmt.Println("hahaha")
	econ := economy.GetEconInstance()
	_ = transaction_manager.GetTransactionManagerInstance()

	// Create and add markets
	eMkt := market.NewConsumedGoodsMarket(
		"Enjoyment",
		MIN_INIT_E_PRICE,
		MAX_INIT_E_PRICE,
	)

	nMkt := market.NewConsumedGoodsMarket(
		"Necessity",
		MIN_INIT_N_PRICE,
		MAX_INIT_N_PRICE,
	)

	lMkt := market.NewLaborMarket()
	cMkt := market.NewCapitalMarket()

	// add markets to economy
	econ.Markets["Labor"] = lMkt
	econ.Markets["Capital"] = cMkt
	econ.Markets["Enjoyment"] = eMkt
	econ.Markets["Necessity"] = nMkt

	// Create and add firms
	cFirm := agent.NewCapitalFirm(
		CFIRM_INIT_CHECKING,
		CFIRM_INIT_WAGEBUDGET,
		lMkt,
		cMkt,
	)
	cFirms := []*agent.CapitalFirm{cFirm}

	eFirms := []*agent.EnjoymentFirm{}
	nFirms := []*agent.NecessityFirm{}
	for i := 0; i < NUM_EFIRMS; i++ {
		eFirm := agent.NewEnjoymentFirm(
			EFIRM_INIT_CHECKING,
			EFIRM_INIT_OUTPUT,
			EFIRM_INIT_WAGEBUDGET,
			EFIRM_INIT_CAPITAL,
			cFirms,
			eMkt,
			lMkt,
			cMkt,
		)
		eFirms = append(eFirms, eFirm)
	}
	for i := 0; i < NUM_NFIRMS; i++ {
		nFirm := agent.NewNecessityFirm(
			NFIRM_INIT_CHECKING,
			NFIRM_INIT_OUTPUT,
			NFIRM_INIT_WAGEBUDGET,
			NFIRM_INIT_CAPITAL,
			cFirms,
			nMkt,
			lMkt,
			cMkt,
		)
		nFirms = append(nFirms, nFirm)
	}

	// add firm agents
	econ.Agents = append(econ.Agents, cFirm)
	for _, eFirm := range eFirms {
		econ.Agents = append(econ.Agents, eFirm)
	}
	for _, nFirm := range nFirms {
		econ.Agents = append(econ.Agents, nFirm)
	}

	// Create and add laborers
	for i := 0; i < NUM_LABORERS; i++ {
		var initN float64 = 15
		laborer := agent.NewLaborer(
			LABORER_INIT_E,
			initN,
			LABORER_INIT_CHECKING,
			LABORER_INIT_SAVINGS_RATE,
			eMkt,
			nMkt,
			lMkt,
		)
		econ.Agents = append(econ.Agents, laborer)
	}

	// perform
	lMkt.Perform()
	econ.Run(NUM_STEP)
}
