package models

type CapitalFirm struct {
	/**
	 * technology coefficient in the production function
	 */
	TechCoefficient float64

	/**
	 * sensitivity of output to labor (power on L in the production function
	 */
	Beta float64

	/**
	 * capital price (fixed for now)
	 */
	Price float64

	/**
	 * Firm prop for general props between firm types
	 */
	Firm *Firm
}
