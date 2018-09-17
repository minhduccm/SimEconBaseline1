package economy

import (
	"fmt"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
	market "github.com/ninjadotorg/SimEconBaseline1/market"
)

type Economy struct {
	Agents               []abstraction.Agent
	DeadAgents           []abstraction.Agent
	Markets              map[string]abstraction.Market
	ConsumedGoodsMarkets []*market.ConsumedGoodsMarket
}

var econ *Economy

func GetEconInstance() *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		Agents:               []abstraction.Agent{},
		DeadAgents:           []abstraction.Agent{},
		Markets:              map[string]abstraction.Market{},
		ConsumedGoodsMarkets: []*market.ConsumedGoodsMarket{},
	}
	return econ
}

func (econ *Economy) Run(num_steps int) {
	for i := 0; i < num_steps; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		econ.Step()
	}
}

func (econ *Economy) Step() {
	for _, agent := range econ.Agents {
		agent.Act()
	}
	for _, market := range econ.Markets {
		market.Perform()
	}
	// econ.TimeStep += 1
	common.IncreaseTimeStep()
}

// func (econ *Economy) GetMarket(marketName string) market.Market {
// 	return econ.Markets[marketName]
// }
