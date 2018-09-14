package models

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type Employee struct {
	WalletAddress string
}

type Employer struct {
	Labor         *good.Labor
	WageBudget    float64 // total wage budget
	Name          string  // name of the employer
	WalletAddress string
}

type LaborMarket struct {
	Employees   []*Employee
	Employers   []*Employer
	TotalBudget float64
}
