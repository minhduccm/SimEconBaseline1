package economy

import (
	"fmt"

	"github.com/ninjadotorg/SimEconBaseline1/abstraction"
	"github.com/ninjadotorg/SimEconBaseline1/common"
)

type Economy struct {
	Agents     []abstraction.Agent
	DeadAgents []abstraction.Agent
	Markets    map[string]abstraction.Market
}

var econ *Economy

func GetEconInstance() *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		Agents:     []abstraction.Agent{},
		DeadAgents: []abstraction.Agent{},
		Markets:    map[string]abstraction.Market{},
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
	fmt.Printf("\n--------------------Step %d------------------------ \n", common.TimeStep)
	for _, agent := range econ.Agents {
		agent.PrintLastState()
		agent.Act()
	}
	for _, market := range econ.Markets {
		market.Perform()
	}
	common.IncreaseTimeStep()
}
