package models

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type Firm struct {
	/**
	 *  Firm ID
	 */
	ID string

	/**
	 *  labor owned by the firm
	 */
	Labor *good.Labor

	/**
	 *  max output the firm could produce with the current capital and labor
	 */
	Capacity float64

	/**
	 *  output in the last step
	 */
	Output float64

	/**
	 *  total wage budget in the last step
	 */
	WageBudget float64

	/**
	 *  wage (per worker) in the last step
	 */
	Wage float64

	/**
	 *  total loan in the last step
	 */
	Loan float64

	/**
	 *  revenue in the last step
	 */
	Revenue float64

	/**
	 *  profit in the last step
	 */
	Profit float64

	/**
	 *  marginal profit in the last step
	 */
	MarginalProfit float64

	/**
	 *  cost of capital in the last step
	 */
	CapitalCost float64

	/**
	 *  total cost in the last step
	 */
	TotalCost float64
}
