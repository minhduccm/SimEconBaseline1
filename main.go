package main

import (
	"fmt"

	agentHandlers "github.com/ninjadotorg/SimEconBaseline1/agent/handlers"
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	marketHandlers "github.com/ninjadotorg/SimEconBaseline1/market/handlers"
	marketModels "github.com/ninjadotorg/SimEconBaseline1/market/models"
)

const (
	// number of steps between two printer outputs
	STEP_SIZE = 50

	// number of steps to run
	NUM_STEP = 10000

	NUM_LABORERS = 450
	NUM_EFIRMS   = 10
	NUM_NFIRMS   = 10

	MIN_INIT_E_PRICE = 0.1
	MAX_INIT_E_PRICE = 5
	MIN_INIT_N_PRICE = 0.1
	MAX_INIT_N_PRICE = 5

	EFIRM_INIT_CHECKING   = 100
	EFIRM_INIT_SAVINGS    = -1000
	EFIRM_INIT_OUTPUT     = 40
	EFIRM_INIT_WAGEBUDGET = 100
	EFIRM_INIT_CAPITAL    = 30

	NFIRM_INIT_CHECKING   = 100
	NFIRM_INIT_SAVINGS    = -1000
	NFIRM_INIT_OUTPUT     = 50
	NFIRM_INIT_WAGEBUDGET = 100
	NFIRM_INIT_CAPITAL    = 30

	CFIRM_INIT_WAGEBUDGET = 500
	CFIRM_INIT_CHECKING   = CFIRM_INIT_WAGEBUDGET
	CFIRM_INIT_SAVINGS    = 0

	LABORER_INIT_E            = 0
	LABORER_INIT_CHECKING     = 0
	LABORER_INIT_SAVINGS      = 100
	LABORER_INIT_SAVINGS_RATE = 0.9
)

func main() {
	fmt.Println("hahah")
	econ := economy.GetEconInstance()

	// Create and add markets
	eMkt := marketHandlers.NewConsumedGoodsMarket(
		"Enjoyment",
		MIN_INIT_E_PRICE,
		MAX_INIT_E_PRICE,
	)

	nMkt := marketHandlers.NewConsumedGoodsMarket(
		"Necessity",
		MIN_INIT_N_PRICE,
		MAX_INIT_N_PRICE,
	)

	lMkt := marketHandlers.NewLaborMarket()
	cMkt := marketHandlers.NewCapitalMarket()

	// add markets to economy
	econ.Markets["Labor"] = lMkt
	econ.Markets["Capital"] = cMkt
	econ.Markets["Enjoyment"] = eMkt
	econ.Markets["Necessity"] = nMkt
	econ.ConsumedGoodsMarkets = []*marketModels.ConsumedGoodsMarket{eMkt, nMkt}

	// Create and add firms
	cFirm := agentHandlers.NewCapitalFirm(
		CFIRM_INIT_CHECKING,
		CFIRM_INIT_WAGEBUDGET,
	)
	cFirms := []*agentModels.CapitalFirm{cFirm}

	eFirms := []*agentModels.EnjoymentFirm{}
	nFirms := []*agentModels.NecessityFirm{}
	for i := 0; i < NUM_EFIRMS; i++ {
		eFirm := agentHandlers.NewEnjoymentFirm(
			EFIRM_INIT_CHECKING,
			EFIRM_INIT_OUTPUT,
			EFIRM_INIT_WAGEBUDGET,
			EFIRM_INIT_CAPITAL,
			cFirms,
		)
		eFirms = append(eFirms, eFirm)
	}
	for i := 0; i < NUM_NFIRMS; i++ {
		nFirm := agentHandlers.NewEnjoymentFirm(
			NFIRM_INIT_CHECKING,
			NFIRM_INIT_OUTPUT,
			NFIRM_INIT_WAGEBUDGET,
			NFIRM_INIT_CAPITAL,
			cFirms,
		)
		nFirms = append(nFirms, nFirm)
	}

	// add firm agents
	econ.Agents = append(econ.Agents, cFirm)
	econ.Agents = append(econ.Agents, eFirms...)
	econ.Agents = append(econ.Agents, nFirms...)

	// Create and add laborers
	for i := 0; i < NUM_LABORERS; i++ {
		var initN float64 = 15
		laborer := agentHandlers.NewLaborer(
			LABORER_INIT_E,
			initN,
			LABORER_INIT_CHECKING,
			LABORER_INIT_SAVINGS_RATE,
		)
		econ.Agents = append(econ.Agents, laborer)
	}

	// perform
	lMkt.Perform()
	econ.Run(NUM_STEP)
}
