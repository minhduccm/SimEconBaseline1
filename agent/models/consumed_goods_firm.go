package models

import "github.com/ninjadotorg/SimEconBaseline1/good"

type ConsumedGoodsFirm struct {
	/**
	 * product name: Necessity or Enjoyment
	 */
	ProductName string

	/**
	 * technology coefficient in the production function
	 */
	TechCoefficient float64

	/**
	 * sensitivity of output to labor (power on L in the production function
	 */
	Beta float64

	/**
	 * sensitivity of output to marginal profit
	 */
	Phi float64

	/**
	 * sensitivity of wage to money flow gap
	 */
	Lambda float64

	/**
	 * minimal capacity utilization to allow capital expansion
	 */
	EUtilThreshold float64

	/**
	 * minimal capacity utilization to allow capital replacement
	 */
	RUtilThreshold float64

	/**
	 * product the firm is producing/selling (enjoyment or necessity)
	 */
	Product good.Good

	/**
	 * capital owned by the firm
	 */
	Capital *good.Capital

	/**
	 * quantity of capital
	 */
	CapitalQty float64

	/**
	 * present value of capital
	 */
	CapitalVal float64

	/**
	 * used to calculate average profit
	 */
	// TODO: private Averager pfAvger;

	/**
	 * Firm prop for general props between firm types
	 */
	Firm *Firm
}

type Averager struct {
	Sum  float64 // sum of data
	Size int     // buffer size
	Data []float64
}
