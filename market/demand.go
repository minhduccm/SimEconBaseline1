package market

type Demand interface {
	GetDemand(price float64, consumption float64) float64
}
