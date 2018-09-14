package market

import (
	"math"

	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
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

func NewLaborMarket() *LaborMarket {
	return &LaborMarket{
		Employees:   []*Employee{},
		Employers:   []*Employer{},
		TotalBudget: 0,
	}
}

func (laborMarket *LaborMarket) AddEmployee(agentID, walletAddress string) {
	employee := &Employee{
		AgentID:       agentID,
		WalletAddress: walletAddress,
	}
	laborMarket.Employees = append(laborMarket.Employees, employee)
}

func (laborMarket *LaborMarket) AddEmployer(
	agentID string,
	walletAddress string,
	labor *good.Labor,
	wageBudget float64,
) {
	employer := &Employer{
		Labor:         labor,
		WageBudget:    wageBudget,
		AgentID:       agentID,
		WalletAddress: walletAddress,
	}
	laborMarket.Employers = append(laborMarket.Employers, employer)
	laborMarket.TotalBudget += wageBudget
}

func (laborMarket *LaborMarket) Perform() {
	// TODO: should shuffle employers & employees

	econ := economy.GetEconInstance()
	var low int = 0
	var sum float64 = 0
	for _, employer := range laborMarket.Employers {
		sum += employer.WageBudget
		high := int(math.Min(1, sum/laborMarket.TotalBudget) * len(laborMarket.Employees))
		wage := employer.WageBudget / (high - low)
		for i := low; i < high; i++ {
			econ.TransactionManager.Pay(
				employer.AgentID,
				laborMarket.Employees[i].AgentID,
				wage,
				common.PRIIC,
			)
			employer.Labor.Increase(1)
		}
		low = high
	}
	// clear
	laborMarket.Employers = []*Employer{}
	laborMarket.Employees = []*Employee{}
	laborMarket.TotalBudget = 0
}
