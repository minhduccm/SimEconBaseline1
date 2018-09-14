package economy

import (
	"fmt"

	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	marketModels "github.com/ninjadotorg/SimEconBaseline1/market/models"
	"github.com/ninjadotorg/SimEconBaseline1/transaction_manager"
)

type Economy struct {
	TimeStep             int
	Agents               []agentModels.Agent
	DeadAgents           []agentModels.Agent
	Markets              map[string]marketModels.Market
	ConsumedGoodsMarkets []*marketModels.ConsumedGoodsMarket
	TransactionManager   *transaction_manager.TransactionManager
}

var econ *Economy

func GetEconInstance() *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		TimeStep:             0,
		Agents:               []agentModels.Agent{},
		DeadAgents:           []agentModels.Agent{},
		Markets:              map[string]marketModels.Market{},
		ConsumedGoodsMarkets: []*marketModels.ConsumedGoodsMarket{},
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

func (econ *Economy) GetMarket(marketName string) marketModels.Market {
	return econ.Markets[marketName]
}
