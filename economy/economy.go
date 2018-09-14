package economy

import (
	"fmt"

	agent "github.com/ninjadotorg/SimEconBaseline1/agent"
	market "github.com/ninjadotorg/SimEconBaseline1/market"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
)

type Economy struct {
	TimeStep             int
	Agents               []agent.Agent
	DeadAgents           []agent.Agent
	Markets              map[string]market.Market
	ConsumedGoodsMarkets []*market.ConsumedGoodsMarket
	TransactionManager   *transaction_manager.TransactionManager
}

var econ *Economy

func GetEconInstance() *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		TimeStep:             0,
		Agents:               []agent.Agent{},
		DeadAgents:           []agent.Agent{},
		Markets:              map[string]market.Market{},
		ConsumedGoodsMarkets: []*market.ConsumedGoodsMarket{},
		TransactionManager: &transaction_manager.TransactionManager{
			WalletAccounts: map[string]*transaction_manager.WalletAccount{},
		},
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
	econ.TimeStep += 1
}

func (econ *Economy) GetMarket(marketName string) market.Market {
	return econ.Markets[marketName]
}
