package models

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type DemandForEnjoyment struct{}
type DemandForNecessity struct{}

type Laborer struct {
	// enjoyment market
	// EMkt *marketModels.ConsumedGoodsMarket

	// // necessity market
	// NMkt *marketModels.ConsumedGoodsMarket

	// // labor market
	// LMkt *marketModels.LaborMarket

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

	// each agent has an unique ID that is also used as the bank account number
	ID string

	// agent status
	IsAlive bool
}
